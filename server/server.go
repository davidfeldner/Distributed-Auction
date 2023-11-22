package main

import (
	"context"
	"distributed-auction/proto"
	"log"
	"net"
	"os"
	"strconv"

	"google.golang.org/grpc"
)
var logger *log.Logger
var port int
//var ip string

var is_over bool 
var highest_bid int32
var highest_bidder string


func main() {
	f, err := os.OpenFile("logfile", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	if len(os.Args) < 2 {
		log.Fatal("Usage : go run server.go <port>")
	}
	//ip = os.Args[1]
	port, err = strconv.Atoi(os.Args[1])
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


type AuctionServer struct {
	proto.UnimplementedAuctionServiceServer
}

func (s *AuctionServer) Bid(ctx context.Context, in *proto.BidRequest) (*proto.Ack, error) {
	if is_over {
		return &proto.Ack{Message: "exception"}, nil
	}
	if in.Bid <= highest_bid {
		return &proto.Ack{Message: "fail"}, nil
	}
	highest_bid = in.Bid
	highest_bidder = in.Bidder
	return &proto.Ack{Message: "success"}, nil
}

func (s *AuctionServer) Result(ctx context.Context, in *proto.EmptyMessage) (*proto.ResultRequest, error) {
	return &proto.ResultRequest{IsOver: is_over, HighestBid: highest_bid, HighestBidder: highest_bidder}, nil
}

func (s *AuctionServer) ServerIsDesynchronized(ctx context.Context, in *proto.ResultRequest) (*proto.Ack, error) {
	if in.HighestBid < highest_bid{
		return &proto.Ack{Message: "success"}, nil	
	}
	if in.HighestBid == highest_bid{
		if in.HighestBidder == highest_bidder{
			return &proto.Ack{Message: "success"}, nil	
		}else{
			is_over = in.IsOver
			highest_bid = in.HighestBid
			highest_bidder = in.HighestBidder	
		}
		return &proto.Ack{Message: "success"}, nil	
	}
	if in.HighestBid > highest_bid{
		is_over = in.IsOver
		highest_bid = in.HighestBid
		highest_bidder = in.HighestBidder
		return &proto.Ack{Message: "success"}, nil	
	}
	return &proto.Ack{Message: "fail"}, nil
}