package main

import (
	"context"
	"fmt"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"os"
)

const mongoPwdEnv = "MONGO_PASSWORD"
const routingDbName = "routing-db"
const targetMapCollName = "target-map"
const groupAnalyzerCollName = "group-leader"

func getRouterDb() *mongo.Database {
	_ = godotenv.Load()
	pwd := os.Getenv(mongoPwdEnv)

	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	uri := fmt.Sprintf(
		"mongodb+srv://ksengupta2911:%s@amberine-mumbai.g57yunl.mongodb.net/?retryWrites=true&w=majority",
		pwd)
	opts := options.Client().ApplyURI(uri).SetServerAPIOptions(serverAPI)

	client, err := mongo.Connect(context.TODO(), opts)
	if err != nil {
		panic(err)
	}

	db := client.Database(routingDbName)

	if err := db.RunCommand(context.TODO(), bson.D{{"ping", 1}}).Err(); err != nil {
		panic(err)
	}
	fmt.Println("Pinged your deployment. You successfully connected to MongoDB!")

	return db
}
