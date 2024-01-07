package routes

import (
	"marvinhagler/handlers"
	"net/http"
)

func ProductsRouter(router *http.ServeMux, env *handlers.ProductsEnv) {
	router.HandleFunc("/products/add", env.AddProductHandler)
	router.HandleFunc("/products/update", env.UpdateProductHandler)
	router.HandleFunc("/products/all", env.GetAllProductsHandler)
	router.HandleFunc("/products/find-product", env.GetOneProductHandler)
	router.HandleFunc("/products/find-by-material", env.GetProductsByMaterialHandler)
	router.HandleFunc("/products/delete-product", env.DeleteOneProductHandler)
}

func MaterialsRouter(router *http.ServeMux, env *handlers.MaterialsEnv) {
	router.HandleFunc("/materials/add", env.AddMaterialHandler)
	router.HandleFunc("/materials/update", env.UpdateMaterialHandler)
	router.HandleFunc("/materials/all", env.GetAllMaterialsHandler)
	router.HandleFunc("/materials/find-material", env.GetOneMaterialHandler)
	router.HandleFunc("/materials/find-by-supplier", env.GetMaterialsBySupplierHandler)
	router.HandleFunc("/materials/delete-material", env.DeleteOneMaterialHandler)
}

func SuppliersRouter(router *http.ServeMux, env *handlers.SuppliersEnv) {
	router.HandleFunc("/suppliers/add", env.AddSupplierHandler)
	router.HandleFunc("/suppliers/update", env.UpdateSupplierHandler)
	router.HandleFunc("/suppliers/all", env.GetAllSuppliersHandler)
	router.HandleFunc("/suppliers/find-supplier", env.GetOneSupplierHandler)
	router.HandleFunc("/suppliers/delete-supplier", env.DeleteOneSupplierHandler)
}

func CertsRouter(router *http.ServeMux, env *handlers.CertsEnv) {
	router.HandleFunc("/certs/add", env.AddCertHandler)
	router.HandleFunc("certs/all", env.GetAllCertsHandler)
}

func CompanyRouter(router *http.ServeMux, env *handlers.CompanyEnv) {
	router.HandleFunc("/company/init", env.InitializeCompanyHandler)
}
