package main

import (
	"encoding/json"
	"net/http"
	"strconv"
)

type api struct {
	store *ProductStore
}

func writeJSON(w http.ResponseWriter, status int, body any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(body)
}

func errorJSON(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}

func (a *api) getProduct(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		errorJSON(w, http.StatusBadRequest, "Invalid ID")
		return
	}
	p, ok := a.store.Get(int64(id))
	if !ok {
		errorJSON(w, http.StatusNotFound, "Product not found")
		return
	}
	writeJSON(w, http.StatusOK, p)
}

func (a *api) createProduct(w http.ResponseWriter, r *http.Request) {
	var input Product
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		errorJSON(w, http.StatusBadRequest, "Invalid JSON: "+err.Error())
		return
	}
	if input.Name == "" || input.Price <= 0 {
		errorJSON(w, http.StatusUnprocessableEntity, "Name field is required and Price field must be a positive number")
		return
	}
	created, err := a.store.Create(input)
	if err != nil {
		errorJSON(w, http.StatusInternalServerError, "Internal error")
		return
	}
	writeJSON(w, http.StatusCreated, created)
}

func (a *api) updateProduct(w http.ResponseWriter, r *http.Request) {
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
	updated, ok, err := a.store.Update(int64(id), input)
	if err != nil {
		errorJSON(w, http.StatusInternalServerError, "Internal error")
		return
	}
	if !ok {
		errorJSON(w, http.StatusNotFound, "Product not found")
		return
	}
	writeJSON(w, http.StatusOK, updated)
}

func (a *api) deleteProduct(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		errorJSON(w, http.StatusBadRequest, "Invalid ID")
		return
	}
	ok, err := a.store.Delete(int64(id))
	if err != nil {
		errorJSON(w, http.StatusInternalServerError, "Internal error")
		return
	}
	if !ok {
		errorJSON(w, http.StatusNotFound, "Product not found")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (a *api) listProducts(w http.ResponseWriter, r *http.Request) {
	products, err := a.store.List()
	if err != nil {
		errorJSON(w, http.StatusInternalServerError, "Internal error")
		return
	}
	writeJSON(w, http.StatusOK, products)
}
