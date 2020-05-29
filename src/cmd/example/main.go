package main

import (
	fmt "fmt"
	log "log"

	uuid "github.com/google/uuid"

	authentication "github.com/jeroenvanmaanen/archetype-go-axon/src/pkg/authentication"
	cache_utils "github.com/jeroenvanmaanen/archetype-go-axon/src/pkg/cache_utils"
	configuration_query "github.com/jeroenvanmaanen/archetype-go-axon/src/pkg/configuration_query"
	example_api "github.com/jeroenvanmaanen/archetype-go-axon/src/pkg/example_api"
	example_command "github.com/jeroenvanmaanen/archetype-go-axon/src/pkg/example_command"
	example_query "github.com/jeroenvanmaanen/archetype-go-axon/src/pkg/example_query"
	trusted "github.com/jeroenvanmaanen/archetype-go-axon/src/pkg/trusted"
	utils "github.com/jeroenvanmaanen/archetype-go-axon/src/pkg/utils"
	axon_utils "github.com/jeroenvanmaanen/dendrite/src/pkg/axon_utils"
	axon_server "github.com/jeroenvanmaanen/dendrite/src/pkg/grpc/axon_server"
)

func main() {
	log.Printf("\n\n\n")
	log.Printf("Start Go Client")

	trusted.Init()
	authentication.Init()
	for k, v := range trusted.GetTrustedKeys() {
		log.Printf("Trusted key: %v: %v", k, v)
	}

	host := "axon-server" // "example-proxy" or "axon-server"
	port := 8124
	clientConnection, streamClient := axon_utils.WaitForServer(host, port, "API")
	defer utils.ReportError("Close clientConnection", clientConnection.Connection.Close)
	log.Printf("Main connection: %v: %v", clientConnection, streamClient)

	// Send a heartbeat
	heartbeat := axon_server.Heartbeat{}
	heartbeatRequest := axon_server.PlatformInboundInstruction_Heartbeat{
		Heartbeat: &heartbeat,
	}
	id := uuid.New()
	instruction := axon_server.PlatformInboundInstruction{
		Request:       &heartbeatRequest,
		InstructionId: id.String(),
	}
	if e := (*streamClient).Send(&instruction); e != nil {
		panic(fmt.Sprintf("Error sending clientInfo %v", e))
	}

	// Initialize cache
	cache_utils.InitializeCache()

	// Handle commands
	commandHandlerConn := example_command.HandleCommands(host, port)
	defer utils.ReportError("Close commandHandlerConn", commandHandlerConn.Connection.Close)

	// Process Events
	eventProcessorConn := example_query.ProcessEvents(host, port)
	defer utils.ReportError("Close eventProcessorConn", eventProcessorConn.Connection.Close)

	configurationEventProcessorConn := configuration_query.ProcessEvents(host, port)
	defer utils.ReportError("Close configurationEventProcessorConn", configurationEventProcessorConn.Connection.Close)

	// Handle queries
	queryHandlerConn := example_query.HandleQueries(host, port)
	defer utils.ReportError("Close queryHandlerConn", queryHandlerConn.Connection.Close)

	// Listen to incoming gRPC requests
	_ = axon_utils.Serve(clientConnection, example_api.RegisterWithServer)
}
