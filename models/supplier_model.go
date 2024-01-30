package models

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"log"
	"os"
	"time"
)

type Supplier struct {
	Id      string `json:"id" bson:"id"`
	Name    string `json:"name" bson:"name"`
	Country string `json:"country" bson:"country"`
	City    string `json:"city" bson:"city"`
}

type SupplierModel struct {
	COLLECTION *mongo.Collection
}

// SupplierModel methods
func (s *SupplierModel) Add(supplier Supplier) error {
	company := os.Getenv("COMPANY")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	caser := cases.Title(language.English)
	companyFirstLMaiusc := caser.String(company)

	var result bson.M
	err := s.COLLECTION.FindOne(ctx, bson.D{{"name", companyFirstLMaiusc}}).Decode(&result)

	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil
		}
		return err
	}

	filter := bson.D{{"name", companyFirstLMaiusc}}
	update := bson.D{{"$push", bson.D{{"suppliers", supplier}}}}

	_, err = s.COLLECTION.UpdateOne(ctx, filter, update)
	if err != nil {
		log.Println("Failed to insert supplier: ", err)
		return err
	}
	return nil

}

func (s *SupplierModel) Update(supplier Supplier) error {
	company := os.Getenv("COMPANY")
	caser := cases.Title(language.English)
	companyFirstLMaiusc := caser.String(company)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{"name": companyFirstLMaiusc, "suppliers.id": supplier.Id}
	update := bson.D{{"$set", bson.D{{"suppliers.$", supplier}}}}

	res, err := s.COLLECTION.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	if res.MatchedCount != 0 {
		fmt.Printf("matched and replaced supplier %v", supplier.Id)
		return nil
	}

	return errors.New("supplier not found")
}

func (s *SupplierModel) GetAll() ([]Supplier, error) {
	company := os.Getenv("COMPANY")
	caser := cases.Title(language.English)
	companyFirstLMaiusc := caser.String(company)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var result bson.M

	err := s.COLLECTION.FindOne(ctx, bson.D{{"name", companyFirstLMaiusc}}).Decode(&result)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}

	suppliersRaw, ok := result["suppliers"]
	if !ok {
		return nil, errors.New("field 'suppliers' not found")
	}

	suppliersJSON, err := json.Marshal(suppliersRaw)
	if err != nil {
		return nil, err
	}

	var suppliers []Supplier
	err = json.Unmarshal(suppliersJSON, &suppliers)
	if err != nil {
		return nil, err
	}

	return suppliers, nil
}

func (s *SupplierModel) GetOne(id string) (*Supplier, error) {
	company := os.Getenv("COMPANY")
	caser := cases.Title(language.English)
	companyFirstLMaiusc := caser.String(company)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	pipeline := bson.A{
		bson.D{
			{"$match", bson.D{{"name", companyFirstLMaiusc}}},
		},
		bson.D{
			{"$unwind", "$suppliers"},
		},
		bson.D{
			{"$match", bson.D{{"suppliers.id", id}}},
		},
	}

	cursor, err := s.COLLECTION.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}

	var results []bson.M
	if err := cursor.All(ctx, &results); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			fmt.Println("Object not found")
			return nil, nil
		} else {
			log.Fatal(err)
			return nil, err
		}
	}

	if len(results) > 0 {
		supplierBytes, err := json.Marshal(results[0]["suppliers"])
		if err != nil {
			return nil, err
		}

		var mySupplier Supplier
		err = json.Unmarshal(supplierBytes, &mySupplier)
		if err != nil {
			return nil, err
		}

		return &mySupplier, nil
	}

	return nil, fmt.Errorf("supplier with ID %v not found", id)
}

func (s *SupplierModel) DeleteOne(id string) error {
	company := os.Getenv("COMPANY")
	caser := cases.Title(language.English)
	companyFirstLMaiusc := caser.String(company)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	supplier, err := s.GetOne(id)
	if err != nil {
		return err
	}

	if supplier != nil {
		filter := bson.D{{"name", companyFirstLMaiusc}}
		update := bson.D{
			{"$pull", bson.D{
				{"suppliers", bson.D{{"id", id}}},
			}},
		}

		res, err := s.COLLECTION.UpdateOne(ctx, filter, update)
		if err != nil {
			return err
		}
		if res.MatchedCount != 0 {
			fmt.Printf("matched and deleted supplier %v", id)
			return nil
		}
	}

	return errors.New("something went wrong")
}
