package configuration

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
    Configuration map[string]string
}

func RestoreProjection(aggregateIdentifier string, clientConnection *axon_utils.ClientConnection) *Projection {
    conn := clientConnection.Connection
    projection := Projection{
        Configuration: make(map[string]string),
    }
    log.Printf("Configuration Projection: %v", projection)

    eventStoreClient := axon_server.NewEventStoreClient(conn)
    requestEvents := axon_server.GetAggregateEventsRequest {
        AggregateId: aggregateIdentifier,
        InitialSequence: 0,
        AllowSnapshots: false,
    }
    log.Printf("Configuration Projection: Request events: %v", requestEvents)
    client, e := eventStoreClient.ListAggregateEvents(context.Background(), &requestEvents)
    if e != nil {
        log.Printf("Configuration Projection: Error while requesting aggregate events: %v", e)
        return nil
    }
    for {
        eventMessage, e := client.Recv()
        if e == io.EOF {
            log.Printf("Configuration Projection: End of stream")
            break
        } else if e != nil {
            log.Printf("Configuration Projection: Error while receiving next event: %v", e)
            break
        }
        log.Printf("Configuration Projection: Received event: %v", eventMessage)
        if eventMessage.Payload != nil {
            log.Printf("Configuration Projection: Payload type: %v", eventMessage.Payload.Type)
            payloadType := eventMessage.Payload.Type
            if (payloadType == "PropertyChangedEvent") {
                event := grpc_example.PropertyChangedEvent{}
                e := proto.Unmarshal(eventMessage.Payload.Data, &event)
                if (e != nil) {
                    log.Printf("Could not unmarshal PropertyChangedEvent")
                }
                projection.Configuration[event.Property.Key] = event.Property.Value
            }
        }
    }

    return &projection
}
