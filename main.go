package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type User struct {
	Email       string `json:"email"`
	AccessToken string `json:"token"`
}

type Users struct {
	Users []User
}

func new() *Users {
	return &Users{}
}

func main() {
	http.HandleFunc("/user", createUser).Methods(http.MethodPost)
	log.Fatal(http.ListenAndServe(":8080", nil)

	setup()
}



func createUser(r *Users) (user *User) {
	
	r.Users = append(r.Users, user)
	

}