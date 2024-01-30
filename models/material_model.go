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

type Material struct {
	Id          string   `json:"id" bson:"id"`
	Name        string   `json:"name" bson:"name"`
	Supplier    Supplier `json:"supplier" bson:"supplier"`
	Origin      string   `json:"origin" bson:"origin"`
	Sustainable bool     `json:"sustainable" bson:"sustainable"`
	Details     string   `json:"details" bson:"details"`
	LastOrder   string   `json:"lastOrder" bson:"lastOrder"`
}

type MaterialModel struct {
	COLLECTION *mongo.Collection
}

// MaterialModel methods
func (m *MaterialModel) Add(material Material) error {
	company := os.Getenv("COMPANY")
	caser := cases.Title(language.English)
	companyFirstLMaiusc := caser.String(company)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var result bson.M
	err := m.COLLECTION.FindOne(ctx, bson.D{{"name", companyFirstLMaiusc}}).Decode(&result)

	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil
		}
		return err
	}

	filter := bson.D{{"name", companyFirstLMaiusc}}
	update := bson.D{{"$push", bson.D{{"materials", material}}}}

	_, err = m.COLLECTION.UpdateOne(ctx, filter, update)
	if err != nil {
		log.Println("Failed to insert material: ", err)
		return err
	}
	return nil
}

func (m *MaterialModel) Update(material Material) error {
	company := os.Getenv("COMPANY")
	caser := cases.Title(language.English)
	companyFirstLMaiusc := caser.String(company)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{"name": companyFirstLMaiusc, "materials.id": material.Id}
	update := bson.D{{"$set", bson.D{{"materials.$", material}}}}

	res, err := m.COLLECTION.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	if res.MatchedCount != 0 {
		log.Printf("matched and replaced material %v", material.Id)
		return nil
	}

	return errors.New("material not found")
}

func (m *MaterialModel) GetAll() ([]Material, error) {
	company := os.Getenv("COMPANY")
	caser := cases.Title(language.English)
	companyFirstLMaiusc := caser.String(company)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var result bson.M

	err := m.COLLECTION.FindOne(ctx, bson.D{{"name", companyFirstLMaiusc}}).Decode(&result)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}

	materialsRaw, ok := result["materials"]
	if !ok {
		return nil, errors.New("field 'materials' not found")
	}

	materialsJSON, err := json.Marshal(materialsRaw)
	if err != nil {
		return nil, err
	}

	var materials []Material
	err = json.Unmarshal(materialsJSON, &materials)
	if err != nil {
		return nil, err
	}

	return materials, nil
}

func (m *MaterialModel) GetOne(id string) (*Material, error) {
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
			{"$unwind", "$materials"},
		},
		bson.D{
			{"$match", bson.D{{"materials.id", id}}},
		},
	}

	cursor, err := m.COLLECTION.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}

	var results []bson.M
	if err := cursor.All(ctx, &results); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			log.Println("Object not found")
			return nil, nil
		} else {
			log.Fatal(err)
			return nil, err
		}
	}

	if len(results) > 0 {
		materialBytes, err := json.Marshal(results[0]["materials"])
		if err != nil {
			return nil, err
		}

		var myMaterial Material
		err = json.Unmarshal(materialBytes, &myMaterial)
		if err != nil {
			return nil, err
		}

		return &myMaterial, nil
	}

	return nil, fmt.Errorf("material with ID %v not found", id)
}

func (m *MaterialModel) GetBySupplier(supplierID string) (*[]Material, error) {
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
			{"$unwind", "$materials"},
		},
		bson.D{
			{"$match", bson.D{{"materials.supplier.id", supplierID}}},
		},
	}

	cursor, err := m.COLLECTION.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}

	var results []bson.M
	if err := cursor.All(ctx, &results); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			log.Println("Object not found")
			return nil, nil
		} else {
			log.Fatal(err)
			return nil, err
		}
	}

	if len(results) > 0 {
		var myMaterials []Material
		for _, material := range results {
			var mat Material
			materialBytes, err := json.Marshal(material["materials"])
			if err != nil {
				return nil, err
			}
			err = json.Unmarshal(materialBytes, &mat)
			if err != nil {
				return nil, err
			}
			myMaterials = append(myMaterials, mat)
		}

		return &myMaterials, nil
	}

	return nil, fmt.Errorf("materials for supplier with ID %v not found", supplierID)
}

func (m *MaterialModel) DeleteOne(id string) error {
	company := os.Getenv("COMPANY")
	caser := cases.Title(language.English)
	companyFirstLMaiusc := caser.String(company)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	material, err := m.GetOne(id)
	if err != nil {
		return err
	}

	if material != nil {
		filter := bson.D{{"name", companyFirstLMaiusc}}
		update := bson.D{
			{"$pull", bson.D{
				{"materials", bson.D{{"id", id}}},
			}},
		}

		res, err := m.COLLECTION.UpdateOne(ctx, filter, update)
		if err != nil {
			return err
		}
		if res.MatchedCount != 0 {
			log.Printf("matched and deleted material %v", id)
			return nil
		}
	}

	return errors.New("something went wrong")
}
