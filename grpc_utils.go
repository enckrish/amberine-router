package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"io"
	"net"
	"router/pb"
	"strings"
)

const serverPort = 50051

type RouterServer struct {
	pb.UnimplementedRouterServer

	// Better to use a persistent db, e.g. Redis
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
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", serverPort))
	if err != nil {
		return err
	}
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
	id := uuid.New().String()

	if req.Id != nil {
		return nil, errors.New("id filled")
	}
	r.id2service[id] = strings.ToLower(req.Service)
	reqWithId := &pb.InitRequest_Type0{Id: &pb.UUID{Id: id}, Service: req.Service, HistorySize: req.HistorySize}

	_, err := r.prod.PublishInitRequest(reqWithId, r.id2service)
	res := pb.InitResponse_Type0{Id: &pb.UUID{Id: id}}
	return &res, err
}
