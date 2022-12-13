package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	pb "github.com/toposware/go-topos-sequencer-client/client/proto"
	"google.golang.org/grpc"
)

func main() {
	fmt.Println("Starting go client")
	url := "127.0.0.1:4001"

	conn, err := grpc.Dial(url, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	// create stream
	client := pb.NewFrostAPIServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	stream, err := client.WatchFrostMessages(ctx)
	if err != nil {
		log.Fatalf("openn stream error %v", err)
	}

	waitc := make(chan struct{})
	go func() {
		fmt.Println("Starting receive loop")
		for {
			in, err := stream.Recv()
			if err == io.EOF {
				close(waitc)
				return
			}
			if err != nil {
				log.Fatalf("client.WatchFrostMessages failed: %v", err)
			}
			log.Printf("Got message", in)
		}
	}()

	event := &pb.WatchFrostMessagesRequest{RequestId: &pb.UUID{}, Command: &pb.WatchFrostMessagesRequest_OpenStream_{
		OpenStream: &pb.WatchFrostMessagesRequest_OpenStream{
			ValidatorIds: []*pb.PolygonEdgeValidator{},
		},
	}}

	if err := stream.Send(event); err != nil {
		log.Fatalf("client.EventStream: stream.Send(%v) failed: %v", event, err)
	}
	stream.CloseSend()
	<-waitc

	fmt.Println("Ending go client")

}
