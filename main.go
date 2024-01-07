package main

import (
	"context"
	"errors"
	"github.com/joho/godotenv"
	"log"
	"marvinhagler/db"
	"marvinhagler/handlers"
	"marvinhagler/models"
	"marvinhagler/routes"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	client, collection, err := db.ConnectMongoDB()
	if err != nil {
		log.Fatal(err)
		return
	}
	defer db.DisconnectMongoDB(client)

	productsEnv := &handlers.ProductsEnv{Products: &models.ProductModel{DB: client, COLLECTION: collection}}
	materialsEnv := &handlers.MaterialsEnv{Materials: &models.MaterialModel{DB: client, COLLECTION: collection}}
	suppliersEnv := &handlers.SuppliersEnv{Suppliers: &models.SupplierModel{DB: client, COLLECTION: collection}}
	certsEnv := &handlers.CertsEnv{Certs: &models.CertModel{DB: client, COLLECTION: collection}}
	companyEnv := &handlers.CompanyEnv{Company: &models.CompanyModel{DB: client, COLLECTION: collection}}

	mux := http.NewServeMux()
	routes.ProductsRouter(mux, productsEnv)
	routes.MaterialsRouter(mux, materialsEnv)
	routes.SuppliersRouter(mux, suppliersEnv)
	routes.CertsRouter(mux, certsEnv)
	routes.CompanyRouter(mux, companyEnv)

	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	// Start server
	go func() {
		log.Printf("Server listening on %v\n", server.Addr)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatal(err)
		}
	}()

	// Wait signal for closing db connection
	interruptChan := make(chan os.Signal, 1)
	signal.Notify(interruptChan, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	<-interruptChan

	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Close server and handling err
	if err := server.Shutdown(ctx); err != nil {
		log.Fatal(err)
	}

	log.Println("Server stopped")
}
