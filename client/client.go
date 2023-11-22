package main

import (
	"bufio"
	"context"
	"distributed-auction/proto"
	"errors"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Host struct {
	Host string
	Port string
}
var logger *log.Logger

var hosts = []Host{}
	
var AuctionHosts []proto.AuctionServiceClient
var bidderId string
var maxBid int32

func readHostsFromFile() {
	bidderId = uuid.New().String()
	file, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		hostString := scanner.Text()
		var hostSlice = strings.Split(hostString, ":")
		hosts = append(hosts, Host{Host: hostSlice[0], Port: hostSlice[1]})
	}
}

func resultRequestIsEqual(a *proto.ResultRequest, b *proto.ResultRequest) bool {
	return a.IsOver == b.IsOver && a.HighestBid == b.HighestBid && a.HighestBidder == b.HighestBidder
}

func createConnetionToServers(server Host) {
	conn, err := grpc.Dial(server.Host + ":" + server.Port, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Printf("could not connect: %v", err)
	}
	log.Printf("Connection State: %s", conn.GetState().String())
	//defer conn.Close()

	NewAuctionClient := proto.NewAuctionServiceClient(conn)
	AuctionHosts = append(AuctionHosts, NewAuctionClient)
}

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Usage : go run client.go <hosts file>")
	}
	f, err := os.OpenFile("logfile", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	bidderId = uuid.New().String()
	maxBid = rand.Int31n(10000)
	// Make a log file
	logger = log.New(os.Stdout, "UUID: " + bidderId + " ", log.LstdFlags)
	logger.SetOutput(f)
	
	readHostsFromFile()
	for _, host := range hosts {
		createConnetionToServers(host)
	}
	
	for {
		time.Sleep(time.Duration(rand.Intn(4)+1) * time.Second)
		results, err := ResultAll()
		if err != nil {
			log.Fatalf("Error when getting results: %s", err)
		}
		mostCommon := getMostCommonResult(results)
		log.Printf("Auction is over: %t", mostCommon.IsOver)
		log.Printf("Highest bidder is %s", mostCommon.HighestBidder)
		log.Printf("Highest bid is %d", mostCommon.HighestBid)
		for index, result := range results {
			if !resultRequestIsEqual(result, mostCommon) {
				SendServerIsDesynchronized(mostCommon, AuctionHosts[index])
			}
		}
		if mostCommon.IsOver {
			log.Printf("Auction is over!\nHighest bidder is %s with bid %d", mostCommon.HighestBidder, mostCommon.HighestBid)
			break
		}
		if mostCommon.HighestBidder == bidderId {
			log.Printf("I am the highest bidder!")
			continue
		}
		bid, err := calculateNextBid(mostCommon.HighestBid)
		if err != nil {
			if err.Error() == "max bid reached" {
				log.Printf("Max bid reached: %d", maxBid)
			} else {
				log.Fatalf("Error when calculating next bid: %s", err)		
			}
		}
		AckList := bidAll(bid)
		checkforAckError(AckList);
		
	}
}
func getMostCommonResult(results []*proto.ResultRequest) *proto.ResultRequest {
	resultMap := make(map[*proto.ResultRequest]int32)
	for _, result := range results {
		resultMap[result]++
		
	}
	var max int32
	var maxResult *proto.ResultRequest
	for res, value := range resultMap {
		if value > max {
			max = value
			maxResult = res
		}
	}
	return maxResult
}

func checkforAckError(acks []*proto.Ack) *proto.Ack {
	ackMap := make(map[*proto.Ack]int32)
	for _, ack := range acks {
		ackMap[ack]++
	}
	var max int32
	var maxAck *proto.Ack

	for ack, value := range ackMap{ 
		if value > max {
			max = value
			maxAck = ack
		}
	}
	return maxAck
}

func calculateNextBid(currentBid int32) (int32, error) {
	if currentBid >= maxBid { 
		return 0, errors.New("max bid reached")
	}
	multiplier := rand.Int31n(10)+1
	var proposedBid = currentBid+1 * (1+multiplier/100)
	if(proposedBid > maxBid) {
		return maxBid, nil
	}
	return proposedBid, nil
}

func bidAll(bid int32) ([]*proto.Ack) {
	results := []*proto.Ack{}
	for _, ac := range AuctionHosts {
		results = append(results, Bid(bid, ac))
	}
	return results
}

func Bid(amount int32, ac proto.AuctionServiceClient) *proto.Ack {
	request := &proto.BidRequest{Bidder: bidderId, Bid: amount}
	res, err := ac.Bid(context.Background(), request)
	log.Printf("Bid: %d", amount)
	if err != nil {
		log.Printf("Error when calling Bid: %s", err)
	}
	log.Printf("Bid-response from server: %s", res)
	return res
}

func ResultAll() ([]*proto.ResultRequest, error) {
	results := []*proto.ResultRequest{}
	for _, ac := range AuctionHosts {
		res, err := Result(ac)
		if err == nil {
			results = append(results, res)
		}
	}
	return results, nil
}

func Result(ac proto.AuctionServiceClient) (*proto.ResultRequest, error) {
	res, err := ac.Result(context.Background(), &proto.EmptyMessage{})
	if err != nil {
		return nil, err
	}
	log.Printf("Result-response from server: %s", res)
	return res, nil
}

func SendServerIsDesynchronized(correctValue *proto.ResultRequest, ac proto.AuctionServiceClient) {
	ack, err := ac.ServerIsDesynchronized(context.Background(), correctValue)
	if err != nil {
		log.Printf("Error when calling ServerIsDesynchronized: %s", err)
	}
	log.Printf("SendServerIsDesynchronized-response from server: %s", ack.Message)
}

