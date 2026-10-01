package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "ok")
	})
	mux.HandleFunc("POST /products", createProduct)
	mux.HandleFunc("GET /products", listProducts)
	mux.HandleFunc("GET /products/{id}", getProduct)
	mux.HandleFunc("DELETE /products/{id}", deleteProduct)
	mux.HandleFunc("PATCH /products/{id}", updateProduct)
	srv := &http.Server{Addr: ":3000", Handler: withLogging(mux)}
	log.Println("Listening on :3000")
	srv.ListenAndServe()
}

func withLogging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	})
}
