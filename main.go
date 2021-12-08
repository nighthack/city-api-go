package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

// model for user endpoint
type User struct {
	Email string `json:"email"`
}

// fake db to temp store users
var users []User

// checks if json is empty or not
func (u *User) IsEmpty() bool {
	return u.Email == ""
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/user", createUser).Methods("POST")
	fmt.Println("Server running at port 8080")
	log.Fatal(http.ListenAndServe(":8080", r))

}

func createUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Body == nil {
		json.NewEncoder(w).Encode("Must send data")
	}

	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		log.Fatal(err)
	}
	if user.IsEmpty() {
		json.NewEncoder(w).Encode("Invalid! Enter user email.")
		return
	}
	users = append(users, user)
	json.NewEncoder(w).Encode(user)

}
