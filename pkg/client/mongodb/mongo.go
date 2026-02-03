package mongodb

import (
	"context"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// MongoConfig holds MongoDB connection configuration
type MongoConfig struct {
	URI             string
	Username        string
	Password        string
	MaxPoolSize     uint64
	MaxConnIdleTime time.Duration
}

// InitMongoDB initializes a MongoDB client and returns a client and a cleanup function.
// The cleanup function must be called when the client is no longer needed.
func InitMongoDB(config MongoConfig) (*mongo.Client, func()) {
	ctx := context.Background()

	clientOptions := []*options.ClientOptions{
		options.Client().ApplyURI(config.URI),
		options.Client().SetMaxPoolSize(config.MaxPoolSize),
		options.Client().SetMaxConnIdleTime(config.MaxConnIdleTime),
	}

	// Only set auth if credentials are provided
	if config.Username != "" && config.Password != "" {
		clientOptions = append(clientOptions,
			options.Client().SetAuth(options.Credential{
				Username: config.Username,
				Password: config.Password,
			}))
	}

	// Create a new MongoDB client.
	cli, err := mongo.Connect(ctx, clientOptions...)
	if err != nil {
		panic("failed to connect to MongoDB: " + err.Error())
	}

	// Define a cleanup function to disconnect from MongoDB.
	cleanup := func() { cli.Disconnect(ctx) }

	// Ping the primary to test the connection.
	err = cli.Ping(ctx, readpref.Primary())
	if err != nil {
		cleanup()
		panic("failed to ping MongoDB: " + err.Error())
	}

	log.Info("init mongodb success")
	return cli, cleanup
}

type Mongo struct {
	db      *mongo.Database
	cli     *mongo.Client
	cleanup func()
}

// NewMongo creates a new Mongo instance with default configuration
// Deprecated: Use NewMongoWithConfig for better configuration control
func NewMongo(url string, dbName string) *Mongo {
	config := MongoConfig{
		URI:             url,
		MaxPoolSize:     100,
		MaxConnIdleTime: 10 * time.Second,
	}
	cli, cleanup := InitMongoDB(config)
	db := cli.Database(dbName)
	return &Mongo{db: db, cli: cli, cleanup: cleanup}
}

// NewMongoWithConfig creates a new Mongo instance with custom configuration
func NewMongoWithConfig(config MongoConfig, dbName string) *Mongo {
	cli, cleanup := InitMongoDB(config)
	db := cli.Database(dbName)
	return &Mongo{db: db, cli: cli, cleanup: cleanup}
}

func (m *Mongo) Close() {
	m.cleanup()
}

func (m *Mongo) Database() *mongo.Database {
	return m.db
}

func (m *Mongo) Collection(collectionName string) *mongo.Collection {
	return m.db.Collection(collectionName)
}

func InitMongoCollection(db *mongo.Database, collectionName string) *mongo.Collection {
	return db.Collection(collectionName)
}

func InitMongoDatabase(cli *mongo.Client, dbName string) *mongo.Database {
	return cli.Database(dbName)
}
