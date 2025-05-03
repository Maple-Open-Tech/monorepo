// cloud/backend/internal/ipe/repo/evaluation/impl.go
package evaluation

import (
	"context"
	"log"

	"go.uber.org/zap"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"

	"github.com/Maple-Open-Tech/monorepo/cloud/backend/config"
	dom_evaluation "github.com/Maple-Open-Tech/monorepo/cloud/backend/internal/ipe/domain/evaluation"
)

type evaluationImpl struct {
	Logger                  *zap.Logger
	DbClient                *mongo.Client
	EvaluationCollection    *mongo.Collection
	NeighbourhoodCollection *mongo.Collection
	PropertyPhotoCollection *mongo.Collection
}

func NewRepository(appCfg *config.Configuration, loggerp *zap.Logger, client *mongo.Client) dom_evaluation.EvaluationRepository {
	// Get collections
	evalCollection := client.Database(appCfg.DB.IncomePropertyEvaluatorName).Collection("evaluations")
	neighbourhoodCollection := client.Database(appCfg.DB.IncomePropertyEvaluatorName).Collection("neighbourhoods")
	photoCollection := client.Database(appCfg.DB.IncomePropertyEvaluatorName).Collection("property_photos")

	// Create indexes for evaluation collection
	_, err := evalCollection.Indexes().CreateMany(context.TODO(), []mongo.IndexModel{
		{Keys: bson.D{{Key: "propertyId", Value: 1}}, Options: options.Index().SetUnique(true)},
		{Keys: bson.D{{Key: "clientId", Value: 1}}},
		{Keys: bson.D{{Key: "presenterId", Value: 1}}},
		{Keys: bson.D{{Key: "ownerId", Value: 1}}},
		{Keys: bson.D{{Key: "neighbourhoodId", Value: 1}}},
	})
	if err != nil {
		log.Fatal(err)
	}

	// Create indexes for neighbourhood collection
	_, err = neighbourhoodCollection.Indexes().CreateMany(context.TODO(), []mongo.IndexModel{
		{Keys: bson.D{{Key: "recordName", Value: 1}}},
		{Keys: bson.D{{Key: "recordUniqueId", Value: 1}}, Options: options.Index().SetUnique(true)},
		{Keys: bson.D{{Key: "city", Value: 1}}},
		{Keys: bson.D{{Key: "province", Value: 1}}},
		{Keys: bson.D{{Key: "country", Value: 1}}},
	})
	if err != nil {
		log.Fatal(err)
	}

	// Create indexes for property photo collection
	_, err = photoCollection.Indexes().CreateMany(context.TODO(), []mongo.IndexModel{
		{Keys: bson.D{{Key: "evaluationId", Value: 1}}},
		{Keys: bson.D{{Key: "photoCategory", Value: 1}}},
		{Keys: bson.D{{Key: "photoUniqueId", Value: 1}}, Options: options.Index().SetUnique(true)},
	})
	if err != nil {
		log.Fatal(err)
	}

	return &evaluationImpl{
		Logger:                  loggerp,
		DbClient:                client,
		EvaluationCollection:    evalCollection,
		NeighbourhoodCollection: neighbourhoodCollection,
		PropertyPhotoCollection: photoCollection,
	}
}
