package main

import (
    context "context"
    fmt "fmt"
    log "log"
    control "github.com/jeroenvm/archetype-nix-go/src/pkg/control"
    axonserver "github.com/jeroenvm/archetype-nix-go/src/pkg/grpc/axonserver"
    example "github.com/jeroenvm/archetype-nix-go/src/pkg/example"
)

func main() {
    host := "example-proxy" // "example-proxy" or "axon-server"
    port := 8124
    conn, clientInfo := example.WaitForServerAndReturnClientInfo(host, port, "API")
    defer conn.Close()

    log.Printf("Main connection: %v", conn)

    // Open stream
    client := axonserver.NewPlatformServiceClient(conn)
    streamClient, e := client.OpenStream(context.Background())
    if e != nil {
        panic(fmt.Sprintf("Could not open stream %v", e))
    }

    // Send client info
    var instruction axonserver.PlatformInboundInstruction
    registrationRequest := axonserver.PlatformInboundInstruction_Register{
        Register: clientInfo,
    }
    instruction.Request = &registrationRequest
    if e = streamClient.Send(&instruction); e != nil {
        panic(fmt.Sprintf("Error sending clientInfo %v", e))
    }

    // Listen to messages from Axon Server in a separate go routine
    go control.Listen(&streamClient)

    // Send a heartbeat
    heartbeat := axonserver.Heartbeat{}
    heartbeatRequest := axonserver.PlatformInboundInstruction_Heartbeat{
        Heartbeat: &heartbeat,
    }
    instruction.Request = &heartbeatRequest
    if e = streamClient.Send(&instruction); e != nil {
        panic(fmt.Sprintf("Error sending clientInfo %v", e))
    }

    // Handle commands
    commandHandlerConn := example.HandleCommands(host, port)
    defer commandHandlerConn.Close()

    // Listen to incoming gRPC requests
    example.Serve(conn, clientInfo)
}
