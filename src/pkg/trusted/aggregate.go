package trusted

import (
    log "log"

    proto "github.com/golang/protobuf/proto"

    axon_server "github.com/jeroenvm/archetype-go-axon/src/pkg/grpc/axon_server"
    axon_utils "github.com/jeroenvm/archetype-go-axon/src/pkg/axon_utils"
    grpc_example "github.com/jeroenvm/archetype-go-axon/src/pkg/grpc/example"
)

const AggregateIdentifier = "trusted-keys-aggregate"

func HandleRegisterTrustedKeyCommand(commandMessage *axon_server.Command, stream axon_server.CommandService_OpenStreamClient, clientConnection *axon_utils.ClientConnection) {
    command := grpc_example.RegisterTrustedKeyCommand{}
    e := proto.Unmarshal(commandMessage.Payload.Data, &command)
    if (e != nil) {
        log.Printf("Could not unmarshal RegisterTrustedKeyCommand")
    }

    projection := RestoreProjection(AggregateIdentifier, clientConnection)

    currentValue := projection.TrustedKeys[command.PublicKey.Name]
    newValue := command.PublicKey.PublicKey
    if newValue == currentValue {
        axon_utils.CommandRespond(stream, commandMessage.MessageIdentifier)
        return
    }

    var eventType string
    var event axon_utils.Event
    if len(newValue) > 0 {
        eventType = "TrustedKeyAddedEvent"
        event = &TrustedKeyAddedEvent{
            grpc_example.TrustedKeyAddedEvent{
                PublicKey: command.PublicKey,
            },
        }
    } else {
        eventType = "TrustedKeyRemovedEvent"
        event = &TrustedKeyRemovedEvent{
            grpc_example.TrustedKeyRemovedEvent{
                Name: command.PublicKey.Name,
            },
        }
    }
    log.Printf("Trusted aggregate: emit: %v: %v", eventType, event)
    axon_utils.AppendEvent(event, AggregateIdentifier, projection, clientConnection)
    axon_utils.CommandRespond(stream, commandMessage.MessageIdentifier)
}

func HandleRegisterKeyManagerCommand(commandMessage *axon_server.Command, stream axon_server.CommandService_OpenStreamClient, clientConnection *axon_utils.ClientConnection) {
    command := grpc_example.RegisterKeyManagerCommand{}
    e := proto.Unmarshal(commandMessage.Payload.Data, &command)
    if (e != nil) {
        log.Printf("Could not unmarshal RegisterKeyManagerCommand")
    }

    projection := RestoreProjection(AggregateIdentifier, clientConnection)

    currentValue := projection.KeyManagers[command.PublicKey.Name]
    newValue := command.PublicKey.PublicKey
    if newValue == currentValue {
        return
    }

    var eventType string
    var event axon_utils.Event
    if len(newValue) > 0 {
        eventType = "KeyMangerAddedEvent"
        event = &KeyManagerAddedEvent{
            grpc_example.KeyManagerAddedEvent{
                PublicKey: command.PublicKey,
            },
        }
    } else {
        eventType = "KeyManagerRemovedEvent"
        event = &KeyManagerRemovedEvent{
            grpc_example.KeyManagerRemovedEvent{
                Name: command.PublicKey.Name,
            },
        }
    }
    log.Printf("Trusted aggregate: emit: %v: %v", eventType, event)
    axon_utils.AppendEvent(event, AggregateIdentifier, projection, clientConnection)
    axon_utils.CommandRespond(stream, commandMessage.MessageIdentifier)
}
