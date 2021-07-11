package examples

import (
	"context"
	"log"

	"github.com/go-broadcast/broadcast"
	"github.com/go-broadcast/examples/service"
	"google.golang.org/protobuf/types/known/emptypb"
)

type ChatService struct {
	service.UnimplementedChatServiceServer
	Broadcaster broadcast.Broadcaster
}

func (s *ChatService) SendMessage(ctx context.Context, request *service.SendMessageRequest) (*emptypb.Empty, error) {
	event := &service.ChatMessage{
		Contents: request.Message,
		From:     request.User,
	}

	s.Broadcaster.ToRoom(event, "chat-room", "user"+request.User)

	return &emptypb.Empty{}, nil
}

func (s *ChatService) Subscribe(request *service.SubscribeRequest, server service.ChatService_SubscribeServer) error {
	sub := s.Broadcaster.Subscribe(func(data interface{}) {
		message, ok := data.(*service.ChatMessage)

		if !ok {
			log.Println("Invalid message")
			return
		}

		server.Send(message)
	})

	s.Broadcaster.JoinRoom(sub, "chat-room", "user"+request.User)
	<-server.Context().Done()
	s.Broadcaster.Unsubscribe(sub)

	return nil
}
