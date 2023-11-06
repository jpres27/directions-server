package main

import (
	"context"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Directions struct {
	ID          primitive.ObjectID `bson:"id,omitempty"`
	Destination string             `json:"destination" bson:"destination"`
	How         string             `json:"how" bson:"how"`
}

func main() {
	fmt.Println("Server starting...")

	//Connect to database
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}
	uri := os.Getenv("MONGODB_URI")
	if uri == "" {
		log.Fatal("You must set your 'MONGODB_URI' environment variable. See\n\t https://www.mongodb.com/docs/drivers/go/current/usage-examples/#environment-variable")
	}
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()
	fmt.Println("Connected to database...")

	// Access database
	coll := client.Database("directions").Collection("mainlib")

	// request handler which displays site
	h1 := func(w http.ResponseWriter, r *http.Request) {
		tmpl := template.Must(template.ParseFiles("index.html")) // fix this if needed after reimplementing frontend
		filter := bson.D{{}}
		cursor, err := coll.Find(context.TODO(), filter)
		if err != nil {
			panic(err)
		}

		var allDirections []bson.M
		if err = cursor.All(context.TODO(), &allDirections); err != nil {
			log.Fatal(err)
		}
		for _, allDirections := range allDirections {
			fmt.Println(allDirections)
		}
		tmpl.Execute(w, allDirections)
	}

	// request handler which accepts submission and adds to database
	h2 := func(w http.ResponseWriter, r *http.Request) {
		destination := r.PostFormValue("destination")
		how := r.PostFormValue("how")
		tmpl := template.Must(template.ParseFiles("index.html"))
		tmpl.ExecuteTemplate(w, "directions-list-element", Directions{Destination: destination, How: how})
		newDirection := Directions{Destination: destination, How: how}
		result, err := coll.InsertOne(context.TODO(), newDirection)
		if err != nil {
			panic(err)
		}
		fmt.Printf("Document inserted with ID: %s\n", result.InsertedID)
	}

	http.HandleFunc("/", h1)
	http.HandleFunc("/submit/", h2)

	log.Fatal(http.ListenAndServe(":8000", nil))
}
