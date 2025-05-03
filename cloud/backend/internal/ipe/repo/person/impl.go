// cloud/backend/internal/ipe/repo/person/impl.go
package person

import (
	"context"
	"log"

	"go.uber.org/zap"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"

	"github.com/Maple-Open-Tech/monorepo/cloud/backend/config"
	dom_person "github.com/Maple-Open-Tech/monorepo/cloud/backend/internal/ipe/domain/person"
)

type personImpl struct {
	Logger              *zap.Logger
	DbClient            *mongo.Client
	ClientCollection    *mongo.Collection
	PresenterCollection *mongo.Collection
	OwnerCollection     *mongo.Collection
	CompareCollection   *mongo.Collection
}

func NewRepository(appCfg *config.Configuration, loggerp *zap.Logger, client *mongo.Client) dom_person.PersonRepository {
	// Get collections
	clientCollection := client.Database(appCfg.DB.IncomePropertyEvaluatorName).Collection("clients")
	presenterCollection := client.Database(appCfg.DB.IncomePropertyEvaluatorName).Collection("presenters")
	ownerCollection := client.Database(appCfg.DB.IncomePropertyEvaluatorName).Collection("owners")
	compareCollection := client.Database(appCfg.DB.IncomePropertyEvaluatorName).Collection("compares")

	// Create indexes for client collection
	_, err := clientCollection.Indexes().CreateMany(context.TODO(), []mongo.IndexModel{
		{Keys: bson.D{{Key: "personName", Value: 1}}},
		{Keys: bson.D{{Key: "email", Value: 1}}},
		{Keys: bson.D{
			{Key: "personName", Value: "text"},
			{Key: "email", Value: "text"},
			{Key: "address", Value: "text"},
		}},
	})
	if err != nil {
		log.Fatal(err)
	}

	// Create indexes for presenter collection
	_, err = presenterCollection.Indexes().CreateMany(context.TODO(), []mongo.IndexModel{
		{Keys: bson.D{{Key: "personName", Value: 1}}},
		{Keys: bson.D{{Key: "email", Value: 1}}},
		{Keys: bson.D{
			{Key: "personName", Value: "text"},
			{Key: "email", Value: "text"},
			{Key: "address", Value: "text"},
		}},
	})
	if err != nil {
		log.Fatal(err)
	}

	// Create indexes for owner collection
	_, err = ownerCollection.Indexes().CreateMany(context.TODO(), []mongo.IndexModel{
		{Keys: bson.D{{Key: "personName", Value: 1}}},
		{Keys: bson.D{{Key: "email", Value: 1}}},
		{Keys: bson.D{
			{Key: "personName", Value: "text"},
			{Key: "email", Value: "text"},
			{Key: "address", Value: "text"},
		}},
	})
	if err != nil {
		log.Fatal(err)
	}

	// Create indexes for compare collection
	_, err = compareCollection.Indexes().CreateMany(context.TODO(), []mongo.IndexModel{
		{Keys: bson.D{{Key: "clientId", Value: 1}}},
		{Keys: bson.D{{Key: "presenterId", Value: 1}}},
		{Keys: bson.D{{Key: "recordUniqueId", Value: 1}}},
	})
	if err != nil {
		log.Fatal(err)
	}

	return &personImpl{
		Logger:              loggerp,
		DbClient:            client,
		ClientCollection:    clientCollection,
		PresenterCollection: presenterCollection,
		OwnerCollection:     ownerCollection,
		CompareCollection:   compareCollection,
	}
}
