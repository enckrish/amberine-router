package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"io"
	"net"
	"router/pb"
	"strings"
)

const grpcPort = 50051

type RouterServer struct {
	pb.UnimplementedRouterServer

	// TODO replace with a persistent db, e.g. Redis
	id2service map[string]string
	prod       *KafkaProducer
}

func NewRouterServer() *RouterServer {
	r := RouterServer{}
	r.id2service = make(map[string]string)
	r.prod = NewKafkaProducer(true)
	return &r
}

func (r *RouterServer) Start() error {
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", grpcPort))
	if err != nil {
		return err
	}
	// TODO add authentication
	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)
	pb.RegisterRouterServer(grpcServer, r)
	reflection.Register(grpcServer)

	return grpcServer.Serve(lis)
}
func (r *RouterServer) RouteLog_Type0(stream pb.Router_RouteLog_Type0Server) error {
	for {
		in, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}

		_, err = r.prod.PublishAnalysisRequest(in, r.id2service)
		if err != nil {
			return err
		}
		err = stream.Send(&pb.AnalyzerResponse{Committed: true})
		if err != nil {
			return err
		}
	}
}

func (r *RouterServer) Init_Type0(
	_ context.Context,
	req *pb.InitRequest_Type0,
) (*pb.InitResponse_Type0, error) {
	fmt.Printf("Initialize Request: %+v\n", req)
	id := uuid.New().String()
	if req.StreamId != "" {
		return nil, errors.New("id filled")
	}
	r.id2service[id] = strings.ToLower(req.Service)
	reqWithId := &pb.InitRequest_Type0{StreamId: id, Service: req.Service, HistorySize: req.HistorySize}

	_, err := r.prod.PublishInitRequest(reqWithId, r.id2service)
	if err != nil {
		panic(err)
	}
	res := pb.InitResponse_Type0{StreamId: id}
	return &res, err
}

func (r *RouterServer) Admin_setTargets(
	ctx context.Context,
	req *pb.SetTargetsRequest,
) (*pb.BoolResult, error) {
	coll := routerDb.Collection(targetMapCollName)

	req.Service = strings.ToLower(req.Service)
	filter := bson.D{{"service", req.Service}}
	update := bson.M{"$set": bson.M{"targets": req.Targets}}
	upsert := true
	opts := options.UpdateOptions{Upsert: &upsert}
	_, err := coll.UpdateOne(ctx, filter, update, &opts)
	if err != nil {
		return &pb.BoolResult{Done: false}, err
	}

	//fmt.Println("Target update result:", result)
	return &pb.BoolResult{Done: true}, nil
}

func (r *RouterServer) Admin_setNewGroupAnalyzer(
	ctx context.Context,
	req *pb.SetNewGroupAnalyzerReq,
) (*pb.BoolResult, error) {
	coll := routerDb.Collection(groupAnalyzerCollName)

	filter := bson.D{{"group_id", req.GroupId}}
	update := bson.M{"$set": bson.M{"leader": req.NewSelfId}}
	upsert := true
	opts := options.UpdateOptions{Upsert: &upsert}
	_, err := coll.UpdateOne(ctx, filter, update, &opts)
	if err != nil {
		return &pb.BoolResult{Done: false}, err
	}

	//fmt.Println("Leader update result:", result)
	return &pb.BoolResult{Done: true}, nil
}
