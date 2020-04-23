package trusted

import (
    log "log"

    grpc "google.golang.org/grpc"
    proto "github.com/golang/protobuf/proto"

    axon_server "github.com/jeroenvm/archetype-go-axon/src/pkg/grpc/axon_server"
    axon_utils "github.com/jeroenvm/archetype-go-axon/src/pkg/axon_utils"
    grpc_example "github.com/jeroenvm/archetype-go-axon/src/pkg/grpc/example"
)

const AggregateIdentifier = "trusted-keys-aggregate"

func HandleRegisterTrustedKeyCommand(commandMessage *axon_server.Command, stream axon_server.CommandService_OpenStreamClient, conn *grpc.ClientConn) {
    command := grpc_example.RegisterTrustedKeyCommand{}
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
        event := grpc_example.TrustedKeyAddedEvent{
            PublicKey: command.PublicKey,
        }
        log.Printf("Trusted aggregate: emit: %v", event)
        data, e = proto.Marshal(&event)
    } else {
        eventType = "TrustedKeyRemovedEvent"
        event := grpc_example.TrustedKeyRemovedEvent{
            Name: command.PublicKey.Name,
        }
        log.Printf("Trusted aggregate: emit: %v", event)
        data, e = proto.Marshal(&event)
    }

    if e != nil {
        log.Printf("Server: Error while marshalling event: %v", e)
        return
    }
    serializedEvent := axon_server.SerializedObject{
        Type: eventType,
        Data: data,
    }

    axon_utils.AppendEvent(&serializedEvent, AggregateIdentifier, conn)
    axon_utils.CommandRespond(stream, commandMessage.MessageIdentifier)
}

func HandleRegisterKeyManagerCommand(commandMessage *axon_server.Command, stream axon_server.CommandService_OpenStreamClient, conn *grpc.ClientConn) {
    command := grpc_example.RegisterKeyManagerCommand{}
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
        event := grpc_example.KeyManagerAddedEvent{
            PublicKey: command.PublicKey,
        }
        log.Printf("Trusted aggregate: emit: %v", event)
        data, e = proto.Marshal(&event)
    } else {
        eventType = "KeyManagerRemovedEvent"
        event := grpc_example.KeyManagerRemovedEvent{
            Name: command.PublicKey.Name,
        }
        log.Printf("Trusted aggregate: emit: %v", event)
        data, e = proto.Marshal(&event)
    }

    if e != nil {
        log.Printf("Server: Error while marshalling event: %v", e)
        return
    }
    serializedEvent := axon_server.SerializedObject{
        Type: eventType,
        Data: data,
    }

    axon_utils.AppendEvent(&serializedEvent, AggregateIdentifier, conn)
    axon_utils.CommandRespond(stream, commandMessage.MessageIdentifier)
}
