// github.com/Maple-Open-Tech/monorepo/cloud/backend/internal/papercloud/repo/user/update.go
package user

import (
	"context"

	"go.uber.org/zap"

	"go.mongodb.org/mongo-driver/v2/bson"

	dom_user "github.com/Maple-Open-Tech/monorepo/cloud/backend/internal/papercloud/domain/user"
)

// func (impl userStorerImpl) DeleteByID(ctx context.Context, id primitive.ObjectID) error {
// 	_, err := impl.Collection.DeleteOne(ctx, bson.M{"_id": id})
// 	if err != nil {
// 		impl.Logger.Error("database failed deletion error",
// 			zap.Any("error", err))
// 		return err
// 	}
// 	return nil
// }
//
// func (impl userStorerImpl) CheckIfExistsByID(ctx context.Context, id primitive.ObjectID) (bool, error) {
// 	filter := bson.M{"_id": id}
// 	count, err := impl.Collection.CountDocuments(ctx, filter)
// 	if err != nil {
// 		impl.Logger.Error("database check if exists by ID error", zap.Any("error", err))
// 		return false, err
// 	}
// 	return count >= 1, nil
// }
//
// func (impl userStorerImpl) CheckIfExistsByEmail(ctx context.Context, email string) (bool, error) {
// 	filter := bson.M{"email": email}
// 	count, err := impl.Collection.CountDocuments(ctx, filter)
// 	if err != nil {
// 		impl.Logger.Error("database check if exists by email error", zap.Any("error", err))
// 		return false, err
// 	}
// 	return count >= 1, nil
// }

func (impl userStorerImpl) UpdateByID(ctx context.Context, m *dom_user.User) error {
	filter := bson.M{"_id": m.ID}

	update := bson.M{ // DEVELOPERS NOTE: https://stackoverflow.com/a/60946010
		"$set": m,
	}

	// execute the UpdateOne() function to update the first matching document
	_, err := impl.Collection.UpdateOne(ctx, filter, update)
	if err != nil {
		impl.Logger.Error("database update by id error", zap.Any("error", err))
		return err
	}

	// // display the number of documents updated
	// impl.Logger.Debug("number of documents updated", zap.Int64("modified_count", result.ModifiedCount))

	return nil
}

// // buildCountMatchStage creates the match stage for the aggregation pipeline
// func (s *userStorerImpl) buildCountMatchStage(filter *dom_user.UserFilter) bson.M {
// 	match := bson.M{}
//
// 	if filter.Status != 0 {
// 		match["status"] = filter.Status
// 	}
//
// 	// Date range filtering
// 	if filter.CreatedAtStart != nil || filter.CreatedAtEnd != nil {
// 		createdAtFilter := bson.M{}
// 		if filter.CreatedAtStart != nil {
// 			createdAtFilter["$gte"] = filter.CreatedAtStart
// 		}
// 		if filter.CreatedAtEnd != nil {
// 			createdAtFilter["$lte"] = filter.CreatedAtEnd
// 		}
// 		match["created_at"] = createdAtFilter
// 	}
//
// 	// Text search for name
// 	if filter.Name != nil && *filter.Name != "" {
// 		match["$text"] = bson.M{"$search": *filter.Name}
// 	}
//
// 	return match
// }
//
// // hasActiveFilters checks if any filters besides tenant_id are active
// func (s *userStorerImpl) hasActiveFilters(filter *dom_user.UserFilter) bool {
// 	return filter.Name != nil ||
// 		filter.Status != 0 ||
// 		filter.CreatedAtStart != nil ||
// 		filter.CreatedAtEnd != nil
// }
//
// func (s *userStorerImpl) CountByFilter(ctx context.Context, filter *dom_user.UserFilter) (uint64, error) {
// 	if filter == nil {
// 		return 0, errors.New("filter cannot be nil")
// 	}
//
// 	// For exact counts with filters
// 	if s.hasActiveFilters(filter) {
// 		pipeline := []bson.D{
// 			{{Key: "$match", Value: s.buildCountMatchStage(filter)}},
// 			{{Key: "$count", Value: "total"}},
// 		}
//
// 		// Use aggregation with allowDiskUse for large datasets
// 		opts := options.Aggregate().SetAllowDiskUse(true)
// 		cursor, err := s.Collection.Aggregate(ctx, pipeline, opts)
// 		if err != nil {
// 			return 0, err
// 		}
// 		defer cursor.Close(ctx)
//
// 		// Decode the result
// 		var results []struct {
// 			Total int64 `bson:"total"`
// 		}
// 		if err := cursor.All(ctx, &results); err != nil {
// 			return 0, err
// 		}
//
// 		if len(results) == 0 {
// 			return 0, nil
// 		}
//
// 		return uint64(results[0].Total), nil
// 	}
//
// 	// For unfiltered counts (much faster for basic tenant-only filtering)
// 	countOpts := options.Count().SetHint("tenant_id_1_created_at_-1")
// 	match := bson.M{}
//
// 	count, err := s.Collection.CountDocuments(ctx, match, countOpts)
// 	if err != nil {
// 		return 0, err
// 	}
//
// 	return uint64(count), nil
// }
//
// func (impl userStorerImpl) buildMatchStage(filter *dom_user.UserFilter) bson.M {
// 	match := bson.M{}
//
// 	// Handle cursor-based pagination
// 	if filter.LastID != nil && filter.LastCreatedAt != nil {
// 		match["$or"] = []bson.M{
// 			{
// 				"created_at": bson.M{"$lt": filter.LastCreatedAt},
// 			},
// 			{
// 				"created_at": filter.LastCreatedAt,
// 				"_id":        bson.M{"$lt": filter.LastID},
// 			},
// 		}
// 	}
//
// 	// Add other filters
// 	if filter.Status != 0 {
// 		match["status"] = filter.Status
// 	}
//
// 	if filter.ProfileVerificationStatus != 0 {
// 		match["profile_verification_status"] = filter.ProfileVerificationStatus
// 	}
//
// 	// if !filter.UserID.IsZero() {
// 	// 	match["user_id"] = filter.UserID
// 	// }
//
// 	if filter.CreatedAtStart != nil || filter.CreatedAtEnd != nil {
// 		createdAtFilter := bson.M{}
// 		if filter.CreatedAtStart != nil {
// 			createdAtFilter["$gte"] = filter.CreatedAtStart
// 		}
// 		if filter.CreatedAtEnd != nil {
// 			createdAtFilter["$lte"] = filter.CreatedAtEnd
// 		}
// 		match["created_at"] = createdAtFilter
// 	}
//
// 	// Text search for name
// 	if filter.Name != nil && *filter.Name != "" {
// 		match["$text"] = bson.M{"$search": *filter.Name}
// 	}
//
// 	return match
// }
//
// func (impl userStorerImpl) ListByFilter(ctx context.Context, filter *dom_user.UserFilter) (*dom_user.UserFilterResult, error) {
// 	if filter == nil {
// 		return nil, errors.New("filter cannot be nil")
// 	}
//
// 	// Default limit if not specified
// 	if filter.Limit <= 0 {
// 		filter.Limit = 100
// 	}
//
// 	// Request one more document than needed to determine if there are more results
// 	limit := filter.Limit + 1
//
// 	// Build the aggregation pipeline
// 	pipeline := make([]bson.D, 0)
//
// 	// Match stage - initial filtering
// 	matchStage := bson.D{{"$match", impl.buildMatchStage(filter)}}
// 	pipeline = append(pipeline, matchStage)
//
// 	// Sort stage
// 	sortStage := bson.D{{"$sort", bson.D{
// 		{"created_at", -1},
// 		{"_id", -1},
// 	}}}
// 	pipeline = append(pipeline, sortStage)
//
// 	// Limit stage
// 	limitStage := bson.D{{"$limit", limit}}
// 	pipeline = append(pipeline, limitStage)
//
// 	// Execute aggregation
// 	opts := options.Aggregate().SetAllowDiskUse(true)
// 	cursor, err := impl.Collection.Aggregate(ctx, pipeline, opts)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer cursor.Close(ctx)
//
// 	// Decode results
// 	var Users []*dom_user.User
// 	if err := cursor.All(ctx, &Users); err != nil {
// 		return nil, err
// 	}
//
// 	// Handle empty results case
// 	if len(Users) == 0 {
// 		// For debugging purposes only.
// 		// impl.Logger.Debug("Empty list", zap.Any("filter", filter))
// 		return &dom_user.UserFilterResult{
// 			Users:   make([]*dom_user.User, 0),
// 			HasMore: false,
// 		}, nil
// 	}
//
// 	// Check if there are more results
// 	hasMore := false
// 	if len(Users) > int(filter.Limit) {
// 		hasMore = true
// 		Users = Users[:len(Users)-1]
// 	}
//
// 	// Get last document info for next page
// 	lastDoc := Users[len(Users)-1]
//
// 	return &dom_user.UserFilterResult{
// 		Users:         Users,
// 		HasMore:       hasMore,
// 		LastID:        lastDoc.ID,
// 		LastCreatedAt: lastDoc.CreatedAt,
// 	}, nil
// }
