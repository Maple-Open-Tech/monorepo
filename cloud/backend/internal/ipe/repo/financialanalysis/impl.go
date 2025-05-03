// cloud/backend/internal/ipe/repo/financialanalysis/impl.go
package financialanalysis

import (
	"context"
	"log"

	"go.uber.org/zap"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"

	"github.com/Maple-Open-Tech/monorepo/cloud/backend/config"
	dom_financial "github.com/Maple-Open-Tech/monorepo/cloud/backend/internal/ipe/domain/financialanalysis"
)

type financialAnalysisImpl struct {
	Logger     *zap.Logger
	DbClient   *mongo.Client
	Collection *mongo.Collection
}

func NewRepository(appCfg *config.Configuration, loggerp *zap.Logger, client *mongo.Client) dom_financial.FinancialRepository {
	// Get collection
	collection := client.Database(appCfg.DB.IncomePropertyEvaluatorName).Collection("financial_analyses")

	// Create indexes
	_, err := collection.Indexes().CreateMany(context.TODO(), []mongo.IndexModel{
		{Keys: bson.D{{Key: "propertyId", Value: 1}}, Options: options.Index().SetUnique(true)},
		{Keys: bson.D{{Key: "purchasePrice", Value: 1}}},
		{Keys: bson.D{{Key: "annualGrossIncome", Value: 1}}},
		{Keys: bson.D{{Key: "annualNetIncome", Value: 1}}},
		{Keys: bson.D{{Key: "capRateWithoutMortgage", Value: 1}}},
	})
	if err != nil {
		// Fatal error on startup to meet requirements of google/wire framework
		log.Fatal(err)
	}

	return &financialAnalysisImpl{
		Logger:     loggerp,
		DbClient:   client,
		Collection: collection,
	}
}
