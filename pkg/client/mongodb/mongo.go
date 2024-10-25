package mongodb

import (
	"context"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)


// InitMongoDB initializes a MongoDB client and returns a client and a cleanup function.
// The cleanup function must be called when the client is no longer needed.
func InitMongoDB(uri string) (*mongo.Client, func()) {
	ctx := context.Background()
	// Create a new MongoDB client.
	cli, err := mongo.Connect(ctx,
		options.Client().ApplyURI(uri),
		options.Client().SetMaxPoolSize(100),
		options.Client().SetMaxConnIdleTime(10*time.Second),
		options.Client().SetAuth(options.Credential{Username: "root", Password: "RSeMAN7csL"}),
	)
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

func NewMongo(url string, dbName string) *Mongo {
	cli, cleanup := InitMongoDB(url)
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
