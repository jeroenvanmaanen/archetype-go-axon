package trusted

import (
	log "log"

	proto "github.com/golang/protobuf/proto"

	axon_utils "github.com/jeroenvanmaanen/dendrite/src/pkg/axon_utils"
	axon_server "github.com/jeroenvanmaanen/dendrite/src/pkg/grpc/axon_server"
	grpc_example "github.com/jeroenvm/archetype-go-axon/src/pkg/grpc/example"
)

const AggregateIdentifier = "trusted-keys-aggregate"

func HandleRegisterTrustedKeyCommand(commandMessage *axon_server.Command, stream axon_server.CommandService_OpenStreamClient, clientConnection *axon_utils.ClientConnection) (*axon_utils.Error, error) {
	command := grpc_example.RegisterTrustedKeyCommand{}
	e := proto.Unmarshal(commandMessage.Payload.Data, &command)
	if e != nil {
		log.Printf("Could not unmarshal RegisterTrustedKeyCommand")
		return nil, e
	}

	projection := RestoreProjection(AggregateIdentifier, clientConnection)

	currentValue := projection.TrustedKeys[command.PublicKey.Name]
	newValue := command.PublicKey.PublicKey
	if newValue == currentValue {
		return nil, nil
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
	return axon_utils.AppendEvent(event, AggregateIdentifier, projection, clientConnection)
}

func HandleRegisterKeyManagerCommand(commandMessage *axon_server.Command, stream axon_server.CommandService_OpenStreamClient, clientConnection *axon_utils.ClientConnection) (*axon_utils.Error, error) {
	command := grpc_example.RegisterKeyManagerCommand{}
	e := proto.Unmarshal(commandMessage.Payload.Data, &command)
	if e != nil {
		log.Printf("Could not unmarshal RegisterKeyManagerCommand")
		return nil, e
	}

	projection := RestoreProjection(AggregateIdentifier, clientConnection)

	currentValue := projection.KeyManagers[command.PublicKey.Name]
	newValue := command.PublicKey.PublicKey
	if newValue == currentValue {
		return nil, nil
	}

	var eventType string
	var event axon_utils.Event
	if len(newValue) > 0 {
		eventType = "KeyManagerAddedEvent"
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
	return axon_utils.AppendEvent(event, AggregateIdentifier, projection, clientConnection)
}
