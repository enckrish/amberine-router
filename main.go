package main

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
)

var routerDb *mongo.Database

func main() {
	routerDb = getRouterDb()
	defer func() {
		if err := routerDb.Client().Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()
	server := NewRouterServer()
	err := server.Start()
	if err != nil {
		panic(err)
	}
}
