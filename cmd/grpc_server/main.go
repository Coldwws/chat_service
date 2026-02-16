package main

import (
	"chat_server/internal/config"
	desc "chat_server/pkg/chat_v1"
	"context"
	"fmt"
	"log"
	"net"
	"sync"
	"time"

	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/types/known/emptypb"
)

const grpcPort = 50052

type server struct{
	desc.UnimplementedChatV1Server
}

type Chat struct{
	ID int64
	Usernames []string
}


type Message struct{
	From string
	Text string
	Timestamp time.Time
}

var(
	mu sync.RWMutex
	chats = make(map[int64]Chat)
	messages = make(map[int64][]Message)
	nextID int64 = 1
)

func genChatID()int64{
	mu.Lock()
	defer mu.Unlock()
	id := nextID
	nextID++
	return id
}


func (s *server)SendMessage(ctx context.Context, req *desc.SendMessageRequest)(*emptypb.Empty,error){

	if req.ChatId == 0{
		return nil, fmt.Errorf("chat_id is required")
	}

	if req.From == ""{
		return nil, fmt.Errorf("from is required")
	}
	if req.Text == ""{
		return nil, fmt.Errorf("text is required")
	}

	msg := Message{
		From : req.From,
		Text : req.Text,
		Timestamp : req.Timestamp.AsTime(),
	}

	mu.Lock()
	if _,ok := messages[req.ChatId]; !ok{
		mu.Unlock()
		return nil, fmt.Errorf("chat not found: %d", req.ChatId)
	}

	messages[req.ChatId] = append(messages[req.ChatId],msg)
	mu.Unlock()

		log.Printf("Send Message | chat : %d | from: %s | text: %s | timestamp: %s",
		req.ChatId, req.From, req.Text, req.Timestamp.AsTime())

	return &emptypb.Empty{},nil
}


func (s *server)Create(ctx context.Context, req *desc.CreateRequest)(*desc.CreateResponse, error){

	if req.Usernames == nil{
		return nil, fmt.Errorf("usernames is required")
	}

	chatId := genChatID()
	
	mu.Lock()
	chats[chatId] = Chat{
		ID: chatId,
		Usernames: req.Usernames,
	}
	messages[chatId] = []Message{}
	mu.Unlock()

	log.Printf("Create Chat | id = %d | usernames = %v",chatId, req.Usernames)
	
	return &desc.CreateResponse{
		Id: chatId,}, nil
}

func main(){
	_ =godotenv.Load("local.env")

	cfg := config.LoadConfig()

	lis,err := net.Listen("tcp",cfg.GRPC.Addr())
	if err != nil{
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	reflection.Register(s)
	desc.RegisterChatV1Server(s, &server{})

	log.Printf("Server listening at addr: %v", lis.Addr())

	if err := s.Serve(lis); err != nil{
		log.Fatalf("failed to serve: %v", err)
	}

}