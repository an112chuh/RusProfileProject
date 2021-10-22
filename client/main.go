package main

import (
	"context"
	"fmt"

	pb "../proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"
)

type userData struct {
	INN, KPP, name, HeadName string
}

var address string = "127.0.0.1:8080"

func main() {
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		grpclog.Fatalf("fail to dial: %v", err)
	}
	defer conn.Close()
	var INN string
	INN = "4285493738"
	Client := pb.NewRusProfileServiceClient(conn)
	Request := &pb.INNRequest{
		INN: INN,
	}
	var result, err1 = Client.GetDataByINN(context.Background(), Request)
	if err1 != nil {
		grpclog.Fatalf("fail to dial: %v", err)
	} else {
		if result.Name != "" {
			fmt.Printf("\n%v\n%v\n%v\n%v\n", result.Name, result.INN, result.KPP, result.HeadName)
		} else {
			fmt.Println("Can't find a result")
		}
	}
}
