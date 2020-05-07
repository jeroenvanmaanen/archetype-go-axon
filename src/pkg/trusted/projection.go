package trusted

import (
    context "context"
    io "io"
    log "log"

    proto "github.com/golang/protobuf/proto"

    axon_server "github.com/jeroenvm/archetype-go-axon/src/pkg/grpc/axon_server"
    axon_utils "github.com/jeroenvm/archetype-go-axon/src/pkg/axon_utils"
    grpc_example "github.com/jeroenvm/archetype-go-axon/src/pkg/grpc/example"
)

type Projection struct {
    TrustedKeys map[string]string
    KeyManagers map[string]string
}

func RestoreProjection(aggregateIdentifier string, clientConnection *axon_utils.ClientConnection) *Projection {
    projection := Projection{
        TrustedKeys: make(map[string]string),
        KeyManagers: make(map[string]string),
    }
    log.Printf("Trusted Keys Projection: %v", projection)

    eventStoreClient := axon_server.NewEventStoreClient(clientConnection.Connection)
    requestEvents := axon_server.GetAggregateEventsRequest {
        AggregateId: aggregateIdentifier,
        InitialSequence: 0,
        AllowSnapshots: false,
    }
    log.Printf("Trusted Keys Projection: Request events: %v", requestEvents)
    client, e := eventStoreClient.ListAggregateEvents(context.Background(), &requestEvents)
    if e != nil {
        log.Printf("Trusted Keys Projection: Error while requesting aggregate events: %v", e)
        return nil
    }
    for {
        eventMessage, e := client.Recv()
        if e == io.EOF {
            log.Printf("Trusted Keys Projection: End of stream")
            break
        } else if e != nil {
            log.Printf("Trusted Keys Projection: Error while receiving next event: %v", e)
            break
        }
        log.Printf("Trusted Keys Projection: Received event: %v", eventMessage)
        if eventMessage.Payload != nil {
            log.Printf("Trusted Keys Projection: Payload type: %v", eventMessage.Payload.Type)
            payloadType := eventMessage.Payload.Type
            if (payloadType == "TrustedKeyAddedEvent") {
                event := grpc_example.TrustedKeyAddedEvent{}
                e := proto.Unmarshal(eventMessage.Payload.Data, &event)
                if (e != nil) {
                    log.Printf("Could not unmarshal TrustedKeyAddedEvent")
                }
                projection.TrustedKeys[event.PublicKey.Name] = event.PublicKey.PublicKey
            } else if (payloadType == "TrustedKeyRemovedEvent") {
                event := grpc_example.TrustedKeyRemovedEvent{}
                e := proto.Unmarshal(eventMessage.Payload.Data, &event)
                if (e != nil) {
                    log.Printf("Could not unmarshal TrustedKeyRemovedEvent")
                }
                projection.TrustedKeys[event.Name] = ""
            } else if (payloadType == "KeyManagerAddedEvent") {
                event := grpc_example.KeyManagerAddedEvent{}
                e := proto.Unmarshal(eventMessage.Payload.Data, &event)
                if (e != nil) {
                    log.Printf("Could not unmarshal KeyManagerAddedEvent")
                }
                projection.KeyManagers[event.PublicKey.Name] = event.PublicKey.PublicKey
            } else if (payloadType == "KeyManagerRemovedEvent") {
                event := grpc_example.KeyManagerRemovedEvent{}
                e := proto.Unmarshal(eventMessage.Payload.Data, &event)
                if (e != nil) {
                    log.Printf("Could not unmarshal KeyManagerRemovedEvent")
                }
                projection.KeyManagers[event.Name] = ""
            }
        }
    }

    return &projection
}
