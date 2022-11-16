package main

import (
	"sync"
	request "lucasfth/go-ass5/grpc"
	"context"
)

type serve struct{
	mutex sync.Mutex
	currentBid int32
	isOver bool
	clients map[int32]request.Client
}