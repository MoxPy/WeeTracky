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

type Product struct {
	Id                 string     `json:"id" bson:"id"`
	Name               string     `json:"name" bson:"name"`
	MadeIn             string     `json:"made_in" bson:"made_in"`
	Materials          []Material `json:"materials" bson:"materials"`
	Price              float64    `json:"price" bson:"price"`
	Description        string     `json:"description" bson:"description"`
	SustainablePackage bool       `json:"sustainablePackage" bson:"sustainablePackage"`
}

type ProductModel struct {
	DB         *mongo.Client
	COLLECTION *mongo.Collection
}

// ProductModel methods
func (p *ProductModel) Add(product Product) error {
	company := os.Getenv("COMPANY")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	caser := cases.Title(language.English)
	companyFirstLMaiusc := caser.String(company)
	var result bson.M
	err := p.COLLECTION.FindOne(ctx, bson.D{{"name", companyFirstLMaiusc}}).Decode(&result)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil
		}
		return err
	}

	filter := bson.D{{"name", companyFirstLMaiusc}}
	update := bson.D{{"$push", bson.D{{"products", product}}}}

	_, err = p.COLLECTION.UpdateOne(ctx, filter, update)
	if err != nil {
		log.Println("Failed to insert product: ", err)
		return err
	}
	return nil
}

func (p *ProductModel) Update(product Product) error {
	company := os.Getenv("COMPANY")
	caser := cases.Title(language.English)
	companyFirstLMaiusc := caser.String(company)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{"name": companyFirstLMaiusc, "products.id": product.Id}
	update := bson.D{{"$set", bson.D{{"products.$", product}}}}

	res, err := p.COLLECTION.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	if res.MatchedCount != 0 {
		fmt.Printf("matched and replaced product %v", product.Id)
		return nil
	}

	return errors.New("product not found")
}

func (p *ProductModel) GetAll() ([]Product, error) {
	company := os.Getenv("COMPANY")
	caser := cases.Title(language.English)
	companyFirstLMaiusc := caser.String(company)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var result bson.M

	err := p.COLLECTION.FindOne(ctx, bson.D{{"name", companyFirstLMaiusc}}).Decode(&result)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}

	productsRaw, ok := result["products"]
	if !ok {
		return nil, errors.New("field 'products' not found")
	}

	productsJSON, err := json.Marshal(productsRaw)
	if err != nil {
		return nil, err
	}

	var products []Product
	err = json.Unmarshal(productsJSON, &products)
	if err != nil {
		return nil, err
	}

	return products, nil
}

func (p *ProductModel) GetOne(id string) (*Product, error) {
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
			{"$unwind", "$products"},
		},
		bson.D{
			{"$match", bson.D{{"products.id", id}}},
		},
	}

	cursor, err := p.COLLECTION.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}

	/* bson.M regular map[string]interface{} when encoding and decoding,
	una mappa in cui le chiavi sono stringhe e i valori sono di
	tipo interface{}, il che significa che possono essere di qualsiasi tipo */
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
		productBytes, err := json.Marshal(results[0]["products"])
		if err != nil {
			return nil, err
		}

		var myProduct Product
		err = json.Unmarshal(productBytes, &myProduct)
		if err != nil {
			return nil, err
		}

		return &myProduct, nil
	}

	return nil, fmt.Errorf("product with ID %v not found", id)
}

func (p *ProductModel) GetByMaterial(materialID string) (*[]Product, error) {
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
			{"$unwind", "$products"},
		},
		bson.D{
			{"$match", bson.D{{"products.materials.id", materialID}}},
		},
	}

	cursor, err := p.COLLECTION.Aggregate(ctx, pipeline)
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
		var myProducts []Product
		for _, product := range results {
			var prod Product
			productBytes, err := json.Marshal(product["products"])
			if err != nil {
				return nil, err
			}
			err = json.Unmarshal(productBytes, &prod)
			if err != nil {
				return nil, err
			}
			myProducts = append(myProducts, prod)
		}

		return &myProducts, nil
	}

	return nil, fmt.Errorf("products with material with ID %v not found", materialID)
}

func (p *ProductModel) DeleteOne(id string) error {
	company := os.Getenv("COMPANY")
	caser := cases.Title(language.English)
	companyFirstLMaiusc := caser.String(company)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	product, err := p.GetOne(id)
	if err != nil {
		return err
	}

	if product != nil {
		filter := bson.D{{"name", companyFirstLMaiusc}}
		update := bson.D{
			{"$pull", bson.D{
				{"products", bson.D{{"id", id}}},
			}},
		}

		res, err := p.COLLECTION.UpdateOne(ctx, filter, update)
		if err != nil {
			return err
		}
		if res.MatchedCount != 0 {
			fmt.Printf("matched and deleted product %v", id)
			return nil
		}
	}

	return errors.New("something went wrong")
}
