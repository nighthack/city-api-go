package main

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var client mongo.Client                                                                       // mongo client from setup declared globally
var characterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789") // character runes used for creation of tokens

func setup() { // Setup for mongodb connection using mongo drivers. Returns client.
	clientOptions := options.Client().
		ApplyURI(os.Getenv("DB_URI"))
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
}

func main() {
	err := godotenv.Load() // Loads environment variables from .env
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	setup()
	r := mux.NewRouter()
	r.Handle("/user", UserAuth(createUser)).Methods("POST")  // Route for registering new users
	r.Handle("/suggest", APIAuth(searchCity)).Methods("GET") // Route for searching cities within DB

	fmt.Println("Server running at port 8080")
	log.Fatal(http.ListenAndServe(":8080", r))

}

func createUser(w http.ResponseWriter, r *http.Request) { // POST request that takes in user details, generates an access token and inserts them into the DB
	w.Header().Add("content-type", "application/json")

	var user User
	_ = json.NewDecoder(r.Body).Decode(&user)
	userCollection := client.Database(os.Getenv("CITY_DB")).Collection(os.Getenv("COL_USER")) // Holds 'user' collection from DB
	user.AccessToken = generateToken()                                                        // Generated access token is assigned to user struct

	_, err := userCollection.InsertOne(context.Background(), user) // Details are inserted into the collection
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Fatalln(err, "Failed to create user.")
	}
	json.NewEncoder(w).Encode(user)

}

func searchCity(w http.ResponseWriter, r *http.Request) { // GET request that takes in a url param and returns a list of cities from the DB
	w.Header().Set("Content-Type", "application/json")

	values := r.URL.Query()
	city := values["city_name"]                                                               // Assigns url param from route `/suggest?city_name=` to 'city'
	cityCollection := client.Database(os.Getenv("CITY_DB")).Collection(os.Getenv("COL_CITY")) // Holds 'city' collection from DB

	if len(city) == 0 {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"Error": "City name cannot be empty!"}`))
		return
	}

	filter := bson.D{ //Regex filter for querying the input through the DB
		primitive.E{
			Key: "all_names", Value: primitive.Regex{
				Pattern: city[0], Options: "i",
			},
		},
	}

	cursor, err := cityCollection.Find(r.Context(), filter) // Query to find city is declared with cursor
	if err != nil {
		log.Fatal(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"Error": "Could not execute cursor into query."}`))
		return
	}

	var cityList []City
	for cursor.Next(context.TODO()) { // Runs through each document entry in DB to see if regex matches
		var cities City
		err := cursor.Decode(&cities)
		if err != nil {
			log.Fatal(err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"Error": "Could not run cursor."}`))
			return
		}
		cityList = append(cityList, cities) // List of cities is appended to cityList
	}
	w.WriteHeader(http.StatusOK)
	if len(cityList) == 0 {
		w.Write([]byte(`{"cities": []}`))
		return
	}
	json.NewEncoder(w).Encode(cityList)

}

func init() {
	rand.Seed(time.Now().UnixNano()) // Creates a new seed for randomising the access token
}

func genString(n int) string { // Generates a randomised string from the character runes
	b := make([]rune, n)
	for i := range b {
		b[i] = characterRunes[rand.Intn(len(characterRunes))]
	}
	return string(b)
}

func generateToken(n ...int) string { // Generates a token for when a new user registers their details
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

func APIAuth(endpoint func(http.ResponseWriter, *http.Request)) http.Handler { // Middleware to check if access token is contained within header when running searchCity
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		accessToken := r.Header.Get("x-api-key")

		// Checks token in db
		usersCollection := client.Database(os.Getenv("CITY_DB")).Collection(os.Getenv("COL_USER"))
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
			endpoint(w, r)
		}
	})
}

func UserAuth(endpoint func(http.ResponseWriter, *http.Request)) http.Handler { // Middleware to check if user registering has appropriate auth token
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		accessToken := r.Header.Get("x-access-token")
		if accessToken != os.Getenv("USER_AUTH_TOKEN") {
			http.Error(w, "Unauthorized access", http.StatusUnauthorized)
		} else {
			endpoint(w, r)
		}
	})
}
