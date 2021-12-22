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
	"go.mongodb.org/mongo-driver/bson/primitive"
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

type City struct {
	Name string `bson:"name" json:"name"`
	// CountryName string `bson:"country_name" json:"country"`
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
	// ctx := context.Background()
	w.Header().Set("Content-Type", "application/json")
	DB := setup()
	values := r.URL.Query()
	city := values["city_name"]
	cityCollection := DB.Collection("city")
	fmt.Println(city)

	if len(city) == 0 {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"msg": "Search key required"}`))
		return
	}
	// projection := bson.M{
	// 	"all_names":    *city,
	// 	"country_name": *city,
	// }

	// params := mux.Vars(r)
	// city := params["city"]

	filter := bson.D{
		primitive.E{
			Key: "all_names", Value: primitive.Regex{
				Pattern: city[0], Options: "i",
			},
		},
	}

	cursor, err := cityCollection.Find(r.Context(), filter) // options.Find().SetProjection(projection))
	if err != nil {
		log.Fatal(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"msg": "There was an error. Try again later"}`))
		return
	}

	var cityList []City
	for cursor.Next(context.TODO()) {
		var cities City
		err := cursor.Decode(&cities)
		if err != nil {
			log.Fatal(err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"msg": "There was an error. Try again later"}`))
			return
		}
		cityList = append(cityList, cities)
	}
	w.WriteHeader(http.StatusOK)
	if len(cityList) == 0 {
		w.Write([]byte(`{"cities": []}`))
		return
	}
	json.NewEncoder(w).Encode(cityList)
	// if err = cursor.All(r.Context(), &cityList); err != nil {
	// 	http.Error(w, err.Error(), http.StatusInternalServerError)
	// }

	// for _, cityList := range cityList {
	// 	fmt.Println(cityList["all_names"])
	// 	fmt.Println(cityList["country_name"])
	// }

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
