package main

import (
	"log"
	"net/http"

	"github.com/go-broadcast/broadcast"
	"github.com/go-broadcast/examples"
	"github.com/go-broadcast/examples/service"
	"github.com/improbable-eng/grpc-web/go/grpcweb"
	"google.golang.org/grpc"
)

func main() {
	log.Println("Starting gRPC example...")

	// GRPC server
	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)

	bcast, err := broadcast.New()
	if err != nil {
		log.Fatal(err)
	}
	service.RegisterChatServiceServer(grpcServer, &examples.ChatService{
		Broadcaster: bcast,
	})

	// Wrapped GRPC server
	wrappedGrpc := grpcweb.WrapServer(
		grpcServer,
		grpcweb.WithAllowedRequestHeaders([]string{"*"}),
		grpcweb.WithOriginFunc(func(o string) bool { return true }),
	)

	// Static files
	http.Handle("/", http.FileServer(http.Dir("../../static/grpc")))

	// HTTP server
	httpServer := &http.Server{
		Addr: ":5200",
	}
	httpServer.Handler = http.HandlerFunc(func(rw http.ResponseWriter, request *http.Request) {
		if wrappedGrpc.IsAcceptableGrpcCorsRequest(request) || wrappedGrpc.IsGrpcWebRequest(request) {
			wrappedGrpc.ServeHTTP(rw, request)
			return
		}

		http.DefaultServeMux.ServeHTTP(rw, request)
	})

	// Start HTTP server
	log.Println("Listening on http://localhost:5200")
	httpServer.ListenAndServe()
}
