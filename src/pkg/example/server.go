package example

import (
    context "context"
    fmt "fmt"
    net "net"
    log "log"
    axonserver "github.com/jeroenvm/archetype-nix-go/src/pkg/grpc/axonserver"
    grpcExample "github.com/jeroenvm/archetype-nix-go/src/pkg/grpc/example"
    grpc "google.golang.org/grpc"
    reflection "google.golang.org/grpc/reflection"
)

type GreeterServer struct {
    conn *grpc.ClientConn;
    clientInfo *axonserver.ClientIdentification;
}

func (s *GreeterServer) Greet(c context.Context, greeting *grpcExample.Greeting) (*grpcExample.Acknowledgement, error) {
    log.Printf("Server: Received greeting: %v", (*greeting).Message)
    ack := grpcExample.Acknowledgement{
        Message: "Good day to you too!",
    }
    SubmitCommand("beep", s.conn, s.clientInfo)
    return &ack, nil
}

func Serve(conn *grpc.ClientConn, clientInfo *axonserver.ClientIdentification) {
    port := 8181
    lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
    if err != nil {
        log.Fatalf("Server: Failed to listen: %v", err)
        panic("Server: Panic!")
    }
    log.Printf("Server: Listening on port: %d", port)
    grpcServer := grpc.NewServer()
    grpcExample.RegisterGreeterServiceServer(grpcServer, &GreeterServer{conn,clientInfo})
    reflection.Register(grpcServer)
    // ... // determine whether to use TLS
    grpcServer.Serve(lis)
}