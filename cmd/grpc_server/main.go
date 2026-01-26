package main

import (
	desc "chat_server/pkg/chat_v1"
	"context"
	"fmt"
	"log"
	"net"
	"google.golang.org/protobuf/types/known/emptypb"
	"github.com/brianvoe/gofakeit"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const grpcPort = 50052

type server struct{
	desc.UnimplementedChatV1Server
}

func (s *server)SendMessage(ctx context.Context, req *desc.SendMessageRequest)(*emptypb.Empty,error){
	if req.From == ""{
		return nil, fmt.Errorf("from is required")
	}
	if req.Text == ""{
		return nil, fmt.Errorf("text is required")

	}

		log.Printf("Send Message | from : %s | text: %s | timestamp: %s",
		req.From, req.Text, req.Timestamp.AsTime())

	return &emptypb.Empty{},nil
}


func (s *server)Create(ctx context.Context, req *desc.CreateRequest)(*desc.CreateResponse, error){
	if req.Usernames == nil{
		return nil, fmt.Errorf("usernames is required")
	}

	chatId := gofakeit.Int64()
	 
	return &desc.CreateResponse{
		Id: chatId,
	}, nil

}

func main(){
	lis,err := net.Listen("tcp",fmt.Sprintf(":%d",grpcPort))
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