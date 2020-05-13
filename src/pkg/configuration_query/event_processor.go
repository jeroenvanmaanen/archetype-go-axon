package configuration_query

import (
    context "context"
    log "log"

    proto "github.com/golang/protobuf/proto"
    uuid "github.com/google/uuid"

    authentication "github.com/jeroenvm/archetype-go-axon/src/pkg/authentication"
    axon_server "github.com/jeroenvm/archetype-go-axon/src/pkg/grpc/axon_server"
    axon_utils "github.com/jeroenvm/archetype-go-axon/src/pkg/axon_utils"
    grpc_example "github.com/jeroenvm/archetype-go-axon/src/pkg/grpc/example"
    trusted "github.com/jeroenvm/archetype-go-axon/src/pkg/trusted"
)

// Redeclare event types, so that they can be extended with event handler methods.
type PropertyChangedEvent    struct { grpc_example.PropertyChangedEvent    }
type TrustedKeyAddedEvent    struct { grpc_example.TrustedKeyAddedEvent    }
type TrustedKeyRemovedEvent  struct { grpc_example.TrustedKeyRemovedEvent  }
type KeyManagerAddedEvent    struct { grpc_example.KeyManagerAddedEvent    }
type KeyManagerRemovedEvent  struct { grpc_example.KeyManagerRemovedEvent  }
type CredentialsAddedEvent   struct { grpc_example.CredentialsAddedEvent   }
type CredentialsRemovedEvent struct { grpc_example.CredentialsRemovedEvent }

// Projection

type Projection struct {
    Configuration map[string]string
}

var projection Projection

func ProcessEvents(host string, port int) *axon_utils.ClientConnection {
    projection = Projection{
        Configuration: make(map[string]string),
    }

    clientConnection, stream := axon_utils.WaitForServer(host, port, "Configuration event processor")
    log.Printf("Connection and client info: %v: %v", clientConnection, stream)

    processorName := "configuration-event-processor"

    if e := registerProcessor(processorName, stream, clientConnection.ClientInfo); e != nil {
        return clientConnection
    }

    go eventProcessorWorker(stream, clientConnection, processorName)

    return clientConnection;
}

func registerProcessor(processorName string, stream *axon_server.PlatformService_OpenStreamClient, clientInfo *axon_server.ClientIdentification) error {
    processorInfo := axon_server.EventProcessorInfo {
        ProcessorName: processorName,
        Mode: "Tracking",
        ActiveThreads: 0,
        Running: true,
        Error: false,
        SegmentStatus: nil,
        AvailableThreads: 1,
    }
    log.Printf("Configuration event processor: event processor info: %v", processorInfo)
    subscriptionRequest := axon_server.PlatformInboundInstruction_EventProcessorInfo {
        EventProcessorInfo: &processorInfo,
    }

    id := uuid.New()
    inbound := axon_server.PlatformInboundInstruction {
        Request: &subscriptionRequest,
        InstructionId: id.String(),
    }
    log.Printf("Configuration event processor: event processor info: instruction ID: %v", inbound.InstructionId)

    e := (*stream).Send(&inbound)
    if e != nil {
        log.Printf("Configuration event processor: Error sending registration", e)
        return e
    }

    e = eventProcessorReceivePlatformInstruction(stream)
    if e != nil {
        log.Printf("Configuration event processor: Error while waiting for acknowledgement of registration")
        return e
    }
    return nil
}

// Map an event name as stored in AxonServer to an object that has two aspects:
// 1. It is a proto.Message, so it can be unmarshalled from the Axon event
// 2. It is an axon_utils.Event, so it can be applied to a projection
func prepareUnmarshal(payloadType string) (event axon_utils.Event) {
    log.Printf("Configuration event processor: Payload type: %v", payloadType)
    switch payloadType {
        case "PropertyChangedEvent":    event = &PropertyChangedEvent{}
        case "TrustedKeyAddedEvent":    event = &TrustedKeyAddedEvent{}
        case "TrustedKeyRemovedEvent":  event = &TrustedKeyRemovedEvent{}
        case "KeyManagerAddedEvent":    event = &KeyManagerAddedEvent{}
        case "KeyManagerRemovedEvent":  event = &KeyManagerRemovedEvent{}
        case "CredentialsAddedEvent":   event = &CredentialsAddedEvent{}
        case "CredentialsRemovedEvent": event = &CredentialsRemovedEvent{}
        default: event = nil
    }
    return event
}

func eventProcessorWorker(stream *axon_server.PlatformService_OpenStreamClient, clientConnection *axon_utils.ClientConnection, processorName string) {
    conn := clientConnection.Connection
    clientInfo := clientConnection.ClientInfo
    var token *int64

    eventStoreClient := axon_server.NewEventStoreClient(conn)
    log.Printf("Configuration event processor worker: Event store client: %v", eventStoreClient)
    client, e := eventStoreClient.ListEvents(context.Background())
    if e != nil {
        log.Printf("Configuration event processor worker: Error while opening ListEvents stream: %v", e)
        return
    }
    log.Printf("Configuration event processor worker: List events client: %v", client)

    getEventsRequest := axon_server.GetEventsRequest{
        NumberOfPermits: 1,
        ClientId: clientInfo.ClientId,
        ComponentName: clientInfo.ComponentName,
        Processor: processorName,
    }
    if token != nil {
        getEventsRequest.TrackingToken = *token + 1
    }
    log.Printf("Configuration event processor worker: Get events request: %v", getEventsRequest)


    log.Printf("Configuration event processor worker: Ready to process events")
    defer func() {
        log.Printf("Configuration event processor worker stopped")
    }()
    for true {
        e = client.Send(&getEventsRequest)
        if e != nil {
            log.Printf("Configuration event processor worker: Error while sending get events request: %v", e)
            return
        }

        eventMessage, e := client.Recv()
        if e != nil {
            log.Printf("Configuration event processor worker: Error while receiving next event: %v", e)
            return
        }
        log.Printf("Configuration event processor worker: Next event message: %v", eventMessage)
        getEventsRequest.TrackingToken = eventMessage.Token

        if eventMessage.Event == nil || eventMessage.Event.Payload == nil {
            continue
        }

        payloadType := eventMessage.Event.Payload.Type
        event := prepareUnmarshal(payloadType)
        if event == nil {
            log.Printf("Configuration event processor worker: Skipped unknown event: %v", payloadType)
            continue
        }
        if e = proto.Unmarshal(eventMessage.Event.Payload.Data, event.(proto.Message)); e != nil {
            log.Printf("Configuration event processor worker: Unmarshalling of %v failed: %v", payloadType, e)
            return
        }
        event.ApplyTo(&projection)
    }
}

// Event Handlers

func (event *PropertyChangedEvent) ApplyTo(projectionWrapper interface{}) {
    projection := projectionWrapper.(*Projection)
    key := event.Property.Key
    value := event.Property.Value
    log.Printf("Configuration event handler: Set property: %v: %v", key, value)
    projection.Configuration[key] = value
}

func (event *TrustedKeyAddedEvent) ApplyTo(projectionWrapper interface{}) {
    trusted.UnsafeSetTrustedKey(event.PublicKey)
}

func (event *TrustedKeyRemovedEvent) ApplyTo(projectionWrapper interface{}) {
    trusted.UnsafeSetTrustedKey(getEmptyPublicKey(event.Name))
}

func (event *KeyManagerAddedEvent) ApplyTo(projectionWrapper interface{}) {
    trusted.UnsafeSetKeyManager(event.PublicKey)
}

func (event *KeyManagerRemovedEvent) ApplyTo(projectionWrapper interface{}) {
    trusted.UnsafeSetKeyManager(getEmptyPublicKey(event.Name))
}

func (event *CredentialsAddedEvent) ApplyTo(projectionWrapper interface{}) {
    authentication.UnsafeSetCredentials(event.Credentials)
}

func (event *CredentialsRemovedEvent) ApplyTo(projectionWrapper interface{}) {
    emptyCredentials := grpc_example.Credentials{
        Identifier: event.Identifier,
        Secret: "",
    }
    authentication.UnsafeSetCredentials(&emptyCredentials)
}

// Public accessor

func GetProperty(key string) string {
    return projection.Configuration[key]
}

// Helper functions

func getEmptyPublicKey(name string) *grpc_example.PublicKey {
    return &grpc_example.PublicKey{
        Name: name,
        PublicKey: "",
    }
}

func eventProcessorReceivePlatformInstruction(stream *axon_server.PlatformService_OpenStreamClient) error {
    log.Printf("Configuration event processor receive platform instruction: Waiting for outbound")
    outbound, e := (*stream).Recv()
    if (e != nil) {
      log.Printf("Configuration event processor receive platform instruction: Error on receive, %v", e)
      return e
    }
    log.Printf("Configuration event processor receive platform instruction: Outbound: %v", outbound)
    return nil
}
