package configuration_command

import (
	log "log"

	proto "github.com/golang/protobuf/proto"

	axon_utils "github.com/jeroenvm/archetype-go-axon/src/pkg/axon_utils"
	axon_server "github.com/jeroenvm/archetype-go-axon/src/pkg/grpc/axon_server"
	grpc_example "github.com/jeroenvm/archetype-go-axon/src/pkg/grpc/example"
)

const AggregateIdentifier = "configuration-aggregate"

func HandleChangePropertyCommand(commandMessage *axon_server.Command, stream axon_server.CommandService_OpenStreamClient, clientConnection *axon_utils.ClientConnection) (*axon_utils.Error, error) {
	command := grpc_example.ChangePropertyCommand{}
	e := proto.Unmarshal(commandMessage.Payload.Data, &command)
	if e != nil {
		log.Printf("Could not unmarshal ChangePropertyCommand")
		return nil, e
	}

	projection := RestoreProjection(AggregateIdentifier, clientConnection)

	key := command.Property.Key
	newValue := command.Property.Value
	oldValue := projection.Configuration[key]

	if newValue != oldValue {
		event := &PropertyChangedEvent{
			grpc_example.PropertyChangedEvent{
				Property: command.Property,
			},
		}
		log.Printf("Trusted aggregate: emit: %v", event)
		return axon_utils.AppendEvent(event, AggregateIdentifier, projection, clientConnection)
	}
	return nil, nil
}
