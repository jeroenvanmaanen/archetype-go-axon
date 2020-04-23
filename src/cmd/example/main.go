package main

import (
    fmt "fmt"
    log "log"

    uuid "github.com/google/uuid"

    authentication "github.com/jeroenvm/archetype-go-axon/src/pkg/authentication"
    axon_server "github.com/jeroenvm/archetype-go-axon/src/pkg/grpc/axon_server"
    axon_utils "github.com/jeroenvm/archetype-go-axon/src/pkg/axon_utils"
    configuration "github.com/jeroenvm/archetype-go-axon/src/pkg/configuration"
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
    conn, clientInfo, streamClient := axon_utils.WaitForServer(host, port, "API")
    defer conn.Close()
    log.Printf("Main connection: %v: %v: %v", conn, clientInfo, streamClient)

    // Listen to messages from Axon Server in a separate go routine
    // go control.Listen(streamClient)

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
    defer commandHandlerConn.Close()

    // Process Events
    eventProcessorConn := example.ProcessEvents(host, port)
    defer eventProcessorConn.Close()

    configurationEventProcessorConn := configuration.ProcessEvents(host, port)
    defer configurationEventProcessorConn.Close()

    // Handle queries
    queryHandlerConn := example.HandleQueries(host, port)
    defer queryHandlerConn.Close()

    // Listen to incoming gRPC requests
    example.Serve(conn, clientInfo)
}
