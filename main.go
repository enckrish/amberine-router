package main

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"router/pb"
	"strings"
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
	testSetup(server, []string{"ABC", "Docker"})
	err := server.Start()
	if err != nil {
		panic(err)
	}
}

func testSetup(server *RouterServer, services []string) {
	const groupId = "general-analyzers"
	const selfId = "general_analyzers-001"
	_, err := server.Admin_setNewGroupAnalyzer(context.Background(), &pb.SetNewGroupAnalyzerReq{GroupId: groupId, NewSelfId: selfId})
	if err != nil {
		panic(err)
	}
	for _, s := range services {
		_, err := server.Admin_setTargets(
			context.Background(),
			&pb.SetTargetsRequest{
				Service: strings.ToLower(s),
				Targets: []string{groupId},
			})
		if err != nil {
			panic(err)
		}
	}
}
