package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"marvinhagler/helpers"
	"marvinhagler/models"
	"net/http"
)

type SuppliersEnv struct {
	Suppliers *models.SupplierModel
}

func (env *SuppliersEnv) AddSupplierHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		var supplierData models.Supplier

		err := json.NewDecoder(r.Body).Decode(&supplierData)
		if err != nil {
			http.Error(w, fmt.Sprintf("JSON Error: %v", err), http.StatusBadRequest)
			return
		}

		supplierData.Id = helpers.GenerateId("S-")

		err = env.Suppliers.Add(supplierData)
		if err != nil {
			thisErr := fmt.Sprintf("%v", err)
			http.Error(w, thisErr, http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)

		err = json.NewEncoder(w).Encode(supplierData)
		if err != nil {
			log.Println("Failed to encode response:", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}

}

func (env *SuppliersEnv) UpdateSupplierHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPut:
		var supplierData models.Supplier

		err := json.NewDecoder(r.Body).Decode(&supplierData)
		if err != nil {
			http.Error(w, fmt.Sprintf("JSON Error: %v", err), http.StatusBadRequest)
			return
		}

		err = env.Suppliers.Update(supplierData)
		if err != nil {
			thisErr := fmt.Sprintf("%v", err)
			http.Error(w, thisErr, http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		err = json.NewEncoder(w).Encode(supplierData)
		if err != nil {
			log.Println("Failed to encode response:", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}

}

func (env *SuppliersEnv) GetAllSuppliersHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		suppliers, err := env.Suppliers.GetAll()
		if err != nil {
			thisErr := fmt.Sprintf("%v", err)
			http.Error(w, thisErr, http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		err = json.NewEncoder(w).Encode(suppliers)
		if err != nil {
			log.Println("Failed to encode response:", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (env *SuppliersEnv) GetOneSupplierHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		// /find-supplier?id=my_id
		id := r.URL.Query().Get("id")
		if len(id) < 20 || len(id) > 25 {
			http.Error(w, "Wrong ID format", http.StatusBadRequest)
			return
		}

		supplier, err := env.Suppliers.GetOne(id)
		if err != nil {
			thisErr := fmt.Sprintf("%v", err)
			http.Error(w, thisErr, http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		err = json.NewEncoder(w).Encode(supplier)
		if err != nil {
			log.Println("Failed to encode response:", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (env *SuppliersEnv) DeleteOneSupplierHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodDelete:
		// /delete-supplier?id=my_id
		id := r.URL.Query().Get("id")
		if len(id) < 20 || len(id) > 25 {
			http.Error(w, "Wrong ID format", http.StatusBadRequest)
			return
		}

		err := env.Suppliers.DeleteOne(id)
		if err != nil {
			thisErr := fmt.Sprintf("%v", err)
			http.Error(w, thisErr, http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		err = json.NewEncoder(w).Encode("Supplier deleted")
		if err != nil {
			log.Println("Failed to encode response:", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
