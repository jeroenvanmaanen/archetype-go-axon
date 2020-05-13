package main

import (
    fmt "fmt"
    log "log"

    uuid "github.com/google/uuid"

    authentication "github.com/jeroenvm/archetype-go-axon/src/pkg/authentication"
    axon_server "github.com/jeroenvm/archetype-go-axon/src/pkg/grpc/axon_server"
    axon_utils "github.com/jeroenvm/archetype-go-axon/src/pkg/axon_utils"
    configuration_query "github.com/jeroenvm/archetype-go-axon/src/pkg/configuration_query"
    example "github.com/jeroenvm/archetype-go-axon/src/pkg/example"
    trusted "github.com/jeroenvm/archetype-go-axon/src/pkg/trusted"
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
    defer clientConnection.Connection.Close()
    log.Printf("Main connection: %v: %v", clientConnection, streamClient)

    // Send a heartbeat
    heartbeat := axon_server.Heartbeat{}
    heartbeatRequest := axon_server.PlatformInboundInstruction_Heartbeat{
        Heartbeat: &heartbeat,
    }
    id := uuid.New()
    instruction := axon_server.PlatformInboundInstruction {
        Request: &heartbeatRequest,
        InstructionId: id.String(),
    }
    if e := (*streamClient).Send(&instruction); e != nil {
        panic(fmt.Sprintf("Error sending clientInfo %v", e))
    }

    // Handle commands
    commandHandlerConn := example.HandleCommands(host, port)
    defer commandHandlerConn.Connection.Close()

    // Process Events
    eventProcessorConn := example.ProcessEvents(host, port)
    defer eventProcessorConn.Connection.Close()

    configurationEventProcessorConn := configuration_query.ProcessEvents(host, port)
    defer configurationEventProcessorConn.Connection.Close()

    // Handle queries
    queryHandlerConn := example.HandleQueries(host, port)
    defer queryHandlerConn.Connection.Close()

    // Listen to incoming gRPC requests
    axon_utils.Serve(clientConnection, example.RegisterWithServer)
}
