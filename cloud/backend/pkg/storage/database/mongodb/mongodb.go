package mongodb

import (
	"context"
	"log"

	"log/slog"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"

	c "github.com/Maple-Open-Tech/monorepo/cloud/backend/config"
)

func NewProvider(appCfg *c.Configuration, logger *slog.Logger) *mongo.Client {
	logger.Debug("storage initializing...")

	// DEVELOPERS NOTE:
	// If you want to turn on more detailed loggined then uncomment the
	// follwoing code and then comment out the other code.
	// - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
	// loggerOptions := options.
	// 	Logger().
	// 	SetComponentLevel(options.LogComponentCommand, options.LogLevelDebug)
	//
	// client, err := mongo.Connect(
	// 	context.TODO(),
	// 	options.Client().
	// 		ApplyURI(appCfg.DB.URI).SetLoggerOptions(loggerOptions))
	//
	// loggerOptions := options.
	// 	Logger().
	// 	SetComponentLevel(options.LogComponentCommand, options.LogLevelDebug)
	// - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -

	// DEVELOPERS NOTE:
	// If you uncommented the ABOVE code then comment out the BOTTOM code.
	// - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
	client, err := mongo.Connect(
		context.TODO(), options.Client().ApplyURI(appCfg.DB.URI))
	// - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -

	if err != nil {
		log.Fatalf("comiccoin.common.storage.database.mongodb.NewProvider: Error: %v\n", err)
	}

	// The MongoDB client provides a Ping() method to tell you if a MongoDB database has been found and connected.
	if err := client.Ping(context.TODO(), readpref.Primary()); err != nil {
		log.Fatal(err)
	}
	logger.Debug("storage initialized successfully")
	return client
}
