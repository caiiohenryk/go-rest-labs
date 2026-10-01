package main

import "sync"

type Product struct {
	ID    int     `json:"id"`
	Name  string  `json:"name"`
	Price float64 `json:"price"`
}

type ProductStore struct {
	mu       sync.Mutex
	nextID   int
	products map[int]Product
}

func NewProductStore() *ProductStore {
	return &ProductStore{nextID: 1, products: make(map[int]Product)}
}

func (s *ProductStore) List() []Product {
	s.mu.Lock()
	defer s.mu.Unlock()
	out := make([]Product, 0, len(s.products))
	for _, p := range s.products {
		out = append(out, p)
	}
	return out
}

func (s *ProductStore) Create(p Product) Product {
	s.mu.Lock()
	defer s.mu.Unlock()
	p.ID = s.nextID
	s.nextID++
	s.products[p.ID] = p
	return p
}

func (s *ProductStore) Get(id int) (Product, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	p, ok := s.products[id]
	return p, ok
}

func (s *ProductStore) Update(id int, p Product) (Product, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.products[id]; ok {
		return Product{}, false
	}
	p.ID = id
	s.products[id] = p
	return p, true
}

func (s *ProductStore) Delete(id int) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.products[id]; !ok {
		return false
	}
	delete(s.products, id)
	return true
}
