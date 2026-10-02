package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

type Product struct {
	ID    int64   `json:"id"`
	Name  string  `json:"name"`
	Price float64 `json:"price"`
}

type ProductStore struct {
	db *sql.DB
}

func NewProductStore(db *sql.DB) *ProductStore {
	return &ProductStore{db: db}
}

func (s *ProductStore) List() ([]Product, error) {
	ctx := context.Background()
	rows, err := s.db.QueryContext(ctx, "SELECT id, name, price FROM products ORDER BY id")
	if err != nil {
		return nil, fmt.Errorf("Error fetching products: %w", err)
	}
	defer rows.Close()
	out := []Product{}
	for rows.Next() {
		var p Product
		if err := rows.Scan(&p.ID, &p.Name, &p.Price); err != nil {
			return nil, fmt.Errorf("Error fetching products: %w", err)
		}
		out = append(out, p)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("Error fetching products: %w", err)
	}
	return out, nil
}

func (s *ProductStore) Create(p Product) (Product, error) {
	ctx := context.Background()
	err := s.db.QueryRowContext(ctx, "INSERT INTO products (name, price) VALUES ($1, $2) RETURNING id", p.Name, p.Price).
		Scan(&p.ID)
	if err != nil {
		return Product{}, fmt.Errorf("Error inserting Product into database: %w", err)
	}
	return p, nil
}

func (s *ProductStore) Get(id int64) (Product, bool) {
	ctx := context.Background()
	var p Product
	err := s.db.QueryRowContext(ctx, "SELECT id, name, price FROM products WHERE id = $1", id).
		Scan(&p.ID, &p.Name, &p.Price)
	if errors.Is(err, sql.ErrNoRows) {
		return Product{}, false
	}
	if err != nil {
		return Product{}, false
	}
	return p, true
}

func (s *ProductStore) Update(id int64, p Product) (Product, bool, error) {
	ctx := context.Background()
	res, err := s.db.ExecContext(ctx,
		"UPDATE products SET name = $1, price = $2 WHERE id = $3",
		p.Name, p.Price, id)
	if err != nil {
		return Product{}, false, fmt.Errorf("atualizando: %w", err)
	}
	n, err := res.RowsAffected()
	if err != nil {
		return Product{}, false, fmt.Errorf("rowsAffected: %w", err)
	}
	if n == 0 {
		return Product{}, false, nil
	}
	novo, ok := s.Get(id)
	return novo, ok, nil
}

func (s *ProductStore) Delete(id int64) (bool, error) {
	ctx := context.Background()
	res, err := s.db.ExecContext(ctx, "DELETE FROM products WHERE id = $1;", id)
	if err != nil {
		return false, fmt.Errorf("Error deleting Product with ID %d: %w", id, err)
	}
	n, err := res.RowsAffected()
	if err != nil {
		return false, fmt.Errorf("Error deleting Product with ID %d: %w", id, err)
	}
	return n > 0, nil
}
