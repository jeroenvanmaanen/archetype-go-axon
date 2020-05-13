package configuration_command

import (
    log "log"

    proto "github.com/golang/protobuf/proto"

    axon_server "github.com/jeroenvm/archetype-go-axon/src/pkg/grpc/axon_server"
    axon_utils "github.com/jeroenvm/archetype-go-axon/src/pkg/axon_utils"
    grpc_example "github.com/jeroenvm/archetype-go-axon/src/pkg/grpc/example"
)

const AggregateIdentifier = "configuration-aggregate"

func HandleChangePropertyCommand(commandMessage *axon_server.Command, stream axon_server.CommandService_OpenStreamClient, clientConnection *axon_utils.ClientConnection) {
    command := grpc_example.ChangePropertyCommand{}
    e := proto.Unmarshal(commandMessage.Payload.Data, &command)
    if (e != nil) {
        log.Printf("Could not unmarshal ChangePropertyCommand")
        axon_utils.ReportError(stream, commandMessage.MessageIdentifier, "EX001", "Could not unmarshal ChangePropertyCommand")
        return
    }

    projection := RestoreProjection(AggregateIdentifier, clientConnection)

    key := command.Property.Key
    newValue := command.Property.Value
    oldValue := projection.Configuration[key]

    if newValue != oldValue {
        event := grpc_example.PropertyChangedEvent{
            Property: command.Property,
        }
        log.Printf("Trusted aggregate: emit: %v", event)
        data, e := proto.Marshal(&event)
        if e != nil {
            log.Printf("Server: Error while marshalling event: %v", e)
            return
        }
        serializedEvent := axon_server.SerializedObject{
            Type: "PropertyChangedEvent",
            Data: data,
        }
        axon_utils.AppendEvent(&serializedEvent, AggregateIdentifier, clientConnection)
    }
    axon_utils.CommandRespond(stream, commandMessage.MessageIdentifier)
}