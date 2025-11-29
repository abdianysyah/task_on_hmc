package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"github.com/go-chi/chi/v5"
	_ "github.com/go-sql-driver/mysql"
)

type Product struct {
	ID string `json:"id"`
	Name string `json:"name"`
	Price float64 `json:"price"`
}

var db *sql.DB

func ConnectMySQL() {
    var err error
    db, err = sql.Open("mysql", "root:Root12345@tcp(localhost:3306)/simple_api?parseTime=true")
    if err != nil {
        log.Fatalf("DB Open error: %v", err)
    }

    err = db.Ping()
    if err != nil {
        log.Fatalf("DB Ping error: %v", err)
    }

    log.Println("Database connected!")
}


func GetProducts(w http.ResponseWriter, r *http.Request)  {
	rows, err := db.Query("SELECT id, name, price FROM products")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var products []Product
	for rows.Next() {
		var p Product
		if err := rows.Scan(&p.ID, &p.Name, &p.Price); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		products = append(products, p)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(products)
}

func CreateProduct(w http.ResponseWriter, r *http.Request)  {
	var p Product

	err := json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, err = db.Exec("INSERT INTO products (id, name, price) VALUES (?, ?, ?),", p.ID , p.Name, p.Price)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func GetProductsByID(W http.ResponseWriter, r *http.Request)  {
	id := chi.URLParam(r, "id")

	var p Product
	err := db.QueryRow("SELECT id, name, price FROM products WHERE id = ?", id).Scan(&p.ID, &p.Name, &p.Price)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(W, "Product not found", http.StatusNotFound)
		} else {
			http.Error(W, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	W.Header().Set("Content-Type", "application/json")
	json.NewEncoder(W).Encode(p)
}

func UpdateProduct(w http.ResponseWriter, r *http.Request)  {
	id := chi.URLParam(r, "id")
	var p Product
	err := json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, err = db.Exec("UPDATE products SET name = ?, price = ? WHERE id = ?", p.Name, p.Price, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(p)
}

func DeleteProduct(w http.ResponseWriter, r *http.Request)  {
	id := chi.URLParam(r, "id")
	_, err := db.Exec("DELETE FROM products WHERE id = ?", id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Product deleted successfully"})
}

func main()  {
	ConnectMySQL()

	r := chi.NewRouter()
	r.Get("/products", GetProducts)
	r.Post("/products", CreateProduct)
	r.Get("/products/{id}", GetProductsByID)
	r.Put("/products/{id}", UpdateProduct)
	r.Delete("/products/{id}", DeleteProduct)

	log.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatal("Server failed to start:", err)
	}
}