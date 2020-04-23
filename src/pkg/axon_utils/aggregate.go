package axon_utils

import (
    context "context"
    fmt "fmt"
    log "log"
    time "time"

    grpc "google.golang.org/grpc"
    uuid "github.com/google/uuid"

    axon_server "github.com/jeroenvm/archetype-go-axon/src/pkg/grpc/axon_server"
)

func SubscribeCommand(commandName string, stream axon_server.CommandService_OpenStreamClient, clientInfo *axon_server.ClientIdentification) {
    id := uuid.New()
    subscription := axon_server.CommandSubscription {
        MessageId: id.String(),
        Command: commandName,
        ClientId: clientInfo.ClientId,
        ComponentName: clientInfo.ComponentName,
    }
    log.Printf("Command handler: Subscription: %v", subscription)
    subscriptionRequest := axon_server.CommandProviderOutbound_Subscribe {
        Subscribe: &subscription,
    }

    outbound := axon_server.CommandProviderOutbound {
        Request: &subscriptionRequest,
    }

    e := stream.Send(&outbound)
    if e != nil {
        panic(fmt.Sprintf("Command handler: Error sending subscription", e))
    }
}

func AppendEvent(message *axon_server.SerializedObject, aggregateId string, conn *grpc.ClientConn) {
    client := axon_server.NewEventStoreClient(conn)

    readRequest := axon_server.ReadHighestSequenceNrRequest {
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
    event := axon_server.Event {
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

func CommandRespond(stream axon_server.CommandService_OpenStreamClient, requestId string) {
    id := uuid.New()
    commandResponse := axon_server.CommandResponse {
        MessageIdentifier: id.String(),
        RequestIdentifier: requestId,
    }
    log.Printf("Command handler: Command response: %v", commandResponse)
    commandResponseRequest := axon_server.CommandProviderOutbound_CommandResponse {
        CommandResponse: &commandResponse,
    }

    outbound := axon_server.CommandProviderOutbound {
        Request: &commandResponseRequest,
    }

    e := stream.Send(&outbound)
    if e != nil {
        panic(fmt.Sprintf("Command handler: Error sending command response", e))
    }
}

func ReportError(stream axon_server.CommandService_OpenStreamClient, requestId string, errorCode string, errorMessageText string) {
    errorMessage := axon_server.ErrorMessage{
        Message: errorMessageText,
        Location: "",
        Details: nil,
    }

    id := uuid.New()
    commandResponse := axon_server.CommandResponse {
        MessageIdentifier: id.String(),
        RequestIdentifier: requestId,
        ErrorCode: errorCode,
        ErrorMessage: &errorMessage,
    }
    log.Printf("Command handler: Command error: %v", commandResponse)
    commandResponseRequest := axon_server.CommandProviderOutbound_CommandResponse {
        CommandResponse: &commandResponse,
    }

    outbound := axon_server.CommandProviderOutbound {
        Request: &commandResponseRequest,
    }

    e := stream.Send(&outbound)
    if e != nil {
        panic(fmt.Sprintf("Command handler: Error sending command error", e))
    }
}

func CommandAddPermits(amount int64, stream axon_server.CommandService_OpenStreamClient, clientId string) {
    flowControl := axon_server.FlowControl {
        ClientId: clientId,
        Permits: amount,
    }
    log.Printf("Command handler: Flow control: %v", flowControl)
    flowControlRequest := axon_server.CommandProviderOutbound_FlowControl {
        FlowControl: &flowControl,
    }

    outbound := axon_server.CommandProviderOutbound {
        Request: &flowControlRequest,
    }

    e := stream.Send(&outbound)
    if e != nil {
        panic(fmt.Sprintf("Command handler: Error sending flow control", e))
    }
}
