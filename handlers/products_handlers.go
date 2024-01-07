package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"marvinhagler/helpers"
	"marvinhagler/models"
	"net/http"
)

type ProductsEnv struct {
	Products *models.ProductModel
}

func (env *ProductsEnv) AddProductHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		var productData models.Product

		err := json.NewDecoder(r.Body).Decode(&productData)
		if err != nil {
			http.Error(w, fmt.Sprintf("JSON Error: %v", err), http.StatusBadRequest)
			return
		}

		productData.Id = helpers.GenerateId("P-")

		err = env.Products.Add(productData)
		if err != nil {
			thisErr := fmt.Sprintf("%v", err)
			http.Error(w, thisErr, http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)

		err = json.NewEncoder(w).Encode(productData)
		if err != nil {
			log.Println("Failed to encode response:", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}

}

func (env *ProductsEnv) UpdateProductHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPut:
		var productData models.Product

		err := json.NewDecoder(r.Body).Decode(&productData)
		if err != nil {
			http.Error(w, fmt.Sprintf("JSON Error: %v", err), http.StatusBadRequest)
			return
		}

		err = env.Products.Update(productData)
		if err != nil {
			thisErr := fmt.Sprintf("%v", err)
			http.Error(w, thisErr, http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		err = json.NewEncoder(w).Encode(productData)
		if err != nil {
			log.Println("Failed to encode response:", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}

}

func (env *ProductsEnv) GetAllProductsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		products, err := env.Products.GetAll()
		if err != nil {
			thisErr := fmt.Sprintf("%v", err)
			http.Error(w, thisErr, http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		err = json.NewEncoder(w).Encode(products)
		if err != nil {
			log.Println("Failed to encode response:", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (env *ProductsEnv) GetOneProductHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		// /find-product?id=my_id
		id := r.URL.Query().Get("id")
		if len(id) < 20 || len(id) > 25 {
			http.Error(w, "Wrong ID format", http.StatusBadRequest)
			return
		}

		product, err := env.Products.GetOne(id)
		if err != nil {
			thisErr := fmt.Sprintf("%v", err)
			http.Error(w, thisErr, http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		err = json.NewEncoder(w).Encode(product)
		if err != nil {
			log.Println("Failed to encode response:", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (env *ProductsEnv) GetProductsByMaterialHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		// /find-by-material?material_id=id
		id := r.URL.Query().Get("material_id")
		if len(id) < 20 || len(id) > 25 {
			http.Error(w, "Wrong ID format", http.StatusBadRequest)
			return
		}

		products, err := env.Products.GetByMaterial(id)
		if err != nil {
			thisErr := fmt.Sprintf("%v", err)
			http.Error(w, thisErr, http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		err = json.NewEncoder(w).Encode(products)
		if err != nil {
			log.Println("Failed to encode response:", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (env *ProductsEnv) DeleteOneProductHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodDelete:
		// /delete-product?id=my_id
		id := r.URL.Query().Get("id")
		if len(id) < 20 || len(id) > 25 {
			http.Error(w, "Wrong ID format", http.StatusBadRequest)
			return
		}

		err := env.Products.DeleteOne(id)
		if err != nil {
			thisErr := fmt.Sprintf("%v", err)
			http.Error(w, thisErr, http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		err = json.NewEncoder(w).Encode("Product deleted")
		if err != nil {
			log.Println("Failed to encode response:", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
