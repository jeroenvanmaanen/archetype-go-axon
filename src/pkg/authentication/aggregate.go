package authentication

import (
	"errors"
	log "log"

	proto "github.com/golang/protobuf/proto"

	axon_utils "github.com/jeroenvanmaanen/dendrite/src/pkg/axon_utils"
	axon_server "github.com/jeroenvanmaanen/dendrite/src/pkg/grpc/axon_server"
	grpc_example "github.com/jeroenvm/archetype-go-axon/src/pkg/grpc/example"
)

const AggregateIdentifier = "credentials-aggregate"

func HandleRegisterCredentialsCommand(commandMessage *axon_server.Command, stream axon_server.CommandService_OpenStreamClient, clientConnection *axon_utils.ClientConnection) (*axon_utils.Error, error) {
	command := grpc_example.RegisterCredentialsCommand{}
	e := proto.Unmarshal(commandMessage.Payload.Data, &command)
	if e != nil {
		log.Printf("Could not unmarshal RegisterCredentialsCommand")
		return nil, e
	}

	projection := RestoreProjection(AggregateIdentifier, clientConnection)

	if CheckKnown(command.Credentials, projection) {
		axon_utils.CommandRespond(stream, commandMessage.MessageIdentifier)
		return nil, errors.New("credentials unknown")
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
