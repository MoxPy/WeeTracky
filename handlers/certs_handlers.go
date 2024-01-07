package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"marvinhagler/helpers"
	"marvinhagler/models"
	"net/http"
)

type CertsEnv struct {
	Certs *models.CertModel
}

func (env *CertsEnv) AddCertHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		var certData models.Cert

		err := json.NewDecoder(r.Body).Decode(&certData)
		if err != nil {
			http.Error(w, fmt.Sprintf("JSON Error: %v", err), http.StatusBadRequest)
			return
		}

		certData.Id = helpers.GenerateId("CERT-")

		err = env.Certs.Add(certData)
		if err != nil {
			thisErr := fmt.Sprintf("%v", err)
			http.Error(w, thisErr, http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)

		err = json.NewEncoder(w).Encode(certData)
		if err != nil {
			log.Println("Failed to encode response:", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}

}

func (env *CertsEnv) GetAllCertsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		certs, err := env.Certs.GetAll()
		if err != nil {
			thisErr := fmt.Sprintf("%v", err)
			http.Error(w, thisErr, http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		err = json.NewEncoder(w).Encode(certs)
		if err != nil {
			log.Println("Failed to encode response:", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
