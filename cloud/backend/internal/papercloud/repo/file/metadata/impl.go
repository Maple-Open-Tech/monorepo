// cloud/backend/internal/papercloud/repo/file/metadata/impl.go
package metadata

import (
	"context"
	"errors"
	"time"

	"go.uber.org/zap"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"

	"github.com/Maple-Open-Tech/monorepo/cloud/backend/config"
	dom_file "github.com/Maple-Open-Tech/monorepo/cloud/backend/internal/papercloud/domain/file"
)

type FileMetadataRepository interface {
	Create(file *dom_file.File) error
	Get(id string) (*dom_file.File, error)
	GetByFileID(fileID string) (*dom_file.File, error)
	GetByCollection(collectionID string) ([]*dom_file.File, error)
	Update(file *dom_file.File) error
	Delete(id string) error
	CheckIfExistsByID(id string) (bool, error)
	CheckIfUserHasAccess(fileID string, userID string) (bool, error)
}

type fileMetadataRepositoryImpl struct {
	Logger     *zap.Logger
	DbClient   *mongo.Client
	Collection *mongo.Collection
}

func NewRepository(appCfg *config.Configuration, loggerp *zap.Logger, client *mongo.Client) FileMetadataRepository {
	// Initialize collection in the PaperCloud database
	fc := client.Database(appCfg.DB.PaperCloudName).Collection("files")

	// Reset indexes for development purposes
	if err := fc.Indexes().DropAll(context.TODO()); err != nil {
		loggerp.Warn("failed deleting all indexes",
			zap.Any("err", err))
	}

	// Create indexes for efficient queries
	_, err := fc.Indexes().CreateMany(context.TODO(), []mongo.IndexModel{
		{Keys: bson.D{
			{Key: "created_at", Value: -1},
		}},
		{Keys: bson.D{
			{Key: "id", Value: 1},
		}, Options: options.Index().SetUnique(true)},
		{Keys: bson.D{
			{Key: "file_id", Value: 1},
		}, Options: options.Index().SetUnique(true)},
		{Keys: bson.D{
			{Key: "collection_id", Value: 1},
			{Key: "created_at", Value: -1},
		}},
		{Keys: bson.D{
			{Key: "owner_id", Value: 1},
			{Key: "created_at", Value: -1},
		}},
	})

	if err != nil {
		loggerp.Error("failed creating indexes error", zap.Any("err", err))
		return nil
	}

	return &fileMetadataRepositoryImpl{
		Logger:     loggerp,
		DbClient:   client,
		Collection: fc,
	}
}

// Create a new file metadata entry
func (impl fileMetadataRepositoryImpl) Create(file *dom_file.File) error {
	ctx := context.Background()

	// Validate file ID
	if file.ID == "" {
		file.ID = uuid.New().String()
	}

	// Set creation time if not set
	if file.CreatedAt.IsZero() {
		file.CreatedAt = time.Now()
	}

	// Set modification time to creation time
	file.ModifiedAt = file.CreatedAt

	// Insert file document
	_, err := impl.Collection.InsertOne(ctx, file)
	if err != nil {
		impl.Logger.Error("database failed create file error",
			zap.Any("error", err),
			zap.String("id", file.ID))
		return err
	}

	return nil
}

// Get file by ID
func (impl fileMetadataRepositoryImpl) Get(id string) (*dom_file.File, error) {
	ctx := context.Background()
	filter := bson.M{"id": id}

	var result dom_file.File
	err := impl.Collection.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		impl.Logger.Error("database get by file id error", zap.Any("error", err))
		return nil, err
	}
	return &result, nil
}

// GetByFileID gets a file by its client-generated FileID
func (impl fileMetadataRepositoryImpl) GetByFileID(fileID string) (*dom_file.File, error) {
	ctx := context.Background()
	filter := bson.M{"file_id": fileID}

	var result dom_file.File
	err := impl.Collection.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		impl.Logger.Error("database get by file_id error", zap.Any("error", err))
		return nil, err
	}
	return &result, nil
}

// GetByCollection gets all files in a collection
func (impl fileMetadataRepositoryImpl) GetByCollection(collectionID string) ([]*dom_file.File, error) {
	ctx := context.Background()
	filter := bson.M{"collection_id": collectionID}

	cursor, err := impl.Collection.Find(ctx, filter)
	if err != nil {
		impl.Logger.Error("database get files by collection id error", zap.Any("error", err))
		return nil, err
	}
	defer cursor.Close(ctx)

	var files []*dom_file.File
	if err = cursor.All(ctx, &files); err != nil {
		impl.Logger.Error("database decode files error", zap.Any("error", err))
		return nil, err
	}

	return files, nil
}

// Update a file's metadata
func (impl fileMetadataRepositoryImpl) Update(file *dom_file.File) error {
	ctx := context.Background()
	filter := bson.M{"id": file.ID}

	// Update modification time
	file.ModifiedAt = time.Now()

	update := bson.M{
		"$set": file,
	}

	_, err := impl.Collection.UpdateOne(ctx, filter, update)
	if err != nil {
		impl.Logger.Error("database update file error",
			zap.Any("error", err),
			zap.String("id", file.ID))
		return err
	}

	return nil
}

// Delete a file's metadata
func (impl fileMetadataRepositoryImpl) Delete(id string) error {
	ctx := context.Background()
	filter := bson.M{"id": id}

	_, err := impl.Collection.DeleteOne(ctx, filter)
	if err != nil {
		impl.Logger.Error("database failed deletion error",
			zap.Any("error", err),
			zap.String("id", id))
		return err
	}

	return nil
}

// CheckIfExistsByID checks if a file exists by ID
func (impl fileMetadataRepositoryImpl) CheckIfExistsByID(id string) (bool, error) {
	ctx := context.Background()
	filter := bson.M{"id": id}

	count, err := impl.Collection.CountDocuments(ctx, filter)
	if err != nil {
		impl.Logger.Error("database check if exists by ID error", zap.Any("error", err))
		return false, err
	}
	return count >= 1, nil
}

// CheckIfUserHasAccess checks if a user has access to a file
func (impl fileMetadataRepositoryImpl) CheckIfUserHasAccess(fileID string, userID string) (bool, error) {
	// ctx := context.Background()

	// First get the file to find its owner and collection ID
	file, err := impl.Get(fileID)
	if err != nil {
		impl.Logger.Error("database get file error", zap.Any("error", err))
		return false, err
	}
	if file == nil {
		return false, nil
	}

	// Direct access if user is the owner
	if file.OwnerID == userID {
		return true, nil
	}

	// We need to check collection sharing permissions
	// This would require a separate call to the collection repository
	// For now, we'll just return false, but in a real implementation
	// we would check if the collection is shared with the user

	// Note: In a complete implementation, we would inject the collection repository
	// and check if the collection is shared with the user

	return false, errors.New("collection access check not implemented")
}
