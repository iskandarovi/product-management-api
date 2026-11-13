package handlers

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"product-management-api/models"
	"product-management-api/repositories"
	"strconv"
	"strings"

	"github.com/jackc/pgx/v5"
)

type ProductHandler struct {
	Repo *repositories.ProductRepository
}

func (h *ProductHandler) Create(w http.ResponseWriter, r *http.Request) {
	var input productInput

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "notogri json malumotlar yuborildi", http.StatusBadRequest)
		return
	}
	if err := input.validate(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if input.Name == "" || input.Price <= 0 {
		http.Error(w, "mahsulot nomi va narxi togri kiritilishi kerak", http.StatusBadRequest)
		return
	}

	newProduct := models.Product{
		Name:     input.Name,
		Price:    input.Price,
		Quantity: input.Quantity,
	}

	createdProduct, err := h.Repo.Create(r.Context(), newProduct)
	if err != nil {
		log.Printf("POST /products da xatolik: %v", err)
		http.Error(w, "Mahsulotni saqlashda server xatosi", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(createdProduct)
}

func (h *ProductHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	produc, err := h.Repo.GetAll(r.Context())
	if err != nil {
		log.Printf("Get /products: %v", err)
		http.Error(w, "Mahsulotlarni olishda xatolik", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(produc)
}

func (h *ProductHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil || id <= 0 {
		http.Error(w, "ID noto'g'ri", http.StatusBadRequest)
		return
	}

	product, err := h.Repo.GetByID(r.Context(), id)
	if err != nil {
		if err == pgx.ErrNoRows {
			http.Error(w, "produc topilmadi", http.StatusNotFound)
			return
		}

		log.Printf("GET /products/%d: %v", id, err)
		http.Error(w, "Server xatosi", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(product)
}

func (h *ProductHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil || id <= 0 {
		http.Error(w, "id xato!", http.StatusBadRequest)
		return
	}
	var input productInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Json xato kiritilgan", http.StatusBadRequest)
		return
	}
	if err := input.validate(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	updatedProduct, err := h.Repo.Update(
		r.Context(),
		id,
		input.Name,
		input.Price,
		input.Quantity,
	)
	if err != nil {

		if err == pgx.ErrNoRows {
			http.Error(w, "Mahsulot topilmadi", 404)
			return
		}
		log.Printf("PUT /products/%d xatolik: %v", id, err)
		http.Error(w, "seerver xatoligi", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(updatedProduct); err != nil {
		http.Error(w, "jsonda xatolik", http.StatusInternalServerError)
		return
	}
}

func (h *ProductHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil || id <= 0 {
		http.Error(w, "ID xato", http.StatusBadRequest)
		return
	}

	err = h.Repo.Delete(r.Context(), id)
	if err != nil {
		if err == pgx.ErrNoRows {
			http.Error(w, "Mahsulot topilmadi", 404)
			return
		}

		log.Printf("DELETE /products/%d xatolik: %v", id, err)
		http.Error(w, "Server xatosi", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

type productInput struct {
	Name     string  `json:"name"`
	Price    float64 `json:"price"`
	Quantity int     `json:"quantity"`
}

func (p productInput) validate() error {
	if strings.TrimSpace(p.Name) == "" {
		return errors.New("nomni bosh kiritmang")
	}
	if p.Price <= 0 {
		return errors.New("narx musbat bolsin")
	}
	if p.Quantity < 0 {
		return errors.New("miqdor musbat bolsin")
	}
	return nil
}
