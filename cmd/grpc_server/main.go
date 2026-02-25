package main

import (
	"chat_server/internal/config"
	desc "chat_server/pkg/chat_v1"
	"context"
	"log"
	"net"
	"os"
	"google.golang.org/protobuf/types/known/emptypb"
	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
)

type server struct{
	db *pgxpool.Pool
	desc.UnimplementedChatV1Server
}

func (s *server) Create(ctx context.Context, req *desc.CreateRequest)(*desc.CreateResponse,error){
	usernames := req.GetUsernames()

	if len(usernames) == 0{
		return nil, status.Error(codes.InvalidArgument,"usernames is required")
	}

	tx,err := s.db.Begin(ctx)
	if err != nil{
		return nil, status.Error(codes.Internal,"failed to begin transaction")
	}
	defer func ()  {
		_ = tx.Rollback(ctx)
	}()

	qbChat := sq.Insert("chats").PlaceholderFormat(sq.Dollar).
	Columns("created_at").
	Values(sq.Expr("now()")).
	Suffix("RETURNING id")
	
	q1,a1,err := qbChat.ToSql()
	if err != nil{
		return nil, status.Errorf(codes.Internal, "failed to build query: %v",err)
	}
	var chatID int64
	if err := tx.QueryRow(ctx,q1,a1...).Scan(&chatID);err != nil{
		return nil, status.Error(codes.Internal,"failed to create chat")
	}

	for _, u := range usernames{
		if u == ""{
			return nil, status.Error(codes.InvalidArgument,"Username cannot be empty")
		}
		qbU := sq.Insert("chat_users").
		PlaceholderFormat(sq.Dollar).
		Columns("chat_id","username").Values(chatID,u)

		q2,a2,err := qbU.ToSql()
		if err != nil{
			return nil, status.Error(codes.Internal,"Failed to build query")
		}

		if _, err := tx.Exec(ctx,q2,a2...); err != nil{
			return nil, status.Error(codes.Internal,"failed to add user to chat")
		}
	}
	if err:= tx.Commit(ctx);err != nil{
		return nil, status.Error(codes.Internal,"failed to commit tx")
	}

	return &desc.CreateResponse{Id: chatID},nil
}

func (s *server) Delete(ctx context.Context, req *desc.DeleteRequest)(*emptypb.Empty,error){
	deleteID := req.GetId()
	log.Printf("Delete chat with ID: ",deleteID)

	if deleteID == 0{
		return nil, status.Error(codes.InvalidArgument,"ID is required")
	}
	qb := sq.Delete("chats").PlaceholderFormat(sq.Dollar).Where(sq.Eq{"id":deleteID})

	query,args,err := qb.ToSql()
	if err != nil{
		return nil,status.Error(codes.Internal,"Failed to build query")
	}

	ct, err := s.db.Exec(ctx,query,args...)
	if err != nil{
		return nil,status.Error(codes.Internal,"Failed to delete chat")
	} 
	if ct.RowsAffected() == 0{
		return nil,status.Error(codes.NotFound,"Chat not found")
	}
	return &emptypb.Empty{},nil
}


func (s *server)SendMessage(ctx context.Context, req *desc.SendMessageRequest)(*emptypb.Empty,error){
	if req.GetChatId() == 0{
		return nil, status.Error(codes.InvalidArgument,"chat_id is required")
	}
	if req.GetFrom() == ""{
		return nil, status.Error(codes.InvalidArgument,"from is required")
	}

	if req.GetText() == ""{
		return nil, status.Error(codes.InvalidArgument,"text is required")
	}

	if req.GetTimestamp() == nil{
		return nil, status.Error(codes.InvalidArgument,"timestamp is required")
	}

	qb := sq.Insert("messages").
	PlaceholderFormat(sq.Dollar).Columns("chat_id","sender","text","created_at").
	Values(req.GetChatId(),req.GetFrom(),req.GetText(),req.GetTimestamp().AsTime())

	query,args,err := qb.ToSql()
	if err != nil{
		return nil, status.Error(codes.Internal,"Failed to build query")
	}
	_, err = s.db.Exec(ctx,query,args...)
	if err != nil{
		return nil,status.Error(codes.Internal,"Failed to insert message")
	}
	return &emptypb.Empty{},nil
}

func main(){
	if f := os.Getenv("ENV_FILE"); f!= ""{
		_ = godotenv.Load(f)
	}

	cfg := config.LoadConfig()
	ctx := context.Background()

	poll,err := pgxpool.Connect(ctx,cfg.PG.DSN())
	if err != nil{
		log.Fatalf("Failed to connect database: %v",err)
	}
	defer poll.Close()

	lis,err := net.Listen("tcp",cfg.GRPC.Addr())
	if err != nil{
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	reflection.Register(s)
	srv := &server{db : poll}
	desc.RegisterChatV1Server(s, srv)

	log.Printf("Server listening at addr: %v", lis.Addr())

	if err := s.Serve(lis); err != nil{
		log.Fatalf("failed to serve: %v", err)
	}

}