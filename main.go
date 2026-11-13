package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"product-management-api/handlers"
	"product-management-api/repositories"

	"github.com/joho/godotenv"

	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("env ulanishda xatolik ")
	}
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgres://postgres:postgres@localhost:5432/products_db?sslmode=disable"
	}
	ctx := context.Background()
	pool, err := pgxpool.New(ctx, dbURL)
	if err != nil {
		log.Fatalf("Bazaga ulanib bo'lmadi: %v", err)
	}
	defer pool.Close()
	if err := pool.Ping(ctx); err != nil {
		log.Fatalf("Bazaga ping o'tmadi: %v", err)
	}
	productRepo := &repositories.ProductRepository{DB: pool}
	productHandler := &handlers.ProductHandler{Repo: productRepo}
	mux := http.NewServeMux()

	mux.HandleFunc("POST /products", productHandler.Create)
	mux.HandleFunc("GET /products", productHandler.GetAll)
	mux.HandleFunc("GET /products/{id}", productHandler.GetByID)
	mux.HandleFunc("PUT /products/{id}", productHandler.Update)
	mux.HandleFunc("DELETE /products/{id}", productHandler.Delete)
	mux.HandleFunc("GET /products/search", productHandler.Search)
	mux.HandleFunc("GET /products/sort", productHandler.Sort)
	mux.HandleFunc("GET /products/pagination", productHandler.Pagination)
	log.Println("Server :8080 portida ishga tushdi...")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatalf("Server xatoligi: %v", err)
	}
}
