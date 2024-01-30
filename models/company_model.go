package models

import (
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"time"
)

type Company struct {
	ID        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Name      string             `json:"name" bson:"name"`
	Products  []Product          `json:"products" bson:"products"`
	Materials []Material         `json:"materials" bson:"materials"`
	Suppliers []Supplier         `json:"suppliers" bson:"suppliers"`
	Certs     []Cert             `json:"certs" bson:"certs"`
}

type CompanyModel struct {
	COLLECTION *mongo.Collection
}

func (c *CompanyModel) Initialize(company Company) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	newCompany, err := c.COLLECTION.InsertOne(ctx, company)

	if err != nil {
		log.Fatal(err)
		return err
	}

	log.Printf("Created new company with ID %v\n", newCompany.InsertedID)
	return nil
}
