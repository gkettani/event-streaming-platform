package main

import (
	"log"
	"net/rpc"

	"github.com/gkettani/event-streaming-platform/pkg/common"
)

func main() {
	client, err := rpc.DialHTTP("tcp", "localhost:1234")
	if err != nil {
		log.Fatal("Error connecting to server:", err)
	}

	offset := 1 // Start polling from offset 1
	for {
		msg, err := retrieveMessage(client, offset)
		if err != nil {
			continue
		}
		log.Printf("Polled message: Offset=%d, Type=%s, Key=%s, Msg=%d\n", msg.Offset, msg.Type, msg.Key, msg.Msg)
		offset++ // Increment offset for next polling
	}
}

func retrieveMessage(client *rpc.Client, offset int) (common.MessageBody, error) {
	var msg common.MessageBody
	err := client.Call("Server.Poll", offset, &msg)
	if err != nil {
		return common.MessageBody{}, err
	}
	return msg, nil
}
