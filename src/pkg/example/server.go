package example

import (
    context "context"
    fmt "fmt"
    net "net"
    log "log"
    grpcExample "github.com/jeroenvm/archetype-nix-go/src/pkg/grpc/example"
    grpc "google.golang.org/grpc"
    reflection "google.golang.org/grpc/reflection"
)

type GreeterServer struct {
    conn *grpc.ClientConn;
}

func (s *GreeterServer) Greet(c context.Context, greeting *grpcExample.Greeting) (*grpcExample.Acknowledgement, error) {
    fmt.Println(fmt.Sprintf("Received greeting: %v", (*greeting).Message))
    ack := grpcExample.Acknowledgement{
        Message: "Good day to you too!",
    }
    SubmitCommand("bliep", s.conn)
    return &ack, nil
}

func Serve(conn *grpc.ClientConn) {
    port := 8181
    lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
    if err != nil {
        log.Fatalf("failed to listen: %v", err)
    }
    log.Printf("Listening on port: %d", port)
    grpcServer := grpc.NewServer()
    grpcExample.RegisterGreeterServiceServer(grpcServer, &GreeterServer{conn})
    reflection.Register(grpcServer)
    // ... // determine whether to use TLS
    grpcServer.Serve(lis)
}