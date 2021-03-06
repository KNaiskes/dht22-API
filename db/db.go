package db

import (
	"context"
	"fmt"
	"github.com/KNaiskes/measurementsTH-API/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"time"
)

const (
	URI            = "mongodb://dht22_db:27017" // dht22_db is the name of the docker container
	TIMEOUT        = 10 * time.Second
	GET_TIMEOUT    = 30 * time.Second
	INSERT_TIMEOUT = 5 * time.Second
	DATABASE       = "measurements"
	COLLECTION     = "dht22"
)

var (
	collection  *mongo.Collection = client.Database(DATABASE).Collection(COLLECTION)
	ctx                           = context.TODO()
	client, err                   = mongo.Connect(ctx, options.Client().ApplyURI(URI))
	_, cancel                     = context.WithTimeout(ctx, TIMEOUT)
)

func Connect() {
	defer cancel()

	err = client.Ping(ctx, readpref.Primary())
	defer cancel()

	if err != nil {
		panic(err)
	} else {
		fmt.Println("Connected to db")
	}
}

func GetAll() []models.Measurement {
	results := []models.Measurement{}

	ctx, cancel = context.WithTimeout(context.Background(), GET_TIMEOUT)
	defer cancel()

	cur, err := collection.Find(ctx, bson.D{})
	if err != nil {
		fmt.Println(err)
	}

	defer cur.Close(ctx)

	for cur.Next(ctx) {
		row := models.Measurement{}

		err = cur.Decode(&row)
		if err != nil {
			fmt.Println(err)
		}
		results = append(results, row)
	}

	return results
}

func GetOne(id string) (m models.Measurement, e bool) {
	var result models.Measurement
	var empty bool = false

	ctx, cancel = context.WithTimeout(context.Background(), GET_TIMEOUT)
	defer cancel()

	cur, err := collection.Find(ctx, bson.D{})
	if err != nil {
		fmt.Println(err)
	}

	filter := bson.M{"_id": id}
	err = collection.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		fmt.Println(err)
	}

	if (models.Measurement{}) == result {
		empty = true
	}

	defer cur.Close(ctx)

	return result, empty
}

func GetAllByName(name string) []models.Measurement {
	results := []models.Measurement{}

	ctx, cancel = context.WithTimeout(context.Background(), GET_TIMEOUT)
	defer cancel()

	cur, err := collection.Find(ctx, bson.M{"name": name})
	if err != nil {
		fmt.Println(err)
	}

	defer cur.Close(ctx)

	for cur.Next(ctx) {
		row := models.Measurement{}

		err = cur.Decode(&row)
		if err != nil {
			fmt.Println(err)
		}
		results = append(results, row)
	}

	return results
}

func InsertNewMeasurement(m models.Measurement) {
	ctx, cancel = context.WithTimeout(context.Background(), INSERT_TIMEOUT)
	defer cancel()
	_, err := collection.InsertOne(ctx, m)
	if err != nil {
		fmt.Println(err)
	}
}

func NameExists(name string) bool {
	var result models.Measurement
	var exists bool = true

	ctx, cancel = context.WithTimeout(context.Background(), GET_TIMEOUT)
	defer cancel()

	cur, err := collection.Find(ctx, bson.D{})
	if err != nil {
		fmt.Println(err)
	}

	filter := bson.M{"name": name}
	err = collection.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		fmt.Println(err)
	}

	if result.ID == "" {
		exists = false
	}

	defer cur.Close(ctx)

	return exists
}
