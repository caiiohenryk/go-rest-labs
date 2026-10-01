package main

import (
	"encoding/json"
	"net/http"
	"strconv"
)

func writeJSON(w http.ResponseWriter, status int, body any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(body)
}

func errorJSON(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}

var store = NewProductStore()

func getProduct(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		errorJSON(w, http.StatusBadRequest, "Invalid ID")
		return
	}
	p, ok := store.Get(id)
	if !ok {
		errorJSON(w, http.StatusNotFound, "Product not found")
		return
	}
	writeJSON(w, http.StatusOK, p)
}

func createProduct(w http.ResponseWriter, r *http.Request) {
	var input Product
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		errorJSON(w, http.StatusBadRequest, "Invalid JSON: "+err.Error())
		return
	}
	if input.Name == "" || input.Price <= 0 {
		errorJSON(w, http.StatusUnprocessableEntity, "Name field is required and Price field must be a positive number")
		return
	}
	writeJSON(w, http.StatusCreated, store.Create(input))
}

func updateProduct(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		errorJSON(w, http.StatusBadRequest, "Invalid ID: ")
		return
	}
	var input Product
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		errorJSON(w, http.StatusBadRequest, "Invalid JSON: "+err.Error())
		return
	}
	if input.Name == "" || input.Price <= 0 {
		errorJSON(w, http.StatusUnprocessableEntity, "Name field is required and Price field must be a positive number")
		return
	}
	updated, ok := store.Update(id, input)
	if !ok {
		errorJSON(w, http.StatusNotFound, "Product not found")
		return
	}
	writeJSON(w, http.StatusOK, updated)
}

func deleteProduct(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		errorJSON(w, http.StatusBadRequest, "Invalid ID")
		return
	}
	if !store.Delete(id) {
		errorJSON(w, http.StatusNotFound, "Product not found")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func listProducts(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, store.List())
}
