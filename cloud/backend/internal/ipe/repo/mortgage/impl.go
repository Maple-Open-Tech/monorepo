// cloud/backend/internal/ipe/repo/mortgage/impl.go
package mortgage

import (
	"context"
	"log"

	"go.uber.org/zap"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"

	"github.com/Maple-Open-Tech/monorepo/cloud/backend/config"
	dom_mortgage "github.com/Maple-Open-Tech/monorepo/cloud/backend/internal/ipe/domain/mortgage"
)

type mortgageImpl struct {
	Logger     *zap.Logger
	DbClient   *mongo.Client
	Collection *mongo.Collection
}

func NewRepository(appCfg *config.Configuration, loggerp *zap.Logger, client *mongo.Client) dom_mortgage.MortgageRepository {
	// Get collection
	collection := client.Database(appCfg.DB.IncomePropertyEvaluatorName).Collection("mortgages")

	// Create indexes
	_, err := collection.Indexes().CreateMany(context.TODO(), []mongo.IndexModel{
		{Keys: bson.D{{Key: "financialAnalysisId", Value: 1}}, Options: options.Index().SetUnique(true)},
		{Keys: bson.D{{Key: "loanAmount", Value: 1}}},
		{Keys: bson.D{{Key: "annualInterestRate", Value: 1}}},
		{Keys: bson.D{{Key: "amortizationYear", Value: 1}}},
	})
	if err != nil {
		// Fatal error on startup to meet requirements of google/wire framework
		log.Fatal(err)
	}

	return &mortgageImpl{
		Logger:     loggerp,
		DbClient:   client,
		Collection: collection,
	}
}
