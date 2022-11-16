package main

import (
	"context"
	"log"
	"os"
	"strconv"
	"sync"
	"fmt"
	"net"
	request "github.com/lucasfth/go-ass5/grpc"
	"golang.org/x/net/context/ctxhttp"
	"google.golang.org/grpc"
)

func main() {
	log.SetFlags(0)

	serverId, _ := strconv.ParseInt(os.Args[1], 10, 32)
	crashBoolean, _ := strconv.ParseBool(os.Args[2])
	ownPort := int32(serverId) + 5000

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	s := &server{
		ownPort: 		ownPort,
		currentBid: 	0,
		isOver: 		false,
		clients: 		make(map[int32]request.ClientHandshake),
		ctx: 			ctx,
	}

	// Create listener on port ownPort
	list, err := net.Listen("tcp", fmt.Sprintf(":%v", ownPort))
	if err != nil {
		log.Fatalf("Failed to listen on port: %v", err)
	}
	grpcServer := grpc.NewServer()
	request.RegisterBiddingServiceServer(grpcServer, s)
	
	go func() {
		if err := grpcServer.Serve(list); err != nil {
			log.Fatalf("failed to server %v", err)
		}
	}()

	for {

	}
}

func (s *server) WelcomeClient(ctx context.Context, req *request.ClientHandshake) (*request.BidResponse, error) {
	
	
	
}

type server struct{
	mutex 		sync.Mutex
	ownPort 	int32
	currentBid 	int32
	isOver 		bool
	clients 	map[int32]request.ClientHandshake
	ctx 		context.Context
}