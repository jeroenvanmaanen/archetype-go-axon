package main

import (
    context "context"
    fmt "fmt"
    time "time"
    control "github.com/jeroenvm/archetype-nix-go/src/pkg/control"
    axonserver "github.com/jeroenvm/archetype-nix-go/src/pkg/grpc/axonserver"
    example "github.com/jeroenvm/archetype-nix-go/src/pkg/example"
    grpc "google.golang.org/grpc"
)

func main() {
    clientInfo := axonserver.ClientIdentification {
        ClientId: "12345",
        ComponentName: "GoClient",
        Version: "0.0.1",
    }
    conn, client := waitForServer("example-proxy", 8124, clientInfo)
    defer conn.Close()

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

    // Listen to incoming gRPC requests
    example.Serve(conn)
}

func waitForServer(host string, port int, clientInfo axonserver.ClientIdentification) (*grpc.ClientConn, axonserver.PlatformServiceClient) {
    serverAddress := fmt.Sprintf("%s:%d", host, port)
    fmt.Println(fmt.Sprintf("Client identification: %v", clientInfo))
    d, _ := time.ParseDuration("3s")
    for {
        conn, e := grpc.Dial(serverAddress, grpc.WithInsecure())
        if e == nil {
            // Get platform server
            client := axonserver.NewPlatformServiceClient(conn)
            response, e := client.GetPlatformServer(context.Background(), &clientInfo)
            if e == nil {
                fmt.Println("Connected")
                fmt.Println(response)
                fmt.Println(response.SameConnection)
                if !response.SameConnection {
                    panic(fmt.Sprintf("Need to setup a new connection %v", e))
                }
                return conn, client
            }
        }
        time.Sleep(d)
        fmt.Println(".")
    }
}
