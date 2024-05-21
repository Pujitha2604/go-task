package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ssm"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client
var parameterStore *ssm.SSM

func main() {
	var err error
	// Initialize MongoDB client
	clientOptions := options.Client().ApplyURI("mongodb+srv://admin:admin@cluster0.0elhpdy.mongodb.net/?retryWrites=true&w=majority&appName=Cluster0")
	client, err = mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	// Initialize AWS session
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))
	parameterStore = ssm.New(sess)

	// Create document in Parameter Store
	err = createDocumentInParameterStore()
	if err != nil {
		log.Fatal(err)
	}

	// Retrieve data from Parameter Store
	doc, err := getDocumentFromParameterStore()
	if err != nil {
		log.Fatal(err)
	}

	// Insert data into MongoDB
	err = insertDocumentIntoMongoDB(doc)
	if err != nil {
		log.Fatal(err)
	}

	// Start HTTP server
	fmt.Println("Server started on port 8080")
	http.HandleFunc("/", fetchData)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func fetchData(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Document fetched from Parameter Store and inserted into MongoDB successfully")
}

func createDocumentInParameterStore() error {
	input := &ssm.PutParameterInput{
		Name:      aws.String("/myapp/config1"),
		Value:     aws.String(`{"name": "abc", "age": 5, "city": "chicago"}`),
		Type:      aws.String("String"),
		Overwrite: aws.Bool(true),
	}

	_, err := parameterStore.PutParameter(input)
	return err
}

func getDocumentFromParameterStore() (map[string]interface{}, error) {
	input := &ssm.GetParameterInput{
		Name: aws.String("/myapp/config1"),
	}

	result, err := parameterStore.GetParameter(input)
	if err != nil {
		return nil, err
	}

	var doc map[string]interface{}
	err = json.Unmarshal([]byte(*result.Parameter.Value), &doc)
	if err != nil {
		return nil, err
	}

	return doc, nil
}

func insertDocumentIntoMongoDB(doc map[string]interface{}) error {
	collection := client.Database("TestingDatabase").Collection("Table1")
	_, err := collection.InsertOne(context.Background(), doc)
	if err != nil {
		return err
	}

	fmt.Println("Document inserted into MongoDB successfully")
	return nil
}
