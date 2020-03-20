package example

import (
    context "context"
    errors "errors"
    fmt "fmt"
    net "net"
    log "log"
    axonserver "github.com/jeroenvm/archetype-nix-go/src/pkg/grpc/axonserver"
    grpcExample "github.com/jeroenvm/archetype-nix-go/src/pkg/grpc/example"
    grpc "google.golang.org/grpc"
    proto "github.com/golang/protobuf/proto"
    reflection "google.golang.org/grpc/reflection"
)

type GreeterServer struct {
    conn *grpc.ClientConn;
    clientInfo *axonserver.ClientIdentification;
}

func (s *GreeterServer) Greet(c context.Context, greeting *grpcExample.Greeting) (*grpcExample.Acknowledgement, error) {
    message := (*greeting).Message
    log.Printf("Server: Received greeting: %v", message)
    ack := grpcExample.Acknowledgement{
        Message: "Good day to you too!",
    }
    command := grpcExample.GreetCommand {
        AggregateIdentifier: "single_aggregate",
        Message: greeting,
    }
    data, err := proto.Marshal(&command)
    if err != nil {
        log.Printf("Server: Error while marshalling command")
        return nil, errors.New("Marshalling error")
    }
    serializedCommand := axonserver.SerializedObject{
        Type: "GreetCommand",
        Data: data,
    }
    SubmitCommand(&serializedCommand, s.conn, s.clientInfo)
    return &ack, nil
}

func (s *GreeterServer) Greetings(empty *grpcExample.Empty, greetingsServer grpcExample.GreeterService_GreetingsServer) error {
    greeting := grpcExample.Greeting {
        Message: "Hello, World!",
    }
    log.Printf("Greetings streamed reply: %v", greeting)
    greetingsServer.Send(&greeting)
    log.Printf("Greetings streamed reply sent")
    return nil
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