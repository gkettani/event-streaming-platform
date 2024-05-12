package main

import (
	"log"
	"net/rpc"
	"net/rpc/jsonrpc"

	"github.com/gkettani/event-streaming-platform/pkg/common"
)

func main() {
	client, err := jsonrpc.Dial("tcp", "localhost:1234")
	if err != nil {
		log.Fatal("Error connecting to server:", err)
	}

	offset := 1 // Start polling from offset 1
	for {
		msg, err := retrieveMessage(client, offset)
		if err != nil {
			log.Println("Error retrieving message:", err)
			continue
		}
		log.Printf("Polled message: Offset=%d, Type=%s, Key=%s, Msg=%d\n", msg.Offset, msg.Type, msg.Key, msg.Msg)
		offset++ // Increment offset for next polling
	}
}

func retrieveMessage(client *rpc.Client, offset int) (common.Message, error) {
	var msg common.Message
	err := client.Call("Server.Poll", offset, &msg)
	if err != nil {
		return common.Message{}, err
	}
	return msg, nil
}
