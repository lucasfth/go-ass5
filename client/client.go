package main

import (
	"context"
	"fmt"
	"log"

	request "github.com/lucasfth/go-ass5/grpc"

	"google.golang.org/grpc"
)

type client struct {
	name string
}

func main(){
	log.SetFlags(0)
	ctx := context.Background()
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithInsecure())

	conn, err := grpc.Dial(":8080", opts...)
	if err != nil {
		log.Fatal(err)
	}

	client := request.NewBiddingServiceClient(conn)

}

func sendBid(){

}

func RequestCurrentResult(){

}