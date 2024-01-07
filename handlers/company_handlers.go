package handlers

import (
	"encoding/json"
	"fmt"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"log"
	"marvinhagler/models"
	"net/http"
	"os"
)

type CompanyEnv struct {
	Company *models.CompanyModel
}

func (env *CompanyEnv) InitializeCompanyHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:

		companyFromEnv := os.Getenv("COMPANY")
		caser := cases.Title(language.English)
		companyNameCaser := caser.String(companyFromEnv)

		newCompany := models.Company{
			Name:      companyNameCaser,
			Products:  []models.Product{},
			Materials: []models.Material{},
			Suppliers: []models.Supplier{},
			Certs:     []models.Cert{},
		}

		err := env.Company.Initialize(newCompany)
		if err != nil {
			thisErr := fmt.Sprintf("%v", err)
			http.Error(w, thisErr, http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)

		err = json.NewEncoder(w).Encode(newCompany.Name)
		if err != nil {
			log.Println("Failed to encode response:", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}

}
