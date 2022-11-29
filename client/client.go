package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"math/rand"
	"os"
	"time"

	request "github.com/lucasfth/go-ass5/grpc"

	"google.golang.org/grpc"
)

type client struct {
	name 			string
	downedServers 	map[int32]bool
	downedServerInt int32
	servers 		[]request.BiddingServiceClient
}

func main(){
	log.SetFlags(0)
	ctx := context.Background()
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithInsecure())

	c := &client{downedServerInt: 0}

	log.Printf("Enter username below:")
	fmt.Scanln(&c.name)
	log.Printf("Welcome %s", c.name)
	
	c.downedServers = make (map[int32]bool)

	for i := 0; i < 3; i++ { // Will iterate through ports 5001, 5002, 5003
		dialNum  := int32(5001 + i)
		dialNumString := fmt.Sprintf(":%v", dialNum) 

		conn, err := grpc.Dial(dialNumString, grpc.WithInsecure())
		if err != nil {
			log.Fatal(err)
		}
		
		// create stream
		client := request.NewBiddingServiceClient(conn)
		in := &request.ClientHandshake{ClientPort: dialNum, Name: c.name} 
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
		c.downedServers[dialNum] = false
		time.Sleep(4 * time.Second)
	}

	for { // Communication loop
		rand.Seed(time.Now().UnixNano()) // ensure "random" number is different each time
		actionType := int32(rand.Intn(2)) // 0 = bid, 1 = request

		if actionType == 0 {
			// log.Printf("Action type: Bid")
			randomBid := int32(rand.Intn(1000))
			c.sendBids(randomBid)
			time.Sleep(4 * time.Second)
		} else {
			// log.Printf("Action type: Request")
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
	
	if logicResponse == "Exception" {
		log.Print("Tried to bid but found EXCEPTION")
		os.Exit(1)
	}

	log.Printf("---------Bid %v was %s", bid, logicResponse)
}

func (c *client) sendBid(iteration int32, bid int32) (string, error) {
	in := &request.Bid{Name: c.name, Amount: bid}
	stream, err := c.servers[iteration].SendBid(context.Background(), in)
	if err != nil {
		serverDown(iteration, c)
		return "nil", err
	}
	resp, err := stream.Recv()
	return resp.GetResponse(), err
}

func serverDown (iteration int32, c *client) {
	if (!c.downedServers[5001 + iteration]) {
		log.Printf("Server %v is down", (5001 + iteration))
	}
	c.downedServers[5001 + iteration] = true
}

func (c *client) requestCurrentResults() (currentRelaventBid int32){
	var highestBid int32 
	for i := 0; i < len(c.servers); i++ { // Request current result from all servers
		resp, err := c.requestCurrentResult(int32(i))
		if err != nil {
			c.downedServerInt++
			serverDown(int32(i), c)
			continue
		}
		if (c.downedServerInt == int32(len(c.servers))){
			log.Printf("Tried to request but found EXCEPTION")
			os.Exit(1)
		}
		if (resp.IsOver) {
			auctionFinished(resp, c.name)
		}
		highestBid = resp.HighestBid
	}
	log.Printf("---------Current highest bid is %v", highestBid)
	return highestBid
}

func auctionFinished(resp *request.RequestResponse, name string) {
	if (resp.WinnerName == name) {
		log.Printf("Won the auction with bid %v", resp.HighestBid)
	} else {
		log.Printf("Lost the auction")
	}
	os.Exit(1)
}

func (c *client) requestCurrentResult(iteration int32)(*request.RequestResponse, error){
	in := &request.Request{Name: c.name}
	stream, err := c.servers[iteration].RequestCurrentResult(context.Background(), in)
	if err != nil {
		return nil, err
	}
	resp, err := stream.Recv()
	return resp, err
}

func (c *client) logic(responses []string, bid int32) (string) {
	for i := 0; i < len(responses); i++ {
		// log.Printf("Response was: %s ,on i: %v", responses[i], i)
		if responses[i] == "Success" {
			// log.Printf("Went into success")
			return "Succes"
		} else if responses[i] == "Fail" {
			// log.Printf("Went into fail")
			return "Fail"
		} else if (i == len(responses) - 1) {
			return "Exception"
		}
	}
	return "Fail"
}