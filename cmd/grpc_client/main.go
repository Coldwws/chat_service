package main

import (
	desc "chat_server/pkg/chat_v1"
	"context"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/timestamppb"
)

const (
	address = "localhost:50052"
)

func main() {
	conn, err := grpc.Dial(address,grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil{
		log.Fatalf("Failed to connect to server")
	}
	defer conn.Close()


	c := desc.NewChatV1Client(conn)
	ctx,cancel := context.WithTimeout(context.Background(),time.Second)
	defer cancel()

	s,err := c.SendMessage(ctx, &desc.SendMessageRequest{
		From: "John Doe",
		Text: "Hello World",
		Timestamp: &timestamppb.Timestamp{
			Seconds: time.Now().Unix(),
		},
	})
	if err != nil{
		log.Fatalf("Failed to send message: %v", err)
	}
	log.Printf("Message sent: %v", s)


	cr, err := c.Create(ctx, &desc.CreateRequest{
			Usernames : []string{"John Doe", "Jane Doe"},
	})
	if err != nil{
		log.Printf("Failed to create: %v",err)
	}

	log.Printf("Chat created: %v", cr)
}