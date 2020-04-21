package authentication

import (
    log "log"

    grpc "google.golang.org/grpc"
    proto "github.com/golang/protobuf/proto"

    axonserver "github.com/jeroenvm/archetype-go-axon/src/pkg/grpc/axonserver"
    axonutils "github.com/jeroenvm/archetype-go-axon/src/pkg/axonutils"
    grpcExample "github.com/jeroenvm/archetype-go-axon/src/pkg/grpc/example"
)

const AggregateIdentifier = "credentials-aggregate"

func HandleRegisterCredentialsCommand(commandMessage *axonserver.Command, stream axonserver.CommandService_OpenStreamClient, conn *grpc.ClientConn) {
    command := grpcExample.RegisterCredentialsCommand{}
    e := proto.Unmarshal(commandMessage.Payload.Data, &command)
    if (e != nil) {
        log.Printf("Could not unmarshal RegisterCredentialsCommand")
        return
    }

    projection := RestoreProjection(AggregateIdentifier, conn)

    currentValue := projection.Credentials[command.Credentials.Identifier]
    newValue := command.Credentials.Secret
    if newValue == currentValue {
        return
    }

    var eventType string
    var data []byte
    if len(newValue) > 0 {
        eventType = "CredentialsAddedEvent"
        event := grpcExample.CredentialsAddedEvent{
            Credentials: command.Credentials,
        }
        log.Printf("Credentials aggregate: emit: %v", event)
        data, e = proto.Marshal(&event)
    } else {
        eventType = "CredentialsRemovedEvent"
        event := grpcExample.CredentialsRemovedEvent{
            Identifier: command.Credentials.Identifier,
        }
        log.Printf("Credentials aggregate: emit: %v", event)
        data, e = proto.Marshal(&event)
    }

    if e != nil {
        log.Printf("Server: Error while marshalling event: %v", e)
        return
    }
    serializedEvent := axonserver.SerializedObject{
        Type: eventType,
        Data: data,
    }

    axonutils.AppendEvent(&serializedEvent, AggregateIdentifier, conn)
    axonutils.CommandRespond(stream, commandMessage.MessageIdentifier)
}
