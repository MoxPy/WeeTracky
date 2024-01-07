package models

import (
	"context"
	"encoding/json"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"log"
	"os"
	"time"
)

type Cert struct {
	Id      string `json:"id" bson:"id"`
	Name    string `json:"name" bson:"name"`
	Issuer  string `json:"issuer" bson:"issuer"`
	Details string `json:"details" bson:"details" default:"No details"`
}

type CertModel struct {
	DB         *mongo.Client
	COLLECTION *mongo.Collection
}

// CertModel methods
func (c *CertModel) Add(cert Cert) error {
	company := os.Getenv("COMPANY")
	caser := cases.Title(language.English)
	companyFirstLMaiusc := caser.String(company)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var result bson.M
	err := c.COLLECTION.FindOne(ctx, bson.D{{"name", companyFirstLMaiusc}}).Decode(&result)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil
		}
		return nil
	}

	filter := bson.D{{"name", companyFirstLMaiusc}}
	update := bson.D{{"$push", bson.D{{"certs", cert}}}}

	_, err = c.COLLECTION.UpdateOne(ctx, filter, update)
	if err != nil {
		log.Println("Failed to insert certification: ", err)
		return err
	}
	return nil
}

func (c *CertModel) GetAll() ([]Cert, error) {
	company := os.Getenv("COMPANY")
	caser := cases.Title(language.English)
	companyFirstLMaiusc := caser.String(company)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var result bson.M

	err := c.COLLECTION.FindOne(ctx, bson.D{{"name", companyFirstLMaiusc}}).Decode(&result)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}

	certsRaw, ok := result["certs"]
	if !ok {
		return nil, errors.New("field 'certs' not found")
	}

	certsJSON, err := json.Marshal(certsRaw)
	if err != nil {
		return nil, err
	}

	var certs []Cert
	err = json.Unmarshal(certsJSON, &certs)
	if err != nil {
		return nil, err
	}

	return certs, nil
}
