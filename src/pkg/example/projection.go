package example

import (
    context "context"
    io "io"
    log "log"
    axonserver "github.com/jeroenvm/archetype-nix-go/src/pkg/grpc/axonserver"
    grpc "google.golang.org/grpc"
)

type Projection struct {
    Recording bool
}

func RestoreProjection(aggregateIdentifier string, conn *grpc.ClientConn) *Projection {
    projection := Projection{
        Recording: true,
    }
    log.Printf("Projection: %v", projection)

    eventStoreClient := axonserver.NewEventStoreClient(conn)
    requestEvents := axonserver.GetAggregateEventsRequest {
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
