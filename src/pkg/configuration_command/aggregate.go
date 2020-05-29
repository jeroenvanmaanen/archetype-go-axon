package configuration_command

import (
	log "log"

	proto "github.com/golang/protobuf/proto"

	grpc_example "github.com/jeroenvanmaanen/archetype-go-axon/src/pkg/grpc/example"
	axon_utils "github.com/jeroenvanmaanen/dendrite/src/pkg/axon_utils"
	axon_server "github.com/jeroenvanmaanen/dendrite/src/pkg/grpc/axon_server"
)

const AggregateIdentifier = "configuration-aggregate"

func HandleChangePropertyCommand(commandMessage *axon_server.Command, clientConnection *axon_utils.ClientConnection) (*axon_utils.Error, error) {
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
