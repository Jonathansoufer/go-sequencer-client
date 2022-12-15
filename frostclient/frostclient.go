package frostclient

import (
	"context"
	"io"
	"log"
	"math/rand"

	pb "github.com/topos-network/go-topos-sequencer-client/frostclient/proto"
	grpc "google.golang.org/grpc"
)

type FrostServiceClient struct {
	Conn          *grpc.ClientConn
	Client        pb.FrostAPIServiceClient
	Closec        chan struct{}                       // Signal to message processing loop to shut down
	Inbox         chan *pb.WatchFrostMessagesResponse // Messages received from the service
	Ctx           context.Context
	contextCancel context.CancelFunc
}

func NewFrostServiceClient(serverAddress string, validatorAccount string) (*FrostServiceClient, error) {
	conn, err := grpc.Dial(serverAddress, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("could not connect to frost service: %v", err)
		return nil, err
	}

	// create stream
	client := pb.NewFrostAPIServiceClient(conn)

	ctx, cancel := context.WithCancel(context.Background())

	// Subscribe to frost service
	stream, err := client.WatchFrostMessages(ctx)
	if err != nil {
		log.Fatalf("error opening frost messages stream: %v", err)
		return nil, err
	}

	closec := make(chan struct{})
	inbox := make(chan *pb.WatchFrostMessagesResponse)
	go func() {
		processFrostMessages(closec, inbox, stream)
	}()

	id := pb.UUID{
		MostSignificantBits:  rand.Uint64(),
		LeastSignificantBits: rand.Uint64(),
	}

	request := &pb.WatchFrostMessagesRequest{RequestId: &id, Command: &pb.WatchFrostMessagesRequest_OpenStream_{
		OpenStream: &pb.WatchFrostMessagesRequest_OpenStream{
			ValidatorIds: []*pb.PolygonEdgeValidator{
				&pb.PolygonEdgeValidator{
					Address: validatorAccount,
				},
			},
		},
	}}

	// Sending watch frost message request to the service
	if err := stream.Send(request); err != nil {
		log.Fatalf("client.EventStream: stream.Send(%v) failed: %v", request, err)
		return nil, err
	}

	frostServiceClient := &FrostServiceClient{
		Conn:          conn,
		Client:        client,
		Closec:        closec,
		Inbox:         inbox,
		contextCancel: cancel,
		Ctx:           ctx,
	}

	return frostServiceClient, nil
}

func receive(closec chan struct{},
	stream pb.FrostAPIService_WatchFrostMessagesClient) *pb.WatchFrostMessagesResponse {
	in, err := stream.Recv()
	if err == io.EOF {
		log.Fatalf("client.WatchFrostMessages EOF")
		closec <- struct{}{}
	}
	if err != nil {
		log.Fatalf("client.WatchFrostMessages failed: %v", err)
	}
	return in
}

func processFrostMessages(closec chan struct{}, inbox chan *pb.WatchFrostMessagesResponse,
	stream pb.FrostAPIService_WatchFrostMessagesClient) {
	log.Println("Starting frost messages processing loop")
	for {
		select {
		case inbox <- receive(closec, stream):
		case <-closec:
			log.Printf("Close signal received")
			return
		}
	}
}
