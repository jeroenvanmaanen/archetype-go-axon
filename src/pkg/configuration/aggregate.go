package configuration

import (
    log "log"

    grpc "google.golang.org/grpc"
    proto "github.com/golang/protobuf/proto"

    axon_server "github.com/jeroenvm/archetype-go-axon/src/pkg/grpc/axon_server"
    axon_utils "github.com/jeroenvm/archetype-go-axon/src/pkg/axon_utils"
    grpc_example "github.com/jeroenvm/archetype-go-axon/src/pkg/grpc/example"
)

const AggregateIdentifier = "trusted-keys-aggregate"

func HandleChangePropertyCommand(commandMessage *axon_server.Command, stream axon_server.CommandService_OpenStreamClient, conn *grpc.ClientConn) {
    command := grpc_example.ChangePropertyCommand{}
    e := proto.Unmarshal(commandMessage.Payload.Data, &command)
    if (e != nil) {
        log.Printf("Could not unmarshal ChangePropertyCommand")
        axon_utils.ReportError(stream, commandMessage.MessageIdentifier, "EX001", "Could not unmarshal ChangePropertyCommand")
        return
    }

    projection := RestoreProjection(AggregateIdentifier, conn)

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
        axon_utils.AppendEvent(&serializedEvent, AggregateIdentifier, conn)
    }
    axon_utils.CommandRespond(stream, commandMessage.MessageIdentifier)
}