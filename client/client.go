package main

import (
	"context"
	"distributed-auction/proto"
	"log"
	"os"

	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Host struct {
	Host string
	Port string
}

var hosts = []Host{}
	
var AuctionClient proto.AuctionServiceClient
var bidderId string

func readHostsFromFile() {
	bidderId = uuid.New().String()
	file, err := os.Open(os.Args[3])
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	
	//server ips are introduced with --Server-- and client ips with --Client--
}

func main() {
	if len(os.Args) < 4 {
		log.Fatal("Usage : go run client.go <bid> <request> <hosts file>")
	}
	
	readHostsFromFile()

	conn, err := grpc.Dial("localhost:8080", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("could not connect: %v", err)
	}
	log.Printf("Connection State: %s", conn.GetState().String())
	defer conn.Close()

	AuctionClient = proto.NewAuctionServiceClient(conn)

}

func bid(bid int32) {
	request := &proto.BidRequest{Bidder: "Bidder1", Bid: bid}
	ack, err := AuctionClient.Bid(context.Background(), request)
	if err != nil {
		log.Fatalf("Error when calling ExampleMethod: %s", err)
	}
	log.Printf("Response from server: %s", ack)
}

func request() {
	res, err := AuctionClient.Result(context.Background(), &proto.EmptyMessage{})
	if err != nil {
		log.Fatalf("Error when calling ExampleMethod: %s", err)
	}
	log.Printf("Response from server: %s", res)
}

func SendServerIsDesynchronized() {
	request := &proto.CorrectValue{Value: 1}
	ack, err := AuctionClient.ServerIsDesynchronized(context.Background(), request)
	if err != nil {
		log.Fatalf("Error when calling ExampleMethod: %s", err)
	}
	log.Printf("Response from server: %s", ack)
}