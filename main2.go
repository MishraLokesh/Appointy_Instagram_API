package main

import (
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)



// Main function
func main() {
	// Init router
	r := mux.NewRouter()

	// Hardcoded data - @todo: add database
	users = append(users, User{ID: "1", Name: "Lokesh", Email: "user_One@gmail.com", Password: "yoyo"})
	users = append(users, User{ID: "2", Name: "Aman", Email: "user_Two@gmail.com", Password: "oyoy"})
	users = append(users, User{ID: "3", Name: "Samar", Email: "user_Three@gmail.com", Password: "oyyo"})

	// Route handles & endpoints
	r.HandleFunc("/users", getUsers).Methods("GET")
	r.HandleFunc("/users/{id}", getUser).Methods("GET")
	r.HandleFunc("/users", createUser).Methods("POST")
	r.HandleFunc("/users/{id}", updateUser).Methods("PUT")
	r.HandleFunc("/users/{id}", deleteUser).Methods("DELETE")

	// Start server
	log.Fatal(http.ListenAndServe(":8080", r))
}
