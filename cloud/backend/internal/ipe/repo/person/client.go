// cloud/backend/internal/ipe/repo/person/client.go
package person

import (
	"context"

	"go.uber.org/zap"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"

	dom_person "github.com/Maple-Open-Tech/monorepo/cloud/backend/internal/ipe/domain/person"
)

// FindClientByID retrieves a client by ID
func (impl personImpl) FindClientByID(ctx context.Context, id primitive.ObjectID) (*dom_person.Client, error) {
	filter := bson.M{"_id": id}

	var result dom_person.Client
	err := impl.ClientCollection.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		impl.Logger.Error("database get client by id error", zap.Any("error", err))
		return nil, err
	}
	return &result, nil
}

// FindAllClients retrieves all clients
func (impl personImpl) FindAllClients(ctx context.Context) ([]*dom_person.Client, error) {
	cursor, err := impl.ClientCollection.Find(ctx, bson.M{})
	if err != nil {
		impl.Logger.Error("database find all clients error", zap.Any("error", err))
		return nil, err
	}
	defer cursor.Close(ctx)

	var clients []*dom_person.Client
	if err = cursor.All(ctx, &clients); err != nil {
		impl.Logger.Error("database decode clients error", zap.Any("error", err))
		return nil, err
	}

	return clients, nil
}

// SaveClient creates a new client
func (impl personImpl) SaveClient(ctx context.Context, client *dom_person.Client) (primitive.ObjectID, error) {
	if client.ID == primitive.NilObjectID {
		client.ID = primitive.NewObjectID()
		impl.Logger.Warn("database insert client without ID, created ID now", zap.Any("id", client.ID))
	}

	_, err := impl.ClientCollection.InsertOne(ctx, client)
	if err != nil {
		impl.Logger.Error("database save client error", zap.Any("error", err))
		return primitive.NilObjectID, err
	}

	return client.ID, nil
}

// UpdateClient updates an existing client
func (impl personImpl) UpdateClient(ctx context.Context, client *dom_person.Client) error {
	filter := bson.M{"_id": client.ID}
	update := bson.M{"$set": client}

	_, err := impl.ClientCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		impl.Logger.Error("database update client error", zap.Any("error", err))
		return err
	}

	return nil
}

// DeleteClient removes a client by ID
func (impl personImpl) DeleteClient(ctx context.Context, id primitive.ObjectID) error {
	filter := bson.M{"_id": id}

	_, err := impl.ClientCollection.DeleteOne(ctx, filter)
	if err != nil {
		impl.Logger.Error("database delete client error", zap.Any("error", err))
		return err
	}

	return nil
}
