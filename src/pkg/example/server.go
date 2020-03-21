package example

import (
    context "context"
    errors "errors"
    fmt "fmt"
    io "io"
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
    err = SubmitCommand(&serializedCommand, s.conn, s.clientInfo)
    if err != nil {
        return nil, err
    }
    return &ack, nil
}

func (s *GreeterServer) Greetings(empty *grpcExample.Empty, greetingsServer grpcExample.GreeterService_GreetingsServer) error {
    greeting := grpcExample.Greeting {
        Message: "Hello, World!",
    }
    log.Printf("Server: Greetings streamed reply: %v", greeting)
    greetingsServer.Send(&greeting)
    log.Printf("Server: Greetings streamed reply sent")

    eventStoreClient := axonserver.NewEventStoreClient(s.conn)
    requestEvents := axonserver.GetAggregateEventsRequest {
        AggregateId: "single_aggregate",
        InitialSequence: 0,
        AllowSnapshots: false,
    }
    log.Printf("Server: Request events: %v", requestEvents)
    client, e := eventStoreClient.ListAggregateEvents(context.Background(), &requestEvents)
    if e != nil {
        log.Printf("Server: Error while requesting aggregate events: %v", e)
        return nil
    }
    for {
        eventMessage, e := client.Recv()
        if e == io.EOF {
            log.Printf("Server: End of stream")
            break
        } else if e != nil {
            log.Printf("Server: Error while receiving next event: %v", e)
            break
        }
        log.Printf("Server: Received event: %v", eventMessage)
        if eventMessage.Payload != nil {
            log.Printf("Server: Payload type: %v", eventMessage.Payload.Type)
            if (eventMessage.Payload.Type != "GreetedEvent") {
                continue
            }
            log.Printf("Server: Payload: %v", eventMessage.Payload)
            event := grpcExample.GreetedEvent{}
            e = proto.Unmarshal(eventMessage.Payload.Data, &event)
            if e != nil {
                log.Printf("Server: Error while unmarshalling GreetedEvent")
                continue
            }
            log.Printf("Server: GreetedEvent: %v", event)
            log.Printf("Server: Greetings streamed reply: %v", event.Message)
            greetingsServer.Send(event.Message)
            log.Printf("Server: Greetings streamed reply sent")
        }
    }

    return nil
}

var empty = grpcExample.Empty{}

func (s *GreeterServer) Record(c context.Context, greeting *grpcExample.Empty) (*grpcExample.Empty, error) {
    command := grpcExample.RecordCommand {
        AggregateIdentifier: "single_aggregate",
    }
    data, err := proto.Marshal(&command)
    if err != nil {
        log.Printf("Server: Error while marshalling command")
        return nil, errors.New("Marshalling error")
    }
    serializedCommand := axonserver.SerializedObject{
        Type: "RecordCommand",
        Data: data,
    }
    err = SubmitCommand(&serializedCommand, s.conn, s.clientInfo)
    if err != nil {
        return nil, err
    }
    return &empty, nil
}

func (s *GreeterServer) Stop(c context.Context, greeting *grpcExample.Empty) (*grpcExample.Empty, error) {
    command := grpcExample.StopCommand {
        AggregateIdentifier: "single_aggregate",
    }
    data, err := proto.Marshal(&command)
    if err != nil {
        log.Printf("Server: Error while marshalling command")
        return nil, errors.New("Marshalling error")
    }
    serializedCommand := axonserver.SerializedObject{
        Type: "StopCommand",
        Data: data,
    }
    err = SubmitCommand(&serializedCommand, s.conn, s.clientInfo)
    if err != nil {
        return nil, err
    }
    return &empty, nil
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