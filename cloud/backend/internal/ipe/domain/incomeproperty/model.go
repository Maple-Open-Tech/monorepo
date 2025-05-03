package incomeproperty

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type IncomeProperty struct {
	ID                 primitive.ObjectID `bson:"_id,omitempty"`
	Address            string             `bson:"address"`
	City               string             `bson:"city"`
	Province           string             `bson:"province"`
	Country            string             `bson:"country"`
	PropertyCode       string             `bson:"propertyCode"`
	RecordName         string             `bson:"recordName"`
	RecordCreationDate time.Time          `bson:"recordCreationDate"`
	MainPhotoThumbnail []byte             `bson:"mainPhotoThumbnail,omitempty"`
}
