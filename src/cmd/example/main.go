package main

import (
    context "context"
    "fmt"
    "time"
    control "github.com/jeroenvm/archetype-nix-go/src/pkg/control"
    axonserver "github.com/jeroenvm/archetype-nix-go/src/pkg/grpc/axonserver"
    grpc "google.golang.org/grpc"
)

var _ axonserver.PlatformOutboundInstruction

func main() {
    serverAddress := "axon-server:8124"
    conn, e := grpc.Dial(serverAddress, grpc.WithInsecure())

    if e != nil {
        panic(e)
    }
    defer conn.Close()

    // Get platform server
    client := axonserver.NewPlatformServiceClient(conn)
    clientInfo := axonserver.ClientIdentification {
        ClientId: "12345",
        ComponentName: "GoClient",
        Version: "0.0.1",
    }
    fmt.Println(fmt.Sprintf("Client identification: %v", clientInfo))
    response, e := client.GetPlatformServer(context.Background(), &clientInfo)
    if e != nil {
        panic(fmt.Sprintf("Was not able to get Axon platform server %v", e))
    }
    fmt.Println(response)
    fmt.Println(response.SameConnection)
    if !response.SameConnection {
        panic(fmt.Sprintf("Need to setup a new connection %v", e))
    }

    // Open stream
    streamClient, e := client.OpenStream(context.Background())
    if e != nil {
        panic(fmt.Sprintf("Could not open stream %v", e))
    }

    // Send client info
    var instruction axonserver.PlatformInboundInstruction
    registrationRequest := axonserver.PlatformInboundInstruction_Register{
        Register: &clientInfo,
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

    d, e := time.ParseDuration("10s")
    time.Sleep(d)
}
