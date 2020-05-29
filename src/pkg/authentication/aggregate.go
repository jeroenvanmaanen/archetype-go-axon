package authentication

import (
	log "log"

	proto "github.com/golang/protobuf/proto"

	grpc_example "github.com/jeroenvanmaanen/archetype-go-axon/src/pkg/grpc/example"
	axon_utils "github.com/jeroenvanmaanen/dendrite/src/pkg/axon_utils"
	axon_server "github.com/jeroenvanmaanen/dendrite/src/pkg/grpc/axon_server"
)

const AggregateIdentifier = "credentials-aggregate"

func HandleRegisterCredentialsCommand(commandMessage *axon_server.Command, clientConnection *axon_utils.ClientConnection) (*axon_utils.Error, error) {
	command := grpc_example.RegisterCredentialsCommand{}
	e := proto.Unmarshal(commandMessage.Payload.Data, &command)
	if e != nil {
		log.Printf("Could not unmarshal RegisterCredentialsCommand")
		return nil, e
	}

	projection := RestoreProjection(AggregateIdentifier, clientConnection)

	if CheckKnown(command.Credentials, projection) {
		return nil, nil
	}

	var eventType string
	var event axon_utils.Event
	if len(command.Credentials.Secret) > 0 {
		eventType = "CredentialsAddedEvent"
		event = &CredentialsAddedEvent{
			grpc_example.CredentialsAddedEvent{
				Credentials: command.Credentials,
			},
		}
	} else {
		eventType = "CredentialsRemovedEvent"
		event = &CredentialsRemovedEvent{
			grpc_example.CredentialsRemovedEvent{
				Identifier: command.Credentials.Identifier,
			},
		}
	}
	log.Printf("Credentials aggregate: emit: %v: %v", eventType, event)
	return axon_utils.AppendEvent(event, AggregateIdentifier, projection, clientConnection)
}
