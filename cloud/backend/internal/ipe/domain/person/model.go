package person

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Common fields for all person types
type basePerson struct {
	PersonName      string `bson:"personName" json:"personName"`
	Address         string `bson:"address" json:"address"`
	City            string `bson:"city" json:"city"`
	Province        string `bson:"province" json:"province"`
	Country         string `bson:"country" json:"country"`
	PostalCode      string `bson:"postalCode" json:"postalCode"`
	Email           string `bson:"email" json:"email"`
	OfficeTelNumber string `bson:"officeTelNumber" json:"officeTelNumber"`
	MobileTelNumber string `bson:"mobileTelNumber" json:"mobileTelNumber"`
	FaxTelNumber    string `bson:"faxTelNumber" json:"faxTelNumber"`
	Website         string `bson:"website" json:"website"`
	RecordUniqueID  string `bson:"recordUniqueId" json:"recordUniqueId"`
	LogoPhotoData   []byte `bson:"logoPhotoData,omitempty" json:"logoPhotoData,omitempty"`
}

// Client represents a client
type Client struct {
	ID primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	basePerson
}

// Presenter represents a presenter
type Presenter struct {
	ID primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	basePerson
	RecordUnquieID string `bson:"recordUnquieId" json:"recordUnquieId"`
}

// Owner represents a property owner
type Owner struct {
	ID primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	basePerson
}

// Compare represents a property comparison
type Compare struct {
	ID                       primitive.ObjectID   `bson:"_id,omitempty" json:"id,omitempty"`
	RecordUniqueID           string               `bson:"recordUniqueId" json:"recordUniqueId"`
	ClientID                 primitive.ObjectID   `bson:"clientId" json:"clientId"`
	PresenterID              primitive.ObjectID   `bson:"presenterId" json:"presenterId"`
	SelectedIncomeProperties []primitive.ObjectID `bson:"selectedIncomeProperties" json:"selectedIncomeProperties"`
}
