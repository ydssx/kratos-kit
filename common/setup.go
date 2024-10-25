package common

import (
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
	redisConf := c.Data.Redis
	return redis.NewRedis(&goredis.Options{
		Addr:         redisConf.Addr,
		Password:     redisConf.Password,
		Username:     redisConf.Username,
		ReadTimeout:  redisConf.ReadTimeout.AsDuration(),
		WriteTimeout: redisConf.WriteTimeout.AsDuration(),
		DialTimeout:  redisConf.DialTimeout.AsDuration(),
		DB:           int(redisConf.Db),
	})
}

func NewMysqlDB(c *conf.Bootstrap) (*gorm.DB, error) {
	return mysql.NewDB(c.Data.Database.Source...)
}

func NewCollection(c *conf.Bootstrap) *mongo.Collection {
	mongoConf := c.Data.Mongo
	db := mongodb.NewMongo(mongoConf.Addr, mongoConf.Database)
	return db.Collection(mongoConf.Collection)
}

func SetupLogger(c *conf.Bootstrap) {
	logger.DefaultLogger = logger.NewLogger(logger.NewZapLogger(2,
		c.Log.Path,
		c.Log.Level,
		int(c.Log.MaxSize),
		int(c.Log.MaxAge),
		int(c.Log.MaxBackups),
		c.Log.Compress,
		c.Log.EnableConsole,
	))

	klogger := log.With(logger.NewLogger(logger.NewZapLogger(3,
		c.Log.Path,
		c.Log.Level,
		int(c.Log.MaxSize),
		int(c.Log.MaxAge),
		int(c.Log.MaxBackups),
		c.Log.Compress,
		c.Log.EnableConsole,
	)),
		"traceID", kratos.TraceID())
	log.SetLogger(klogger)
}

var (
	rdbClientOpt asynq.RedisClientOpt
	once         sync.Once
)

func InitRedisOpt(c *conf.Bootstrap) asynq.RedisClientOpt {
	redisConf := c.Data.Redis
	once.Do(func() {
		rdbClientOpt = asynq.RedisClientOpt{
			Addr:     redisConf.Addr,
			Password: redisConf.Password,
			DB:       int(redisConf.Db),
		}
	})
	return rdbClientOpt
}

func InitJobRedisOpt(c *conf.Bootstrap) asynq.RedisClientOpt {
	redisConf := c.Data.JobRedis

	return asynq.RedisClientOpt{
		Addr:     redisConf.Addr,
		Password: redisConf.Password,
		DB:       int(redisConf.Db),
	}
}

func NewAsynqClient(c *conf.Bootstrap) *asynq.Client {
	return asynq.NewClient(InitRedisOpt(c))
}

func NewAsynqInspector(c *conf.Bootstrap) *asynq.Inspector {
	return asynq.NewInspector(InitRedisOpt(c))
}

func NewGoogleCloudStorage(c *conf.Bootstrap) (*storage.GoogleCloudStorage, func()) {
	return storage.NewGoogleCloudStorage(c.Gcs.BucketName, c.Gcs.ProjectId, c.Gcs.CredentialsFile)
}

// NewGeoipDB returns a new GeoipDB.
func NewGeoipDB(c *conf.Bootstrap) (*geoip2.Reader, func()) {
	db, err := geoip2.Open(c.Data.Geoip.Path)
	if err != nil {
		log.Fatal("failed to open GeoIP database: ", err)
	}
	return db, func() {
		db.Close()
	}
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
		ClientID:     c.Google.ClientId,
		ClientSecret: c.Google.ClientSecret,
		RedirectURL:  c.Google.RedirectUrl,
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
	return email.NewEmail(c.Email.Host, int(c.Email.Port), c.Email.Username, c.Email.Password, c.Email.From)
}
