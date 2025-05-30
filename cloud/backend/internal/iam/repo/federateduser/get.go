// github.com/Maple-Open-Tech/monorepo/cloud/backend/internal/iam/repo/federateduser/get.go
package federateduser

import (
	"context"

	"go.uber.org/zap"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"

	dom_user "github.com/Maple-Open-Tech/monorepo/cloud/backend/internal/iam/domain/federateduser"
)

func (impl userStorerImpl) GetByID(ctx context.Context, id primitive.ObjectID) (*dom_user.FederatedUser, error) {
	filter := bson.M{"_id": id}

	var result dom_user.FederatedUser
	err := impl.Collection.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			// This error means your query did not match any documents.
			return nil, nil
		}
		impl.Logger.Error("database get by base user id error", zap.Any("error", err))
		return nil, err
	}
	return &result, nil
}

func (impl userStorerImpl) GetByEmail(ctx context.Context, email string) (*dom_user.FederatedUser, error) {
	filter := bson.M{"email": email}

	var result dom_user.FederatedUser
	err := impl.Collection.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			// This error means your query did not match any documents.
			return nil, nil
		}
		impl.Logger.Error("database get by email error", zap.Any("error", err))
		return nil, err
	}
	return &result, nil
}

func (impl userStorerImpl) GetByVerificationCode(ctx context.Context, verificationCode string) (*dom_user.FederatedUser, error) {
	filter := bson.M{"email_verification_code": verificationCode}

	var result dom_user.FederatedUser
	err := impl.Collection.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			// This error means your query did not match any documents.
			return nil, nil
		}
		impl.Logger.Error("database get by verification code error", zap.Any("error", err))
		return nil, err
	}
	return &result, nil
}
