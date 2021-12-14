package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func setup() mongo.Database {
	clientOptions := options.Client().
		ApplyURI("mongodb+srv://Fiddler46:Fiddler46@cluster0.um5qb.mongodb.net/cities-nighthack?retryWrites=true&w=majority")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Fatal(err)
	}
	//defer client.Disconnect(ctx)
	DB := client.Database("cities-nighthack")
	return *DB
}

// model for user endpoint
type User struct {
	Users []User
	Email string `json:"email"`
}

// fake db to temp store users
var users []User

// checks if json is empty or not
func (u *User) IsEmpty() bool {
	return u.Email == ""
}

type server struct {
	router *mux.Router
	cities *mongo.Collection
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/user", createUser).Methods("POST")
	r.HandleFunc("/suggest", searchCity).Methods("GET")

	fmt.Println("Server running at port 8080")
	log.Fatal(http.ListenAndServe(":8080", r))

}

// func (s *server) handleIndex() http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		cities := s.cities.Find // something like that
// 		// the response
// 	}
// }

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

func searchCity(w http.ResponseWriter, r *http.Request) {
	DB := setup()
	values := r.URL.Query()
	city := values.Get("city_name")
	// params := mux.Vars(r)
	// city := params["city"]

	cityCollection := DB.Collection("city")

	cityList, err := cityCollection.Find(r.Context(), city, options.Find().SetProjection("all_names", "country_name"))
	if err != nil {
		panic(err)
	}
}

// `Find` is supposed to be structured like so
// `func (*mongo.Collection).Find(ctx context.Context, filter interface{}, opts ...*options.FindOptions) (*mongo.Cursor, error)`
// So I must be doing something wrong there cause my param isn't an interface, is it?
