package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
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
	// Users []User
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
	r.HandleFunc("/suggest", searchCity).Methods("GET")

	fmt.Println("Server running at port 8080")
	log.Fatal(http.ListenAndServe(":8080", r))

}

func createUser(w http.ResponseWriter, r *http.Request) {
	DB := setup()
	userCollection := DB.Collection("user")
	// w.Header().Set("Content-Type", "application/json")
	// if r.Body == nil {
	// 	json.NewEncoder(w).Encode("Must send data")
	// }

	var user User
	_ = json.NewDecoder(r.Body).Decode(&user)
	inserted, err := userCollection.InsertOne(context.Background(), user)
	if err != nil {
		log.Fatal(err)
	}
	json.NewEncoder(w).Encode(user)
	fmt.Println("Inserted user into db: ", inserted.InsertedID)
	// err := json.NewDecoder(r.Body).Decode(&user)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// if user.IsEmpty() {
	// 	json.NewEncoder(w).Encode("Invalid! Enter user email.")
	// 	return
	// }
	// users = append(users, user)
	// json.NewEncoder(w).Encode(user)

}

func searchCity(w http.ResponseWriter, r *http.Request) {
	//ctx := context.Background()
	DB := setup()
	values := r.URL.Query()
	city := values.Get("city_name")

	// projection := bson.M{
	// 	"all_names":    *city,
	// 	"country_name": *city,
	// }

	// params := mux.Vars(r)
	// city := params["city"]

	cityCollection := DB.Collection("city")

	cursor, err := cityCollection.Find(r.Context(), bson.E{"$city", city}) // options.Find().SetProjection(projection))
	if err != nil {
		log.Fatal(err)
	}

	var cityList []bson.M
	if err = cursor.All(r.Context(), &cityList); err != nil {
		log.Fatal(err)
	}
	for _, cityList := range cityList {
		fmt.Println(cityList["all_names"])
		fmt.Println(cityList["country_name"])
	}

	// defer cursor.Close(ctx)
	// for cursor.Next(ctx) {
	// 	var cityList bson.M
	// 	if err = cursor.Decode(&cityList); err != nil {
	// 		log.Fatal(err)
	// 	}
	// 	fmt.Println(cityList["all_names"])
	// 	fmt.Println(cityList["country_name"])
	// }
}

func userAuth(w http.ResponseWriter r *http.Request) {
	
}