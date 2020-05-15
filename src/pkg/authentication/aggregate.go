package authentication

import (
    log "log"

    proto "github.com/golang/protobuf/proto"

    axon_server "github.com/jeroenvm/archetype-go-axon/src/pkg/grpc/axon_server"
    axon_utils "github.com/jeroenvm/archetype-go-axon/src/pkg/axon_utils"
    grpc_example "github.com/jeroenvm/archetype-go-axon/src/pkg/grpc/example"
)

const AggregateIdentifier = "credentials-aggregate"

func HandleRegisterCredentialsCommand(commandMessage *axon_server.Command, stream axon_server.CommandService_OpenStreamClient, clientConnection *axon_utils.ClientConnection) {
    command := grpc_example.RegisterCredentialsCommand{}
    e := proto.Unmarshal(commandMessage.Payload.Data, &command)
    if (e != nil) {
        log.Printf("Could not unmarshal RegisterCredentialsCommand")
        axon_utils.ReportError(stream, commandMessage.MessageIdentifier, "EX001", "Could not unmarshal RegisterCredentialsCommand")
        return
    }

    projection := RestoreProjection(AggregateIdentifier, clientConnection)

    if CheckKnown(command.Credentials, projection) {
        axon_utils.CommandRespond(stream, commandMessage.MessageIdentifier)
        return
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
    axon_utils.AppendEvent(event, AggregateIdentifier, projection, clientConnection)
    axon_utils.CommandRespond(stream, commandMessage.MessageIdentifier)
}
