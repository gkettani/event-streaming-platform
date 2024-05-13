package main

import (
	"log"
	"net/rpc"
	"time"

	"github.com/gkettani/event-streaming-platform/pkg/common"
)

func main() {
	client, err := rpc.DialHTTP("tcp", "localhost:1234")
	if err != nil {
		log.Fatal("Error connecting to server:", err)
	}

	for {
		msg := common.MessageBody{
			Type: "send",
			Key:  "k1",
			Msg:  123,
		}
		sendMessage(client, msg)
		time.Sleep(5 * time.Second) // Produce a message every 5 seconds
	}
}

func sendMessage(client *rpc.Client, msg common.MessageBody) {
	var reply common.MessageBody
	err := client.Call("Server.Send", msg, &reply)
	if err != nil {
		log.Println("Error sending message:", err)
	}
	log.Printf("Message Ack: Offset=%d\n", reply.Offset)
}
