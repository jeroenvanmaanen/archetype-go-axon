package trusted

import (
    log "log"

    grpc "google.golang.org/grpc"
    proto "github.com/golang/protobuf/proto"

    axonserver "github.com/jeroenvm/archetype-go-axon/src/pkg/grpc/axonserver"
    axonutils "github.com/jeroenvm/archetype-go-axon/src/pkg/axonutils"
    grpcExample "github.com/jeroenvm/archetype-go-axon/src/pkg/grpc/example"
)

const AggregateIdentifier = "trusted-keys-aggregate"

func HandleRegisterTrustedKeyCommand(commandMessage *axonserver.Command, stream axonserver.CommandService_OpenStreamClient, conn *grpc.ClientConn) {
    command := grpcExample.RegisterTrustedKeyCommand{}
    e := proto.Unmarshal(commandMessage.Payload.Data, &command)
    if (e != nil) {
        log.Printf("Could not unmarshal RegisterTrustedKeyCommand")
    }

    projection := RestoreProjection(AggregateIdentifier, conn)

    currentValue := projection.TrustedKeys[command.PublicKey.Name]
    newValue := command.PublicKey.PublicKey
    if newValue == currentValue {
        return
    }

    var eventType string
    var data []byte
    if len(newValue) > 0 {
        eventType = "TrustedKeyAddedEvent"
        event := grpcExample.TrustedKeyAddedEvent{
            PublicKey: command.PublicKey,
        }
        log.Printf("Trusted aggregate: emit: %v", event)
        data, e = proto.Marshal(&event)
    } else {
        eventType = "TrustedKeyRemovedEvent"
        event := grpcExample.TrustedKeyRemovedEvent{
            Name: command.PublicKey.Name,
        }
        log.Printf("Trusted aggregate: emit: %v", event)
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

func HandleRegisterKeyManagerCommand(commandMessage *axonserver.Command, stream axonserver.CommandService_OpenStreamClient, conn *grpc.ClientConn) {
    command := grpcExample.RegisterKeyManagerCommand{}
    e := proto.Unmarshal(commandMessage.Payload.Data, &command)
    if (e != nil) {
        log.Printf("Could not unmarshal RegisterKeyManagerCommand")
    }

    projection := RestoreProjection(AggregateIdentifier, conn)

    currentValue := projection.KeyManagers[command.PublicKey.Name]
    newValue := command.PublicKey.PublicKey
    if newValue == currentValue {
        return
    }

    var eventType string
    var data []byte
    if len(newValue) > 0 {
        eventType = "KeyMangerAddedEvent"
        event := grpcExample.KeyManagerAddedEvent{
            PublicKey: command.PublicKey,
        }
        log.Printf("Trusted aggregate: emit: %v", event)
        data, e = proto.Marshal(&event)
    } else {
        eventType = "KeyManagerRemovedEvent"
        event := grpcExample.KeyManagerRemovedEvent{
            Name: command.PublicKey.Name,
        }
        log.Printf("Trusted aggregate: emit: %v", event)
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