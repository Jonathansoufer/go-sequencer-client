package main

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"github.com/topos-network/go-topos-sequencer-client/frostclient"
	pb "github.com/topos-network/go-topos-sequencer-client/frostclient/proto"
)

func readString(reader *bufio.Reader) string {
	line, err := reader.ReadString('\n')
	if err != nil {
		log.Fatalf("error reading terminal input: %v", err)
	}
	fmt.Println("Line readed:", line)
	return line
}

func main() {
	serverAddress := os.Args[1]
	validatorAccount := os.Args[2]

	frostServiceClient, err := frostclient.NewFrostServiceClient(serverAddress, validatorAccount)
	if err != nil {
		log.Fatalf("could not instantiate frost client: %v", err)
		return
	}

	// Wait for message from service
	fmt.Println("Starting main loop")

	reader := bufio.NewReader(os.Stdin)
	line := make(chan string)

	go func() {
		var index = 0
		for {
			message := <-line
			fmt.Println("Sending message to server: ", message)
			index = index + 1
			request := &pb.SubmitFrostMessageRequest{
				FrostMessage: &pb.FrostMessage{
					MessageId: string(index),
					From:      validatorAccount,
					Data: &pb.FrostMessageData{
						Data: &pb.FrostMessageData_Value{
							Value: message,
						},
					},
					Signature: "",
				},
			}
			_, err := frostServiceClient.Client.SubmitFrostMessage(frostServiceClient.Ctx, request)
			if err != nil {
				fmt.Println("Error submiting frost message to service: %v", err)
			}

		}
	}()

	for {
		select {
		case message := <-frostServiceClient.Inbox:
			fmt.Println("Received message:", message)
		case line <- readString(reader):

		}
	}

}
