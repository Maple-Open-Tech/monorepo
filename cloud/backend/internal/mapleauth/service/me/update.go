// github.com/Maple-Open-Tech/monorepo/cloud/backend/internal/mapleauth/service/baseuser/service.go
package me

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"go.uber.org/zap"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/v2/mongo"

	"github.com/Maple-Open-Tech/monorepo/cloud/backend/config"
	"github.com/Maple-Open-Tech/monorepo/cloud/backend/config/constants"
	uc_user "github.com/Maple-Open-Tech/monorepo/cloud/backend/internal/mapleauth/usecase/baseuser"
	"github.com/Maple-Open-Tech/monorepo/cloud/backend/pkg/httperror"
)

type UpdateMeRequestDTO struct {
	Email                                          string `bson:"email" json:"email"`
	FirstName                                      string `bson:"first_name" json:"first_name"`
	LastName                                       string `bson:"last_name" json:"last_name"`
	Phone                                          string `bson:"phone" json:"phone,omitempty"`
	Country                                        string `bson:"country" json:"country,omitempty"`
	Region                                         string `bson:"region" json:"region,omitempty"`
	Timezone                                       string `bson:"timezone" json:"timezone"`
	AgreePromotions                                bool   `bson:"agree_promotions" json:"agree_promotions,omitempty"`
	AgreeToTrackingAcrossThirdPartyAppsAndServices bool   `bson:"agree_to_tracking_across_third_party_apps_and_services" json:"agree_to_tracking_across_third_party_apps_and_services,omitempty"`
}

type UpdateMeService interface {
	Execute(sessCtx context.Context, req *UpdateMeRequestDTO) (*MeResponseDTO, error)
}

type updateMeServiceImpl struct {
	config                *config.Configuration
	logger                *zap.Logger
	userGetByIDUseCase    uc_user.UserGetByIDUseCase
	userGetByEmailUseCase uc_user.UserGetByEmailUseCase
	userUpdateUseCase     uc_user.UserUpdateUseCase
}

func NewUpdateMeService(
	config *config.Configuration,
	logger *zap.Logger,
	userGetByIDUseCase uc_user.UserGetByIDUseCase,
	userGetByEmailUseCase uc_user.UserGetByEmailUseCase,
	userUpdateUseCase uc_user.UserUpdateUseCase,
) UpdateMeService {
	return &updateMeServiceImpl{
		config:                config,
		logger:                logger,
		userGetByIDUseCase:    userGetByIDUseCase,
		userGetByEmailUseCase: userGetByEmailUseCase,
		userUpdateUseCase:     userUpdateUseCase,
	}
}

func (svc *updateMeServiceImpl) Execute(sessCtx context.Context, req *UpdateMeRequestDTO) (*MeResponseDTO, error) {
	//
	// Get required from context.
	//

	userID, ok := sessCtx.Value(constants.SessionUserID).(primitive.ObjectID)
	if !ok {
		svc.logger.Error("Failed getting local baseuser id",
			zap.Any("error", "Not found in context: user_id"))
		return nil, errors.New("baseuser id not found in context")
	}

	//
	// STEP 2: Validation
	//

	if req == nil {
		svc.logger.Warn("Failed validation with nothing received")
		return nil, httperror.NewForBadRequestWithSingleField("non_field_error", "Request is required in submission")
	}

	// Sanitization
	req.Email = strings.ToLower(req.Email) // Ensure email is lowercase

	e := make(map[string]string)
	// Add any specific field validations here if needed. Example:
	if req.FirstName == "" {
		e["first_name"] = "First name is required"
	}
	if req.LastName == "" {
		e["last_name"] = "Last name is required"
	}
	if req.Email == "" {
		e["email"] = "Email is required"
	}
	if len(req.Email) > 255 {
		e["email"] = "Email is too long"
	}
	if req.Phone == "" {
		e["phone"] = "Phone confirm is required"
	}
	if req.Country == "" {
		e["country"] = "Country is required"
	}
	if req.Timezone == "" {
		e["timezone"] = "Password confirm is required"
	}
	if len(e) != 0 {
		svc.logger.Warn("Failed validation",
			zap.Any("error", e))
		return nil, httperror.NewForBadRequest(&e)
	}

	//
	// Get related records.
	//

	// Get the baseuser account (aka "Me").
	baseuser, err := svc.userGetByIDUseCase.Execute(sessCtx, userID)
	if err != nil {
		// If it's a "not found" error, it's a critical issue since the ID came from the context.
		if errors.Is(err, mongo.ErrNoDocuments) {
			err := fmt.Errorf("authenticated baseuser does not exist for id: %v", userID.Hex())
			svc.logger.Error("Failed getting authenticated baseuser", zap.Any("error", err))
			return nil, err
		}
		// Handle other potential errors during fetch.
		svc.logger.Error("Failed getting baseuser by ID", zap.Any("error", err))
		return nil, err
	}
	// Defensive check, though GetByID should return ErrNoDocuments if not found.
	if baseuser == nil {
		err := fmt.Errorf("baseuser is nil after lookup for id: %v", userID.Hex())
		svc.logger.Error("Failed getting baseuser", zap.Any("error", err))
		return nil, err
	}

	//
	// Check if the requested email is already taken by another baseuser.
	//
	if req.Email != baseuser.Email {
		existingUser, err := svc.userGetByEmailUseCase.Execute(sessCtx, req.Email)
		if err != nil && !errors.Is(err, mongo.ErrNoDocuments) {
			svc.logger.Error("Failed checking existing email", zap.String("email", req.Email), zap.Any("error", err))
			return nil, err // Internal Server Error
		}
		if existingUser != nil {
			// Email exists. Check if it belongs to the *current* baseuser (which shouldn't happen based on the outer if, but defensive check).
			// The important check is implicit: if existingUser is not nil, the email is taken.
			// We already know req.Email != baseuser.Email, so if existingUser is found, it *must* be another baseuser.
			svc.logger.Warn("Attempted to update to an email already in use",
				zap.String("user_id", userID.Hex()),
				zap.String("existing_user_id", existingUser.ID.Hex()),
				zap.String("email", req.Email))
			e["email"] = "This email address is already in use."
			return nil, httperror.NewForBadRequest(&e)
		}
		// If err is mongo.ErrNoDocuments or existingUser is nil, the email is available.
	}

	//
	// Update local database.
	//

	// Apply changes from request DTO to the baseuser object
	baseuser.Email = req.Email
	baseuser.FirstName = req.FirstName
	baseuser.LastName = req.LastName
	baseuser.Name = fmt.Sprintf("%s %s", req.FirstName, req.LastName)
	baseuser.LexicalName = fmt.Sprintf("%s, %s", req.LastName, req.FirstName)
	baseuser.Phone = req.Phone
	baseuser.Country = req.Country
	baseuser.Region = req.Region
	baseuser.Timezone = req.Timezone
	baseuser.AgreePromotions = req.AgreePromotions
	baseuser.AgreeToTrackingAcrossThirdPartyAppsAndServices = req.AgreeToTrackingAcrossThirdPartyAppsAndServices

	// Persist changes
	if err := svc.userUpdateUseCase.Execute(sessCtx, baseuser); err != nil {
		svc.logger.Error("Failed updating baseuser", zap.Any("error", err), zap.String("user_id", baseuser.ID.Hex()))
		// Consider mapping specific DB errors (like constraint violations) to HTTP errors if applicable
		return nil, err
	}

	svc.logger.Debug("BaseUser updated successfully",
		zap.String("user_id", baseuser.ID.Hex()))

	// Return updated baseuser details
	return &MeResponseDTO{
		ID:              baseuser.ID,
		Email:           baseuser.Email,
		FirstName:       baseuser.FirstName,
		LastName:        baseuser.LastName,
		Name:            baseuser.Name,
		LexicalName:     baseuser.LexicalName,
		Phone:           baseuser.Phone,
		Country:         baseuser.Country,
		Region:          baseuser.Region, // Added Region
		Timezone:        baseuser.Timezone,
		AgreePromotions: baseuser.AgreePromotions,
		AgreeToTrackingAcrossThirdPartyAppsAndServices: baseuser.AgreeToTrackingAcrossThirdPartyAppsAndServices,
	}, nil
}
