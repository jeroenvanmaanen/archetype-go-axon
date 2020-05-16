package example_command

import (
	context "context"
	log "log"

	proto "github.com/golang/protobuf/proto"

	authentication "github.com/jeroenvm/archetype-go-axon/src/pkg/authentication"
	axon_utils "github.com/jeroenvm/archetype-go-axon/src/pkg/axon_utils"
	configuration_command "github.com/jeroenvm/archetype-go-axon/src/pkg/configuration_command"
	axon_server "github.com/jeroenvm/archetype-go-axon/src/pkg/grpc/axon_server"
	grpc_example "github.com/jeroenvm/archetype-go-axon/src/pkg/grpc/example"
	trusted "github.com/jeroenvm/archetype-go-axon/src/pkg/trusted"
)

func HandleCommands(host string, port int) *axon_utils.ClientConnection {
	clientConnection, _ := axon_utils.WaitForServer(host, port, "Command Handler")
	conn := clientConnection.Connection
	clientInfo := clientConnection.ClientInfo

	log.Printf("Command handler: Connection: %v", conn)
	client := axon_server.NewCommandServiceClient(conn)
	log.Printf("Command handler: Client: %v", client)

	stream, e := client.OpenStream(context.Background())
	log.Printf("Command handler: Stream: %v: %v", stream, e)

	axon_utils.SubscribeCommand("GreetCommand", stream, clientInfo)
	axon_utils.SubscribeCommand("RecordCommand", stream, clientInfo)
	axon_utils.SubscribeCommand("StopCommand", stream, clientInfo)
	axon_utils.SubscribeCommand("RegisterTrustedKeyCommand", stream, clientInfo)
	axon_utils.SubscribeCommand("RegisterKeyManagerCommand", stream, clientInfo)
	axon_utils.SubscribeCommand("RegisterCredentialsCommand", stream, clientInfo)
	axon_utils.SubscribeCommand("ChangePropertyCommand", stream, clientInfo)

	go commandWorker(stream, clientConnection)

	return clientConnection
}

func commandWorker(stream axon_server.CommandService_OpenStreamClient, clientConnection *axon_utils.ClientConnection) {
	clientId := clientConnection.ClientInfo.ClientId
	for true {
		axon_utils.CommandAddPermits(1, stream, clientId)

		log.Printf("Command handler: Waiting for command")
		inbound, e := stream.Recv()
		if e != nil {
			log.Printf("Command handler: Error on receive: %v", e)
			break
		}
		log.Printf("Command handler: Inbound: %v", inbound)
		command := inbound.GetCommand()
		if command != nil {
			commandName := command.Name
			log.Printf("Received %v", commandName)
			if commandName == "GreetCommand" {
				handleGreetCommand(command, stream, clientConnection)
			} else if commandName == "RecordCommand" {
				handleRecordCommand(command, stream, clientConnection)
			} else if commandName == "StopCommand" {
				handleStopCommand(command, stream, clientConnection)
			} else if commandName == "RegisterTrustedKeyCommand" {
				trusted.HandleRegisterTrustedKeyCommand(command, stream, clientConnection)
			} else if commandName == "RegisterKeyManagerCommand" {
				trusted.HandleRegisterKeyManagerCommand(command, stream, clientConnection)
			} else if commandName == "RegisterCredentialsCommand" {
				authentication.HandleRegisterCredentialsCommand(command, stream, clientConnection)
			} else if commandName == "ChangePropertyCommand" {
				configuration_command.HandleChangePropertyCommand(command, stream, clientConnection)
			} else {
				log.Printf("Received unknown command: %v", commandName)
			}
		}
	}
}

func handleGreetCommand(command *axon_server.Command, stream axon_server.CommandService_OpenStreamClient, clientConnection *axon_utils.ClientConnection) {
	deserializedCommand := grpc_example.GreetCommand{}
	e := proto.Unmarshal(command.Payload.Data, &deserializedCommand)
	if e != nil {
		log.Printf("Could not unmarshal GreetCommand")
	}

	projection := RestoreProjection(deserializedCommand.AggregateIdentifier, clientConnection)
	if !projection.Recording {
		axon_utils.ReportError(stream, command.MessageIdentifier, "EX001", "Not recording: "+deserializedCommand.AggregateIdentifier)
		return
	}

	event := &GreetedEvent{
		grpc_example.GreetedEvent{
			Message: deserializedCommand.Message,
		},
	}
	axon_utils.AppendEvent(event, deserializedCommand.AggregateIdentifier, projection, clientConnection)
	axon_utils.CommandRespond(stream, command.MessageIdentifier)
}

func handleRecordCommand(command *axon_server.Command, stream axon_server.CommandService_OpenStreamClient, clientConnection *axon_utils.ClientConnection) {
	deserializedCommand := grpc_example.RecordCommand{}
	e := proto.Unmarshal(command.Payload.Data, &deserializedCommand)
	if e != nil {
		log.Printf("Could not unmarshal RecordCommand")
	}
	projection := RestoreProjection(deserializedCommand.AggregateIdentifier, clientConnection)
	if !projection.Recording {
		event := &StartedRecordingEvent{grpc_example.StartedRecordingEvent{}}
		axon_utils.AppendEvent(event, deserializedCommand.AggregateIdentifier, projection, clientConnection)
	}
	axon_utils.CommandRespond(stream, command.MessageIdentifier)
}

func handleStopCommand(command *axon_server.Command, stream axon_server.CommandService_OpenStreamClient, clientConnection *axon_utils.ClientConnection) {
	deserializedCommand := grpc_example.StopCommand{}
	e := proto.Unmarshal(command.Payload.Data, &deserializedCommand)
	if e != nil {
		log.Printf("Could not unmarshal StopCommand")
	}
	projection := RestoreProjection(deserializedCommand.AggregateIdentifier, clientConnection)
	if projection.Recording {
		event := &StoppedRecordingEvent{grpc_example.StoppedRecordingEvent{}}
		axon_utils.AppendEvent(event, deserializedCommand.AggregateIdentifier, projection, clientConnection)
	}
	axon_utils.CommandRespond(stream, command.MessageIdentifier)
}
