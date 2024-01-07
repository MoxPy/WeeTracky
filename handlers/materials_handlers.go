package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"marvinhagler/helpers"
	"marvinhagler/models"
	"net/http"
)

type MaterialsEnv struct {
	Materials *models.MaterialModel
}

func (env *MaterialsEnv) AddMaterialHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		var materialData models.Material

		err := json.NewDecoder(r.Body).Decode(&materialData)
		if err != nil {
			http.Error(w, fmt.Sprintf("JSON Error: %v", err), http.StatusBadRequest)
			return
		}

		materialData.Id = helpers.GenerateId("M-")

		err = env.Materials.Add(materialData)
		if err != nil {
			thisErr := fmt.Sprintf("%v", err)
			http.Error(w, thisErr, http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)

		err = json.NewEncoder(w).Encode(materialData)
		if err != nil {
			log.Println("Failed to encode response:", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}

}

func (env *MaterialsEnv) UpdateMaterialHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPut:
		var materialData models.Material

		err := json.NewDecoder(r.Body).Decode(&materialData)
		if err != nil {
			http.Error(w, fmt.Sprintf("JSON Error: %v", err), http.StatusBadRequest)
			return
		}

		err = env.Materials.Update(materialData)
		if err != nil {
			thisErr := fmt.Sprintf("%v", err)
			http.Error(w, thisErr, http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		err = json.NewEncoder(w).Encode(materialData)
		if err != nil {
			log.Println("Failed to encode response:", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}

}

func (env *MaterialsEnv) GetAllMaterialsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		materials, err := env.Materials.GetAll()
		if err != nil {
			thisErr := fmt.Sprintf("%v", err)
			http.Error(w, thisErr, http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		err = json.NewEncoder(w).Encode(materials)
		if err != nil {
			log.Println("Failed to encode response:", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (env *MaterialsEnv) GetMaterialsBySupplierHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		// /find-by-supplier?supplier_id=id
		id := r.URL.Query().Get("supplier_id")
		if len(id) < 20 || len(id) > 25 {
			http.Error(w, "Wrong ID format", http.StatusBadRequest)
			return
		}

		materials, err := env.Materials.GetBySupplier(id)
		if err != nil {
			thisErr := fmt.Sprintf("%v", err)
			http.Error(w, thisErr, http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		err = json.NewEncoder(w).Encode(materials)
		if err != nil {
			log.Println("Failed to encode response:", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (env *MaterialsEnv) GetOneMaterialHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		// /find-material?id=my_id
		id := r.URL.Query().Get("id")
		if len(id) < 20 || len(id) > 25 {
			http.Error(w, "Wrong ID format", http.StatusBadRequest)
			return
		}

		material, err := env.Materials.GetOne(id)
		if err != nil {
			thisErr := fmt.Sprintf("%v", err)
			http.Error(w, thisErr, http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		err = json.NewEncoder(w).Encode(material)
		if err != nil {
			log.Println("Failed to encode response:", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (env *MaterialsEnv) DeleteOneMaterialHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodDelete:
		// /delete-material?id=my_id
		id := r.URL.Query().Get("id")
		if len(id) < 20 || len(id) > 25 {
			http.Error(w, "Wrong ID format", http.StatusBadRequest)
			return
		}

		err := env.Materials.DeleteOne(id)
		if err != nil {
			thisErr := fmt.Sprintf("%v", err)
			http.Error(w, thisErr, http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		err = json.NewEncoder(w).Encode("Material deleted")
		if err != nil {
			log.Println("Failed to encode response:", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
