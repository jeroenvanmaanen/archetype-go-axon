package axonutils

import (
    context "context"
    fmt "fmt"
    log "log"
    time "time"

    grpc "google.golang.org/grpc"
    uuid "github.com/google/uuid"

    axonserver "github.com/jeroenvm/archetype-go-axon/src/pkg/grpc/axonserver"
)

func SubscribeCommand(commandName string, stream axonserver.CommandService_OpenStreamClient, clientInfo *axonserver.ClientIdentification) {
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

func AppendEvent(message *axonserver.SerializedObject, aggregateId string, conn *grpc.ClientConn) {
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

func CommandRespond(stream axonserver.CommandService_OpenStreamClient, requestId string) {
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

func ReportError(stream axonserver.CommandService_OpenStreamClient, requestId string, errorCode string, errorMessageText string) {
    errorMessage := axonserver.ErrorMessage{
        Message: errorMessageText,
        Location: "",
        Details: nil,
    }

    id := uuid.New()
    commandResponse := axonserver.CommandResponse {
        MessageIdentifier: id.String(),
        RequestIdentifier: requestId,
        ErrorCode: errorCode,
        ErrorMessage: &errorMessage,
    }
    log.Printf("Command handler: Command error: %v", commandResponse)
    commandResponseRequest := axonserver.CommandProviderOutbound_CommandResponse {
        CommandResponse: &commandResponse,
    }

    outbound := axonserver.CommandProviderOutbound {
        Request: &commandResponseRequest,
    }

    e := stream.Send(&outbound)
    if e != nil {
        panic(fmt.Sprintf("Command handler: Error sending command error", e))
    }
}

func CommandAddPermits(amount int64, stream axonserver.CommandService_OpenStreamClient, clientId string) {
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
