package mysql

import (
	"context"
	"database/sql"
	"time"

	"github.com/ydssx/kratos-kit/pkg/errors"
	"github.com/ydssx/kratos-kit/pkg/logger"

	"go.uber.org/zap/zapcore"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/plugin/dbresolver"
)

var db *gorm.DB

// NewDB initializes a new MySQL database connection pool and returns the gorm.DB instance.
// It takes the MySQL DSN as a parameter.
// It configures the gorm logger, prepares statements, sets connection pool limits and logs success.
// Returns the gorm.DB instance and any error.
func NewDB(dsn ...string) (*gorm.DB, error) {
	if len(dsn) == 0 {
		return nil, errors.New("dsn is required")
	}
	dialectors := make([]gorm.Dialector, 0, len(dsn))
	for _, d := range dsn {
		dialectors = append(dialectors, mysql.Open(d))
	}

	var err error
	db, err = gorm.Open(dialectors[0], &gorm.Config{
		Logger:      NewGormLogger(zapcore.InfoLevel, zapcore.InfoLevel, time.Millisecond*200),
		PrepareStmt: true,
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed to connect to mysql")
	}
	if len(dialectors) > 1 {
		db.Use(dbresolver.Register(dbresolver.Config{
			Sources: dialectors[1:],
			Policy:  dbresolver.StrictRoundRobinPolicy(),
		}))
	}
	sqlDB, err := db.DB()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get mysql db")
	}
	sqlDB.SetMaxIdleConns(100)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)
	logger.Info(context.Background(), "init mysql success")
	return db, nil
}

func Transaction(fc func(tx *gorm.DB) error, opts ...*sql.TxOptions) error {
	return db.Transaction(fc, opts...)
}

func GlobalDB() *gorm.DB {
	return db
}

type contextTxKey struct{}

func NewContextWithDB(ctx context.Context) context.Context {
	return context.WithValue(ctx, contextTxKey{}, db)
}

func DBFromContext(ctx context.Context) *gorm.DB {
	return ctx.Value(contextTxKey{}).(*gorm.DB)
}
