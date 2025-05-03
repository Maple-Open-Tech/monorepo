// cloud/backend/internal/ipe/repo/incomeproperty/impl.go
package incomeproperty

import (
	"context"
	"log"

	"go.uber.org/zap"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"

	"github.com/Maple-Open-Tech/monorepo/cloud/backend/config"
	dom_property "github.com/Maple-Open-Tech/monorepo/cloud/backend/internal/ipe/domain/incomeproperty"
)

type incomePropertyImpl struct {
	Logger     *zap.Logger
	DbClient   *mongo.Client
	Collection *mongo.Collection
}

func NewRepository(appCfg *config.Configuration, loggerp *zap.Logger, client *mongo.Client) dom_property.PropertyRepository {
	// Get collection
	collection := client.Database(appCfg.DB.IncomePropertyEvaluatorName).Collection("income_properties")

	// Create indexes
	_, err := collection.Indexes().CreateMany(context.TODO(), []mongo.IndexModel{
		{Keys: bson.D{{Key: "address", Value: 1}}},
		{Keys: bson.D{{Key: "city", Value: 1}}},
		{Keys: bson.D{{Key: "province", Value: 1}}},
		{Keys: bson.D{{Key: "country", Value: 1}}},
		{Keys: bson.D{{Key: "propertyCode", Value: 1}}},
		{Keys: bson.D{
			{Key: "address", Value: "text"},
			{Key: "city", Value: "text"},
			{Key: "recordName", Value: "text"},
		}},
	})
	if err != nil {
		// Fatal error on startup to meet requirements of google/wire framework
		log.Fatal(err)
	}

	return &incomePropertyImpl{
		Logger:     loggerp,
		DbClient:   client,
		Collection: collection,
	}
}
