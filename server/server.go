package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"sync"
	"time"

	request "github.com/lucasfth/go-ass5/grpc"
	"google.golang.org/grpc"
)

func main() {
	var id int32
	log.Printf("Enter id below:")
	fmt.Scanln(&id)
	log.Printf("Welcome %v", id)

	var endHour, endMin int
	log.Printf("Enter auction time below in hour and min:")
	fmt.Scanln(&endHour, &endMin)
	end := time.Date (time.Now().Year(), time.Now().Month(), time.Now().Day(), endHour, endMin, 0, 0, time.Local)

	port := int32(5000 + id)
	portString := fmt.Sprintf(":%v", port)
	lis, err := net.Listen("tcp", portString)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	
	server := &server{
		mutex: sync.Mutex{},
		ownPort: port,
		currentBid: 0,
		currentBidOwner: "",
		isOver: false, 
		auctionEnd: end,
		clients: make(map[int32]request.ClientHandshake),
		ctx: context.Background(),
	}
	
	// create grpc server
	s := grpc.NewServer()
	request.RegisterBiddingServiceServer(s, server)
	
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func(s *server) Handshake(in *request.ClientHandshake, srv request.BiddingService_HandshakeServer) error {
	log.Printf("Handshake 	%s", in.Name)

	resp := &request.BidResponse{}
	resp.Response = "Succes"
	srv.Send(resp);	
	return nil;
}

func (s *server) SendBid(in *request.Bid, srv request.BiddingService_SendBidServer) error{
	s.mutex.Lock()
	defer s.mutex.Unlock()

	resp := &request.BidResponse{}

	if (time.Until(s.auctionEnd) <= 0) {
		s.isOver = true
	}
	if (s.isOver) {
		log.Printf("Bid 	%s with %v but auction over, winner: %s , with: %v", in.Name, in.Amount, s.currentBidOwner, s.currentBid)
		resp.Response = "Fail" 
		srv.Send(resp)
		return nil
	}

	
	if in.Amount > s.currentBid{
		s.currentBid = in.Amount
		s.currentBidOwner = in.Name
		resp.Response = "Success"
	} else if in.Amount <= s.currentBid{
		resp.Response = "Fail"
	} else {
		resp.Response = "Exception"
	}

	log.Printf("Bid 	%s	%s with %v", in.Name, resp.Response, in.Amount)
	
	srv.Send(resp)
	return nil
}

func (s *server) RequestCurrentResult(in *request.Request, srv request.BiddingService_RequestCurrentResultServer) error {
	log.Printf("Request 	%s	highest bid is: %v ,by: %s", in.Name, s.currentBid, s.currentBidOwner)

	if (time.Until(s.auctionEnd) <= 0) {
		s.isOver = true
	}
	
	resp := &request.RequestResponse{}
	resp.HighestBid = s.currentBid
	resp.IsOver = s.isOver
	if s.isOver {
		resp.WinnerName = s.currentBidOwner
	} else {
		resp.WinnerName = ""
	}
	srv.Send(resp);
	return nil;
}

type server struct{
	mutex 			sync.Mutex
	ownPort 		int32
	currentBid 		int32
	currentBidOwner	string
	isOver 			bool
	auctionEnd 		time.Time
	ctx 		context.Context
	request.UnimplementedBiddingServiceServer
	clients 	map[int32]request.ClientHandshake // can probably be removed
}