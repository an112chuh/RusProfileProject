package main

import (
	"context"
	"net"

	pb "../proto"

	"github.com/gocolly/colly"
	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"
)

func main() {
	listener, err := net.Listen("tcp", ":5300")

	if err != nil {
		grpclog.Fatalf("failed to listen: %v", err)
	}
	opts := []grpc.ServerOption{}
	GrpcServer := grpc.NewServer(opts...)
	server := pb.RusProfileServiceServer{}
	pb.RegisterRusProfileServiceServer(GrpcServer, &server{})
	GrpcServer.Serve(listener)
}

const URL = "https://www.rusprofile.ru/search?query="

type server struct{}

func (s *server) GetDataByINN(c context.Context, request *pb.INNRequest) (*pb.UserData, error) {
	var response *pb.UserData
	response = new(pb.UserData)
	var FinalURL string
	FinalURL = URL + request.INN
	collector := colly.NewCollector(
		colly.AllowedDomains(FinalURL),
	)
	collector.OnHTML("span#clip_inn", func(e *colly.HTMLElement) {
		response.INN = e.Text
	})
	collector.OnHTML("span#clip_kpp", func(e *colly.HTMLElement) {
		response.KPP = e.Text
	})
	collector.OnHTML("company-name", func(e *colly.HTMLElement) {
		response.Name = e.Text
	})
	collector.OnHTML("span#company-info__text", func(e *colly.HTMLElement) {
		response.HeadName = e.Text
	})
	return response, nil
}
