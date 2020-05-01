package configuration

import (
    context "context"
    log "log"

    grpc "google.golang.org/grpc"
    proto "github.com/golang/protobuf/proto"
    uuid "github.com/google/uuid"

    authentication "github.com/jeroenvm/archetype-go-axon/src/pkg/authentication"
    axon_server "github.com/jeroenvm/archetype-go-axon/src/pkg/grpc/axon_server"
    axon_utils "github.com/jeroenvm/archetype-go-axon/src/pkg/axon_utils"
    grpc_example "github.com/jeroenvm/archetype-go-axon/src/pkg/grpc/example"
    trusted "github.com/jeroenvm/archetype-go-axon/src/pkg/trusted"
)

var projection Projection

func ProcessEvents(host string, port int) *grpc.ClientConn {
    projection = Projection{
        Configuration: make(map[string]string),
    }

    clientConnection, stream := axon_utils.WaitForServer(host, port, "Configuration event processor")
    log.Printf("Connection and client info: %v: %v", clientConnection, stream)
    conn := clientConnection.Connection

    processorName := "configuration-event-processor"

    if e := registerProcessor(processorName, stream, clientConnection.ClientInfo); e != nil {
        return conn
    }

    go eventProcessorWorker(stream, clientConnection, processorName)

    return conn;
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
        if payloadType == "TrustedKeyAddedEvent" {
            event := grpc_example.TrustedKeyAddedEvent{}
            if e = proto.Unmarshal(eventMessage.Event.Payload.Data, &event); e != nil {
                log.Printf("Configuration event processor worker: Unmarshalling of TrustedKeyAddedEvent failed: %v", e)
                return
            }
            log.Printf("Configuration event processor worker: Payload of TrustedKeyAddedEvent event: %v", event)
            trusted.UnsafeSetTrustedKey(event.PublicKey)
        } else if payloadType == "TrustedKeyRemovedEvent" {
            event := grpc_example.TrustedKeyRemovedEvent{}
            if e = proto.Unmarshal(eventMessage.Event.Payload.Data, &event); e != nil {
                log.Printf("Configuration event processor worker: Unmarshalling of TrustedKeyRemovedEvent failed: %v", e)
                return
            }
            log.Printf("Configuration event processor worker: Payload of TrustedKeyRemovedEvent event: %v", event)
            trusted.UnsafeSetTrustedKey(getEmptyPublicKey(event.Name))
        } else if payloadType == "KeyManagerAddedEvent" {
            event := grpc_example.KeyManagerAddedEvent{}
            if e = proto.Unmarshal(eventMessage.Event.Payload.Data, &event); e != nil {
                log.Printf("Configuration event processor worker: Unmarshalling of KeyManagerAddedEvent failed: %v", e)
                return
            }
            log.Printf("Configuration event processor worker: Payload of KeyManagerAddedEvent event: %v", event)
            trusted.UnsafeSetKeyManager(event.PublicKey)
        } else if payloadType == "KeyManagerRemovedEvent" {
            event := grpc_example.KeyManagerRemovedEvent{}
            if e = proto.Unmarshal(eventMessage.Event.Payload.Data, &event); e != nil {
                log.Printf("Configuration event processor worker: Unmarshalling of KeyManagerRemovedEvent failed: %v", e)
                return
            }
            log.Printf("Configuration event processor worker: Payload of KeyManagerRemovedEvent event: %v", event)
            trusted.UnsafeSetKeyManager(getEmptyPublicKey(event.Name))
        } else if payloadType == "CredentialsAddedEvent" {
            event := grpc_example.CredentialsAddedEvent{}
            if e = proto.Unmarshal(eventMessage.Event.Payload.Data, &event); e != nil {
                log.Printf("Configuration event processor worker: Unmarshalling of CredentialsAddedEvent failed: %v", e)
                return
            }
            log.Printf("Configuration event processor worker: Payload of CredentialsAddedEvent event: %v", event)
            authentication.UnsafeSetCredentials(event.Credentials)
        } else if payloadType == "CredentialsRemovedEvent" {
            event := grpc_example.CredentialsRemovedEvent{}
            if e = proto.Unmarshal(eventMessage.Event.Payload.Data, &event); e != nil {
                log.Printf("Configuration event processor worker: Unmarshalling of CredentialsRemovedEvent failed: %v", e)
                return
            }
            log.Printf("Configuration event processor worker: Payload of CredentialsRemovedEvent event: %v", event)
            emptyCredentials := grpc_example.Credentials{
                Identifier: event.Identifier,
                Secret: "",
            }
            authentication.UnsafeSetCredentials(&emptyCredentials)
        } else if payloadType == "PropertyChangedEvent" {
            event := grpc_example.PropertyChangedEvent{}
            if e = proto.Unmarshal(eventMessage.Event.Payload.Data, &event); e != nil {
                log.Printf("Configuration event processor worker: Unmarshalling of PropertyChangedEvent failed: %v", e)
                return
            }
            log.Printf("Configuration event processor worker: Payload of PropertyChangedEvent event: %v", event)
            projection.Configuration[event.Property.Key] = event.Property.Value
        } else {
            log.Printf("Configuration event processor worker: no processing necessary for payload type: %v", payloadType)
        }
    }
}

func GetProperty(key string) string {
    return projection.Configuration[key]
}

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
