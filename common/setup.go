package common

import (
	"context"
	"os"
	"sync"

	"github.com/ydssx/kratos-kit/common/conf"
	"github.com/ydssx/kratos-kit/constants"
	"github.com/ydssx/kratos-kit/pkg/client/mongodb"
	"github.com/ydssx/kratos-kit/pkg/client/mysql"
	"github.com/ydssx/kratos-kit/pkg/client/redis"
	"github.com/ydssx/kratos-kit/pkg/email"
	"github.com/ydssx/kratos-kit/pkg/limit"
	"github.com/ydssx/kratos-kit/pkg/lock"
	"github.com/ydssx/kratos-kit/pkg/logger"
	"github.com/ydssx/kratos-kit/pkg/middleware/kratos"
	"github.com/ydssx/kratos-kit/pkg/queue"
	"github.com/ydssx/kratos-kit/pkg/storage"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/hibiken/asynq"
	"github.com/oschwald/geoip2-golang"
	goredis "github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"gorm.io/gorm"
)

func NewRedisCLient(c *conf.Bootstrap) (*goredis.Client, error) {
	redisConf := c.Data.GetRedis()
	return redis.NewRedis(&goredis.Options{
		Addr:         redisConf.GetAddr(),
		Password:     redisConf.GetPassword(),
		Username:     redisConf.GetUsername(),
		ReadTimeout:  redisConf.GetReadTimeout().AsDuration(),
		WriteTimeout: redisConf.GetWriteTimeout().AsDuration(),
		DialTimeout:  redisConf.GetDialTimeout().AsDuration(),
		DB:           int(redisConf.GetDb()),
	})
}

func NewMysqlDB(c *conf.Bootstrap) (*gorm.DB, error) {
	return mysql.NewDB(c.Data.GetDatabase().GetSource()...)
}

func NewCollection(c *conf.Bootstrap) *mongo.Collection {
	mongoConf := c.Data.GetMongo()
	db := mongodb.NewMongo(mongoConf.GetAddr(), mongoConf.GetDatabase())
	return db.Collection(mongoConf.GetCollection())
}

func SetupLogger(c *conf.Bootstrap) {
	// 使用 Option 函数配置 logger
	defaultLogger := logger.NewLogger(logger.NewZapLogger(
		logger.WithCallerSkip(3),
		logger.WithLogPath(c.Log.GetPath()),
		logger.WithLevel(c.Log.GetLevel()),
		logger.WithMaxSize(int(c.Log.GetMaxSize())),
		logger.WithMaxAge(int(c.Log.GetMaxAge())),
		logger.WithMaxBackups(int(c.Log.GetMaxBackups())),
		logger.WithCompress(c.Log.GetCompress()),
		logger.WithEnableConsole(c.Log.GetEnableConsole()),
		logger.WithWebhook(c.Webhook.GetUrl()),
	))
	logger.DefaultLogger = defaultLogger

	// 创建 kratos logger
	klogger := log.With(logger.NewLogger(logger.NewZapLogger(
		logger.WithCallerSkip(3),
		logger.WithLogPath(c.Log.GetPath()),
		logger.WithLevel(c.Log.GetLevel()),
		logger.WithMaxSize(int(c.Log.GetMaxSize())),
		logger.WithMaxAge(int(c.Log.GetMaxAge())),
		logger.WithMaxBackups(int(c.Log.GetMaxBackups())),
		logger.WithCompress(c.Log.GetCompress()),
		logger.WithEnableConsole(c.Log.GetEnableConsole()),
		logger.WithWebhook(c.Webhook.GetUrl()),
	)),
		"traceID", kratos.TraceID())
	log.SetLogger(klogger)
}

var (
	rdbClientOpt asynq.RedisClientOpt
	once         sync.Once
)

func InitRedisOpt(c *conf.Bootstrap) asynq.RedisClientOpt {
	redisConf := c.Data.GetRedis()
	once.Do(func() {
		rdbClientOpt = asynq.RedisClientOpt{
			Addr:     redisConf.GetAddr(),
			Password: redisConf.GetPassword(),
			DB:       int(redisConf.GetDb()),
		}
	})
	return rdbClientOpt
}

func InitJobRedisOpt(c *conf.Bootstrap) asynq.RedisClientOpt {
	redisConf := c.Data.GetJobRedis()

	return asynq.RedisClientOpt{
		Addr:     redisConf.GetAddr(),
		Password: redisConf.GetPassword(),
		DB:       int(redisConf.GetDb()),
	}
}

func NewQueueClient(c *conf.Bootstrap) *queue.Client {
	return queue.NewClient(&queue.ConnConfig{
		RedisAddr:     c.Data.GetRedis().GetAddr(),
		RedisPassword: c.Data.GetRedis().GetPassword(),
		RedisDB:       int(c.Data.GetRedis().GetDb()),
		ReadTimeout:   c.Data.GetRedis().GetReadTimeout().AsDuration(),
		WriteTimeout:  c.Data.GetRedis().GetWriteTimeout().AsDuration(),
	})
}

func NewGoogleCloudStorage(c *conf.Bootstrap) (*storage.GoogleCloudStorage, func()) {
	return storage.NewGoogleCloudStorage(c.Gcs.GetBucketName(), c.Gcs.GetProjectId(), c.Gcs.GetCredentialsFile())
}

// NewGeoipDB returns a new GeoipDB.
func NewGeoipDB(ctx context.Context, c *conf.Bootstrap) *geoip2.Reader {
	db, err := geoip2.Open(c.Data.Geoip.Path)
	if err != nil {
		log.Fatal("failed to open GeoIP database: ", err)
	}

	context.AfterFunc(ctx, func() {
		err := db.Close()
		if err != nil {
			log.Error("failed to close GeoIP database: ", err)
		}
	})

	return db
}

func NewRateLimiter(rdb *goredis.Client) *limit.RedisLimiter {
	return limit.NewRedisLimiter(rdb)
}

func NewRedisLocker(rdb *goredis.Client) *lock.RedisLocker {
	return lock.NewLocker(rdb)
}

// 初始化Google OAuth配置
func InitGoogleOAuth(c *conf.Bootstrap) *oauth2.Config {
	return &oauth2.Config{
		ClientID:     c.Google.GetClientId(),
		ClientSecret: c.Google.GetClientSecret(),
		RedirectURL:  c.Google.GetRedirectUrl(),
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
		},
		Endpoint: google.Endpoint,
	}
}

// 设置环境变量
func SetEnv(c *conf.Bootstrap) {
	os.Setenv(string(constants.EnvKeyDingDingWebhook), c.Webhook.GetUrl())
	os.Setenv(string(constants.EnvKeyLogPath), c.Log.GetPath())
	os.Setenv(string(constants.EnvKeyEnv), c.GetEnv())
}

func NewEmail(c *conf.Bootstrap) *email.Email {
	return email.NewEmail(
		c.Email.GetHost(),
		int(c.Email.GetPort()),
		c.Email.GetUsername(),
		c.Email.GetPassword(),
		c.Email.GetFrom(),
	)
}
