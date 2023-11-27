package main

import (
	"context"
	"distributed-auction/proto"
	"log"
	"net"
	"os"
	"strconv"
	"sync"
	"time"

	"google.golang.org/grpc"
)
var logger *log.Logger
var port int
//var ip string
var mu = &sync.Mutex{}

var is_over bool 
var highest_bid int32
var highest_bidder string
var auctionEndsAt time.Time

func main() {
	f, err := os.OpenFile("logfile", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	// new log with normal layout + "SERVER":port, logs to file called logfile
	logger = log.New(os.Stdout, "SERVER:"+strconv.Itoa(port)+" ", log.LstdFlags)
	logger.SetOutput(f)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	if len(os.Args) < 3 {
		log.Fatal("Usage : go run server.go <port> <time auction ends>")
	}
	location, err := time.LoadLocation("Europe/Copenhagen")
	if err != nil {
		log.Fatal(err)
	}
	 
	//ip = os.Args[1]
	port, err = strconv.Atoi(os.Args[1])
	if err != nil {
		log.Fatal("Usage : go run server.go <port> <time auction ends>")
	}
	auctionEndsAt, err = time.ParseInLocation("2006-01-02 15:04:05", os.Args[2], location)
	go isAuctionOver()
	if err != nil {
		log.Fatal("Usage : go run server.go <port>")
	}
	// Make a log file
	logger = log.New(os.Stdout, "PORT: " +   strconv.Itoa(port) + " ", log.LstdFlags)
	logger.SetOutput(f)

	grpcServer := grpc.NewServer()
	proto.RegisterAuctionServiceServer(grpcServer, &AuctionServer{})

	listener, err := net.Listen("tcp", "localhost:" + strconv.Itoa(port))
	if err != nil {
		panic(err)
	}

	if err := grpcServer.Serve(listener); err != nil {
		panic(err)
	}
}

func isAuctionOver() {
	logger.Println("Auction ends at: ", auctionEndsAt)
	logger.Println("until: ", time.Until(auctionEndsAt).Seconds())
	if auctionEndsAt.Before(time.Now()) {
		is_over = true
		return
	}
	time.Sleep(time.Until(auctionEndsAt))
	is_over = true
	logger.Println("Auction is over")
}

type AuctionServer struct {
	proto.UnimplementedAuctionServiceServer
}

func (s *AuctionServer) Bid(ctx context.Context, in *proto.BidRequest) (*proto.Ack, error) {
	mu.Lock()
	defer mu.Unlock()
	if is_over {
		logger.Println("Auction is over")
		return &proto.Ack{Message: "exception"}, nil
	}
	if in.Bid <= highest_bid {
		logger.Printf("Client %s tried to bid %d, but highest bid is %d", in.Bidder, in.Bid, highest_bid)
		return &proto.Ack{Message: "fail"}, nil
	}
	highest_bid = in.Bid
	highest_bidder = in.Bidder
	logger.Printf("Client %s bid %d", in.Bidder, in.Bid)
	return &proto.Ack{Message: "success"}, nil
}

func (s *AuctionServer) Result(ctx context.Context, in *proto.EmptyMessage) (*proto.ResultResponse, error) {
	return &proto.ResultResponse{IsOver: is_over, HighestBid: highest_bid, HighestBidder: highest_bidder}, nil
}

func (s *AuctionServer) ServerIsDesynchronized(ctx context.Context, in *proto.ResultResponse) (*proto.Ack, error) {
	mu.Lock()
	defer mu.Unlock()
	logger.Printf("Server has been made aware that it is desynchronized, attempting fix")
	is_over = in.IsOver
	highest_bid = in.HighestBid
	highest_bidder = in.HighestBidder
	return &proto.Ack{Message: "success"}, nil	
}