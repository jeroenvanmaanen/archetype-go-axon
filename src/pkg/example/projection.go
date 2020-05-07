package example

import (
    context "context"
    io "io"
    log "log"

    axon_server "github.com/jeroenvm/archetype-go-axon/src/pkg/grpc/axon_server"
    axon_utils "github.com/jeroenvm/archetype-go-axon/src/pkg/axon_utils"
)

type Projection struct {
    Recording bool
}

func RestoreProjection(aggregateIdentifier string, clientConnection *axon_utils.ClientConnection) *Projection {
    conn := clientConnection.Connection
    projection := Projection{
        Recording: true,
    }
    log.Printf("Projection: %v", projection)

    eventStoreClient := axon_server.NewEventStoreClient(conn)
    requestEvents := axon_server.GetAggregateEventsRequest {
        AggregateId: aggregateIdentifier,
        InitialSequence: 0,
        AllowSnapshots: false,
    }
    log.Printf("Projection: Request events: %v", requestEvents)
    client, e := eventStoreClient.ListAggregateEvents(context.Background(), &requestEvents)
    if e != nil {
        log.Printf("Projection: Error while requesting aggregate events: %v", e)
        return nil
    }
    for {
        eventMessage, e := client.Recv()
        if e == io.EOF {
            log.Printf("Projection: End of stream")
            break
        } else if e != nil {
            log.Printf("Projection: Error while receiving next event: %v", e)
            break
        }
        log.Printf("Projection: Received event: %v", eventMessage)
        if eventMessage.Payload != nil {
            log.Printf("Projection: Payload type: %v", eventMessage.Payload.Type)
            payloadType := eventMessage.Payload.Type
            if (payloadType == "StartedRecordingEvent") {
                projection.Recording = true
            } else if (payloadType == "StoppedRecordingEvent") {
                projection.Recording = false
            }
        }
    }

    return &projection
}
