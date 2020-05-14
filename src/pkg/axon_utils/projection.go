package axon_utils

import (
    context "context"
    io "io"
    log "log"

    proto "github.com/golang/protobuf/proto"

    axon_server "github.com/jeroenvm/archetype-go-axon/src/pkg/grpc/axon_server"
)

type Event interface {
    ApplyTo(projection interface{})
}

func RestoreProjection(label string, aggregateIdentifier string, createInitialProjection func()interface{}, clientConnection *ClientConnection, prepareUnmarshal func (payloadType string)Event) interface{} {
    conn := clientConnection.Connection
    projection := createInitialProjection()
    log.Printf("%v Projection: %v", label, projection)

    eventStoreClient := axon_server.NewEventStoreClient(conn)
    requestEvents := axon_server.GetAggregateEventsRequest {
        AggregateId: aggregateIdentifier,
        InitialSequence: 0,
        AllowSnapshots: false,
    }
    log.Printf("%v Projection: Request events: %v", label, requestEvents)
    client, e := eventStoreClient.ListAggregateEvents(context.Background(), &requestEvents)
    if e != nil {
        log.Printf("%v Projection: Error while requesting aggregate events: %v", label, e)
        return nil
    }
    for {
        eventMessage, e := client.Recv()
        if e == io.EOF {
            log.Printf("%v Projection: End of stream", label)
            break
        } else if e != nil {
            log.Printf("%v Projection: Error while receiving next event: %v", label, e)
            break
        }
        log.Printf("%v Projection: Received event: %v", label, eventMessage)
        if eventMessage.Payload != nil {
            payloadType := eventMessage.Payload.Type
            event := prepareUnmarshal(payloadType)
            if event == nil {
                log.Printf("%v Projection: unrecognized payload type: %v", label, payloadType)
                continue
            }
            e := proto.Unmarshal(eventMessage.Payload.Data, event.(proto.Message))
            if (e != nil) {
                log.Printf("%v Projection: Could not unmarshal %v", label, eventMessage.Payload.Type)
            }
            event.ApplyTo(projection)
        }
    }
    return projection
}
