package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"math/rand"
	"time"

	request "github.com/lucasfth/go-ass5/grpc"

	"google.golang.org/grpc"
)

type client struct {
	name 			string
	currentBid		int32
	servers 		[]request.BiddingServiceClient
}

func main(){
	log.SetFlags(0)
	ctx := context.Background()
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithInsecure())

	c := &client{}

	log.Printf("Enter username below:")
	fmt.Scanln(&c.name)
	log.Printf("Welcome %s", c.name)

	for i := 0; i < 3; i++ { // Will iterate through ports 5001, 5002, 5003
		dialNum  := int32(5001 + i)
		dialNumString := fmt.Sprintf(":%v", dialNum) 

		conn, err := grpc.Dial(dialNumString, grpc.WithInsecure())
		if err != nil {
			log.Fatal(err)
		}
		
		// create stream
		client := request.NewBiddingServiceClient(conn)
		in := &request.ClientHandshake{ClientPort: dialNum} 
		//bidStream, err := client.SendBid(context.Background(), )
		stream, err := client.Handshake(ctx, in)
		if err != nil {
			log.Fatalf("open stream error %v", err)
		}
		resp, err := stream.Recv()
		if err == io.EOF {
			log.Printf("Connected to server %v and responded %s", dialNum, resp)
		}
		if err != nil {
			log.Fatalf("Cannot receive %v", err)
		}
		c.servers = append(c.servers, client)
		time.Sleep(4 * time.Second)
	}

	for { // Communication loop
		actionType := int32(rand.Intn(2)) // 0 = bid, 1 = request

		if actionType == 0 {
			log.Printf("Action type: Bid")
			randomBid := int32(rand.Intn(1000))
			log.Printf("Will bid with %v", randomBid)
			c.sendBids(randomBid)
			time.Sleep(4 * time.Second)
		} else {
			log.Printf("Action type: Request")
			c.requestCurrentResults()
			time.Sleep(4 * time.Second)
		}
	}
}

func (c *client) sendBids(bid int32){
	responses := make([]string, len(c.servers))
	for i := 0; i < len(c.servers); i++ { // Send bid to all servers
		response, _ := c.sendBid(int32(i), bid)
		responses[i] = response
	}
	logicResponse := c.logic(responses, bid)

	log.Printf("---------Bid %v was %s", bid, logicResponse)
}

func (c *client) sendBid(iteration int32, bid int32) (string, error) {
	in := &request.Bid{Name: c.name, Amount: bid}
	stream, err := c.servers[iteration].SendBid(context.Background(), in)
	if err != nil {
		return "nil", err
	}
	resp, err := stream.Recv()
	return resp.GetResponse(), err
}

func (c *client) requestCurrentResults() (currentRelaventBid int32){
	var highestBid int32; 
	for i := 0; i < len(c.servers); i++ { // Request current result from all servers
		resp, err := c.requestCurrentResult(int32(i))
		if err != nil {
			return
		}
		highestBid = resp.HighestBid
	}
	log.Printf("---------Current highest bid is %v", highestBid)
	return highestBid
}

func (c *client) requestCurrentResult(iteration int32)(*request.RequestResponse, error){
	in := &request.Request{}
	stream, err := c.servers[iteration].RequestCurrentResult(context.Background(), in)
	if err != nil {
		return nil, err
	}
	resp, err := stream.Recv()
	return resp, err
}

func (c *client) logic(responses []string, bid int32) (string) {
	for i := 0; i < len(responses); i++ {
		log.Printf("Response was: %s ,on i: %v", responses[i], i)
		if responses[i] == "Success" {
			c.currentBid = bid
			log.Printf("Went into success")
			return "Succes"
		} else if responses[i] == "Fail" {
			c.currentBid = -1
			log.Printf("Went into fail")
			return "Fail"
		} else if responses[i] == "Exception" {
			continue
		}
	}
	c.currentBid = -1
	return "Fail"
}