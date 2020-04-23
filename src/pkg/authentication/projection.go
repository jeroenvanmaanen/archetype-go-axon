package authentication

import (
    context "context"
    io "io"
    log "log"

    grpc "google.golang.org/grpc"
    proto "github.com/golang/protobuf/proto"

    axon_server "github.com/jeroenvm/archetype-go-axon/src/pkg/grpc/axon_server"
    grpcExample "github.com/jeroenvm/archetype-go-axon/src/pkg/grpc/example"
)

type Projection struct {
    Credentials map[string]string
}

func RestoreProjection(aggregateIdentifier string, conn *grpc.ClientConn) *Projection {
    projection := Projection{
        Credentials: make(map[string]string),
    }
    log.Printf("Credentials Projection: %v", projection)

    eventStoreClient := axon_server.NewEventStoreClient(conn)
    requestEvents := axon_server.GetAggregateEventsRequest {
        AggregateId: aggregateIdentifier,
        InitialSequence: 0,
        AllowSnapshots: false,
    }
    log.Printf("Credentials Projection: Request events: %v", requestEvents)
    client, e := eventStoreClient.ListAggregateEvents(context.Background(), &requestEvents)
    if e != nil {
        log.Printf("Credentials Projection: Error while requesting aggregate events: %v", e)
        return nil
    }
    for {
        eventMessage, e := client.Recv()
        if e == io.EOF {
            log.Printf("Credentials Projection: End of stream")
            break
        } else if e != nil {
            log.Printf("Credentials Projection: Error while receiving next event: %v", e)
            break
        }
        log.Printf("Credentials Projection: Received event: %v", eventMessage)
        if eventMessage.Payload != nil {
            log.Printf("Credentials Projection: Payload type: %v", eventMessage.Payload.Type)
            payloadType := eventMessage.Payload.Type
            if (payloadType == "CredentialsAddedEvent") {
                event := grpcExample.CredentialsAddedEvent{}
                e := proto.Unmarshal(eventMessage.Payload.Data, &event)
                if (e != nil) {
                    log.Printf("Could not unmarshal CredentialsAddedEvent")
                }
                projection.Credentials[event.Credentials.Identifier] = event.Credentials.Secret
            } else if (payloadType == "CredentialsRemovedEvent") {
                event := grpcExample.CredentialsRemovedEvent{}
                e := proto.Unmarshal(eventMessage.Payload.Data, &event)
                if (e != nil) {
                    log.Printf("Could not unmarshal CredentialsRemovedEvent")
                }
                projection.Credentials[event.Identifier] = ""
            }
        }
    }

    return &projection
}
