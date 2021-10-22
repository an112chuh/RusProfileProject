package main

import (
	"context"
	"fmt"
	"net"

	pb "github.com/an112chuh/rusprofileproject/proto"

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
	//	var server pb.&RusProfileServiceServer;
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
	fmt.Println(FinalURL)
	collector := colly.NewCollector()
	collector.OnHTML("span#clip_inn", func(e *colly.HTMLElement) {
		response.INN = e.Text
		fmt.Println(response.INN)
	})
	collector.OnHTML("span#clip_kpp", func(e *colly.HTMLElement) {
		response.KPP = e.Text
		fmt.Println(response.KPP)
	})
	collector.OnHTML("company-name", func(e *colly.HTMLElement) {
		response.Name = e.Text
		fmt.Println(response.Name)
	})
	collector.OnHTML("span#company-info__text", func(e *colly.HTMLElement) {
		response.HeadName = e.Text
		fmt.Println(response.HeadName)
	})
	return response, nil
}
