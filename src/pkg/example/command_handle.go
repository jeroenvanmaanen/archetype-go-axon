package example

import (
    context "context"
    fmt "fmt"
    log "log"
    time "time"
    axonserver "github.com/jeroenvm/archetype-nix-go/src/pkg/grpc/axonserver"
    grpcExample "github.com/jeroenvm/archetype-nix-go/src/pkg/grpc/example"
    grpc "google.golang.org/grpc"
    proto "github.com/golang/protobuf/proto"
    uuid "github.com/google/uuid"
)

func HandleCommands(host string, port int) (conn *grpc.ClientConn) {
    conn, clientInfo := WaitForServer(host, port, "Command Handler")

    log.Printf("Command handler: Connection: %v", conn)
    client := axonserver.NewCommandServiceClient(conn)
    log.Printf("Command handler: Client: %v", client)

    stream, e := client.OpenStream(context.Background())
    log.Printf("Command handler: Stream: %v: %v", stream, e)

    subscribe("GreetCommand", stream, clientInfo)
    subscribe("RecordCommand", stream, clientInfo)
    subscribe("StopCommand", stream, clientInfo)

    go worker(stream, conn, clientInfo.ClientId)

    return conn;
}

func subscribe(commandName string, stream axonserver.CommandService_OpenStreamClient, clientInfo *axonserver.ClientIdentification) {
    id := uuid.New()
    subscription := axonserver.CommandSubscription {
        MessageId: id.String(),
        Command: commandName,
        ClientId: clientInfo.ClientId,
        ComponentName: clientInfo.ComponentName,
    }
    log.Printf("Command handler: Subscription: %v", subscription)
    subscriptionRequest := axonserver.CommandProviderOutbound_Subscribe {
        Subscribe: &subscription,
    }

    outbound := axonserver.CommandProviderOutbound {
        Request: &subscriptionRequest,
    }

    e := stream.Send(&outbound)
    if e != nil {
        panic(fmt.Sprintf("Command handler: Error sending subscription", e))
    }
}

func worker(stream axonserver.CommandService_OpenStreamClient, conn *grpc.ClientConn, clientId string) {
    for true {
        addPermits(1, stream, clientId)

        log.Printf("Command handler: Waiting for command")
        inbound, e := stream.Recv()
        if (e != nil) {
          log.Printf("Command handler: Error on receive, %v", e)
          continue
        }
        log.Printf("Command handler: Inbound: %v", inbound)
        command := inbound.GetCommand()
        if (command != nil) {
            commandName := command.Name
            if (commandName == "GreetCommand") {
                log.Printf("Received GreetCommand")
                handleGreetCommand(command, stream, conn)
            } else if (commandName == "RecordCommand") {
                log.Printf("Received RecordCommand")
                handleRecordCommand(command, stream, conn)
            } else if (commandName == "StopCommand") {
                log.Printf("Received StopCommand")
                handleStopCommand(command, stream, conn)
            } else {
                log.Printf("Received unknown command: %v", commandName)
            }
        }
    }
}

func handleGreetCommand(command *axonserver.Command, stream axonserver.CommandService_OpenStreamClient, conn *grpc.ClientConn) {
    deserializedCommand := grpcExample.GreetCommand{}
    e := proto.Unmarshal(command.Payload.Data, &deserializedCommand)
    if (e != nil) {
        log.Printf("Could not unmarshall GreetCommand")
    }

    event := grpcExample.GreetedEvent {
        Message: deserializedCommand.Message,
    }
    data, err := proto.Marshal(&event)
    if err != nil {
        log.Printf("Server: Error while marshalling event")
        return
    }
    serializedEvent := axonserver.SerializedObject{
        Type: "GreetedEvent",
        Data: data,
    }

    appendEvent(&serializedEvent, deserializedCommand.AggregateIdentifier, conn)
    respond(stream, command.MessageIdentifier)
}

func handleRecordCommand(command *axonserver.Command, stream axonserver.CommandService_OpenStreamClient, conn *grpc.ClientConn) {
    deserializedCommand := grpcExample.RecordCommand{}
    e := proto.Unmarshal(command.Payload.Data, &deserializedCommand)
    if (e != nil) {
        log.Printf("Could not unmarshall RecordCommand")
    }
    event := grpcExample.StartedRecordingEvent {}
    data, err := proto.Marshal(&event)
    if err != nil {
        log.Printf("Server: Error while marshalling event")
        return
    }
    serializedEvent := axonserver.SerializedObject{
        Type: "StartedRecordingEvent",
        Data: data,
    }

    appendEvent(&serializedEvent, deserializedCommand.AggregateIdentifier, conn)
    respond(stream, command.MessageIdentifier)
}

func handleStopCommand(command *axonserver.Command, stream axonserver.CommandService_OpenStreamClient, conn *grpc.ClientConn) {
    deserializedCommand := grpcExample.StopCommand{}
    e := proto.Unmarshal(command.Payload.Data, &deserializedCommand)
    if (e != nil) {
        log.Printf("Could not unmarshall StopCommand")
    }
    event := grpcExample.StoppedRecordingEvent {}
    data, err := proto.Marshal(&event)
    if err != nil {
        log.Printf("Server: Error while marshalling event")
        return
    }
    serializedEvent := axonserver.SerializedObject{
        Type: "StoppedRecordingEvent",
        Data: data,
    }

    appendEvent(&serializedEvent, deserializedCommand.AggregateIdentifier, conn)
    respond(stream, command.MessageIdentifier)
}

func respond(stream axonserver.CommandService_OpenStreamClient, requestId string) {
    id := uuid.New()
    commandResponse := axonserver.CommandResponse {
        MessageIdentifier: id.String(),
        RequestIdentifier: requestId,
    }
    log.Printf("Command handler: Command response: %v", commandResponse)
    commandResponseRequest := axonserver.CommandProviderOutbound_CommandResponse {
        CommandResponse: &commandResponse,
    }

    outbound := axonserver.CommandProviderOutbound {
        Request: &commandResponseRequest,
    }

    e := stream.Send(&outbound)
    if e != nil {
        panic(fmt.Sprintf("Command handler: Error sending command response", e))
    }
}

func addPermits(amount int64, stream axonserver.CommandService_OpenStreamClient, clientId string) {
    flowControl := axonserver.FlowControl {
        ClientId: clientId,
        Permits: amount,
    }
    log.Printf("Command handler: Flow control: %v", flowControl)
    flowControlRequest := axonserver.CommandProviderOutbound_FlowControl {
        FlowControl: &flowControl,
    }

    outbound := axonserver.CommandProviderOutbound {
        Request: &flowControlRequest,
    }

    e := stream.Send(&outbound)
    if e != nil {
        panic(fmt.Sprintf("Command handler: Error sending flow control", e))
    }
}

func appendEvent(message *axonserver.SerializedObject, aggregateId string, conn *grpc.ClientConn) {
    client := axonserver.NewEventStoreClient(conn)

    readRequest := axonserver.ReadHighestSequenceNrRequest {
        AggregateId: aggregateId,
        FromSequenceNr: 0,
    }
    log.Printf("Command handler: Read highest sequence-nr request: %v", readRequest)

    response, e := client.ReadHighestSequenceNr(context.Background(), &readRequest);
    if (e != nil) {
        log.Fatalf("Command handler: Error while reading highest sequence-nr: %v", e)
        return
    }

    log.Printf("Command handler: Response: %v", response)
    next := response.ToSequenceNr + 1;
    log.Printf("Command handler: Next sequence number: %v", next)

    timestamp := time.Now().UnixNano() / 1000000

    id := uuid.New()
    event := axonserver.Event {
        MessageIdentifier: id.String(),
        AggregateIdentifier: aggregateId,
        AggregateSequenceNumber: next,
        AggregateType: "ExampleAggregate",
        Timestamp: timestamp,
        Snapshot: false,
        Payload: message,
    }
    log.Printf("Command handler: Event: %v", event)

    stream, e := client.AppendEvent(context.Background())
    if (e != nil) {
        log.Fatalf("Command handler: Error while preparing to append event: %v", e)
        return
    }

    e = stream.Send(&event)
    if (e != nil) {
        log.Fatalf("Command handler: Error while sending event: %v", e)
        return
    }

    confirmation, e := stream.CloseAndRecv()
    if (e != nil) {
        log.Fatalf("Command handler: Error while sending event: %v", e)
        return
    }

    log.Printf("Command handler: Confirmation: %v", confirmation)
}
