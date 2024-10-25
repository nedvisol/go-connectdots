package main

import (
	"context"
	"fmt"

	"github.com/nedvisol/go-connectdots/cacheditem"
	"github.com/nedvisol/go-connectdots/config"
	"github.com/nedvisol/go-connectdots/downloadmgr"
	"github.com/nedvisol/go-connectdots/graphdb"
	"github.com/nedvisol/go-connectdots/processor"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/fx"
)

func NewMongoClient(ctx context.Context, config *config.Config) *mongo.Client {
	clientOptions := options.Client().ApplyURI(config.MongoUrl) // Replace with your MongoDB URI

	// Connect to MongoDB
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		panic(err)
	}

	// Ping the database to ensure it's connected
	err = client.Ping(ctx, nil)
	if err != nil {
		panic(err)
	}
	fmt.Println("Connected to MongoDB!")
	return client
}

func NewMongoDatabase(client *mongo.Client, config *config.Config) *mongo.Database {
	return client.Database(config.MongoDb)
}

func NewDownloadManagerOptions(db *mongo.Database, config *config.Config) *downloadmgr.DownloadManagerOptions {
	return &downloadmgr.DownloadManagerOptions{
		CacheDir:       config.CacheDir,
		CachedItemRepo: cacheditem.NewCachedItemMongoRepository(db),
		Config:         config,
	}
}

func AppStart(lifecycle fx.Lifecycle, ctx context.Context, congressGov *processor.CongressGovProcessor) {

	lifecycle.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			fmt.Println("Application is stopping. Cleaning up resources...")
			// Perform cleanup actions here (close connections, release resources, etc.)
			return nil
		},
	})

	congressGov.Start()
}

func main() {

	//ctx, _ := context.WithCancel(context.Background())

	app := fx.New(
		fx.Provide(
			config.NewConfig,
			context.Background,
			NewMongoClient,
			NewMongoDatabase,
			NewDownloadManagerOptions,
			downloadmgr.NewDownloadManager,
			graphdb.NewNeo4jGraphService,
			processor.NewCongressGovProcessor,
		),
		fx.Invoke(AppStart),
	)

	app.Run()

}
