package main

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var client mongo.Client
var characterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

func setup() {
	clientOptions := options.Client().
		ApplyURI("mongodb+srv://Fiddler46:Fiddler46@cluster0.um5qb.mongodb.net/cities-nighthack?retryWrites=true&w=majority")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	clientLocal, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	err = clientLocal.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("MongoDB started successfully!")
	client = *clientLocal
}

// model for user endpoint
type User struct {
	Email       string `json:"email"`
	AccessToken string `json:"token"`
}

type City struct {
	Name string `bson:"name" json:"name"`
	// CountryName string `bson:"country_name" json:"country"`
}

func main() {
	setup()
	r := mux.NewRouter()
	r.HandleFunc("/user", createUser).Methods("POST")
	r.Use(APIAuth)
	r.HandleFunc("/suggest", searchCity).Methods("GET")

	fmt.Println("Server running at port 8080")
	log.Fatal(http.ListenAndServe(":8080", r))

}

func createUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("content-type", "application/json")
	var user User
	_ = json.NewDecoder(r.Body).Decode(&user)
	userCollection := client.Database("cities-nighthack").Collection("user")
	user.AccessToken = generateToken()
	inserted, err := userCollection.InsertOne(context.Background(), user)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Fatalln(err, "Failed to create user.")
	}
	json.NewEncoder(w).Encode(inserted)

}

func searchCity(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	values := r.URL.Query()
	city := values["city_name"]
	cityCollection := client.Database("cities-nighthack").Collection("city")

	if len(city) == 0 {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"Error": "City name cannot be empty!"}`))
		return
	}

	filter := bson.D{
		primitive.E{
			Key: "all_names", Value: primitive.Regex{
				Pattern: city[0], Options: "i",
			},
		},
	}

	cursor, err := cityCollection.Find(r.Context(), filter)
	if err != nil {
		log.Fatal(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"Error": "Could not execute cursor into query."}`))
		return
	}

	var cityList []City
	for cursor.Next(context.TODO()) {
		var cities City
		err := cursor.Decode(&cities)
		if err != nil {
			log.Fatal(err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"Error": "Could not run cursor."}`))
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

}

func init() {
	rand.Seed(time.Now().UnixNano())
}

func genString(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = characterRunes[rand.Intn(len(characterRunes))]
	}
	return string(b)
}

func generateToken(n ...int) string {
	characters := 32

	if len(n) > 0 {
		characters = n[0]
	}

	randString := genString(characters)

	h := sha256.New()
	h.Write([]byte(randString))
	generatedToken := h.Sum(nil)

	return fmt.Sprintf("%x", generatedToken)

}

func APIAuth(endpoint http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		accessToken := r.Header.Get("x-api-key")

		// Check token in db
		usersCollection := client.Database("cities-nighthack").Collection("user")
		filter := bson.D{
			primitive.E{
				Key: "accesstoken", Value: accessToken,
			},
		}

		var user User
		err := usersCollection.FindOne(context.TODO(), filter).Decode(&user)
		if err != nil {
			http.Error(w, `Unauthorized access`, http.StatusUnauthorized)
			log.Fatal(err)
		} else {
			endpoint.ServeHTTP(w, r)
		}
	})
}
