package main

import (
	"context"
	"fmt"

	"github.com/nedvisol/go-connectdots/cacheditem"
	"github.com/nedvisol/go-connectdots/config"
	"github.com/nedvisol/go-connectdots/downloadmgr"
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

func AppStart(ctx context.Context, congressGov *processor.CongressGovProcessor) {
	// req, _ := http.NewRequest("GET", "http://ned1:3000/hello.json", nil)
	// dmgr.Download(ctx, req, func(data []byte) { fmt.Printf("downloaded %s\n", string(data)) })
	// dmgr.Download(ctx, req, func(data []byte) { fmt.Printf("downloaded %s\n", string(data)) })
	// dmgr.Download(ctx, req, func(data []byte) { fmt.Printf("downloaded %s\n", string(data)) })
	// dmgr.Download(ctx, req, func(data []byte) { fmt.Printf("downloaded %s\n", string(data)) })
	// dmgr.Download(ctx, req, func(data []byte) { fmt.Printf("downloaded %s\n", string(data)) })
	// dmgr.Download(ctx, req, func(data []byte) { fmt.Printf("downloaded %s\n", string(data)) })
	// dmgr.Download(ctx, req, func(data []byte) { fmt.Printf("downloaded %s\n", string(data)) })

	//dmgr.Wait()

	congressGov.Start()
}

func main() {

	app := fx.New(
		fx.Provide(
			config.NewConfig,
			context.Background,
			NewMongoClient,
			NewMongoDatabase,
			NewDownloadManagerOptions,
			downloadmgr.NewDownloadManager,
			processor.NewCongressGovProcessor,
		),
		fx.Invoke(AppStart),
	)

	app.Run()

}
