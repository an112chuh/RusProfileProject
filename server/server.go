package main

import (
	"context"
	"fmt"
	"net"
	"net/http"

	"github.com/PuerkitoBio/goquery"
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
	collector.OnHTML("clip_inn", func(e *colly.HTMLElement) {
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
	r, err := http.Get(FinalURL)
	if err != nil {
		fmt.Println(err)
	}
	defer r.Body.Close()
	if r.StatusCode != 200 {
		fmt.Println("status code error: %d %s", r.StatusCode, r.Status)
	}
	doc, err := goquery.NewDocumentFromReader(r.Body)
	if err != nil {
		fmt.Println(err)
	}
	doc.Find("span#clip_inn").Each(func(i int, s *goquery.Selection) {
		title := s.Text()
		response.INN = title
		fmt.Println(response.INN)
	})
	doc.Find("span#clip_kpp").Each(func(i int, s *goquery.Selection) {
		title := s.Text()
		response.KPP = title
		fmt.Println(response.KPP)
	})
	doc.Find("div.company-name").Each(func(i int, s *goquery.Selection) {
		title := s.Text()
		response.Name = title
		fmt.Println(response.Name)
	})
	doc.Find("span.company-info__text").EachWithBreak(func(i int, s *goquery.Selection) bool {
		title := s.Text()
		response.HeadName = title
		fmt.Println(response.HeadName)
		return false
	})
	return response, nil
}
