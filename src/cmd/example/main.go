package main

import (
    fmt "fmt"
    log "log"
    authentication "github.com/jeroenvm/archetype-go-axon/src/pkg/authentication"
    axonserver "github.com/jeroenvm/archetype-go-axon/src/pkg/grpc/axonserver"
    example "github.com/jeroenvm/archetype-go-axon/src/pkg/example"
    trusted "github.com/jeroenvm/archetype-go-axon/src/pkg/trusted"
    uuid "github.com/google/uuid"
)

func main() {
    log.Printf("\n\n\n")
    log.Printf("Start Go Client")

    trusted.Init()
    authentication.Init()
    for k, v := range trusted.TrustedKeys {
        log.Printf("Trusted key: %v: %v", k, v)
    }

    host := "axon-server" // "example-proxy" or "axon-server"
    port := 8124
    conn, clientInfo, streamClient := example.WaitForServer(host, port, "API")
    defer conn.Close()
    log.Printf("Main connection: %v: %v: %v", conn, clientInfo, streamClient)

    // Listen to messages from Axon Server in a separate go routine
    // go control.Listen(streamClient)

    // Send a heartbeat
    heartbeat := axonserver.Heartbeat{}
    heartbeatRequest := axonserver.PlatformInboundInstruction_Heartbeat{
        Heartbeat: &heartbeat,
    }
    id := uuid.New()
    instruction := axonserver.PlatformInboundInstruction {
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

    // Handle queries
    queryHandlerConn := example.HandleQueries(host, port)
    defer queryHandlerConn.Close()

    // Listen to incoming gRPC requests
    example.Serve(conn, clientInfo)
}
