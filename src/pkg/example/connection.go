package example

import (
    context "context"
    fmt "fmt"
    log "log"
    time "time"
    axonserver "github.com/jeroenvm/archetype-nix-go/src/pkg/grpc/axonserver"
    grpc "google.golang.org/grpc"
    uuid "github.com/google/uuid"
)

func WaitForServer(host string, port int, qualifier string) (*grpc.ClientConn, *axonserver.ClientIdentification) {
    conn, clientInfo := WaitForServerAndReturnClientInfo(host, port, qualifier)

    // Open stream
    client := axonserver.NewPlatformServiceClient(conn)
    streamClient, e := client.OpenStream(context.Background())
    if e != nil {
        panic(fmt.Sprintf("Connection: Could not open stream %v", e))
    }

    // Send client info
    var instruction axonserver.PlatformInboundInstruction
    registrationRequest := axonserver.PlatformInboundInstruction_Register{
        Register: clientInfo,
    }
    instruction.Request = &registrationRequest
    if e = streamClient.Send(&instruction); e != nil {
        panic(fmt.Sprintf("Connection: Error sending clientInfo %v", e))
    }

    return conn, clientInfo;
}

func WaitForServerAndReturnClientInfo(host string, port int, qualifier string) (*grpc.ClientConn, *axonserver.ClientIdentification) {
    id := uuid.New()
    clientInfo := axonserver.ClientIdentification {
        ClientId: id.String(),
        ComponentName: "GoClient " + qualifier,
        Version: "0.0.1",
    }

    serverAddress := fmt.Sprintf("%s:%d", host, port)
    log.Printf("Connection: Client identification: %v", clientInfo)
    d, _ := time.ParseDuration("3s")
    for {
        conn, e := grpc.Dial(serverAddress, grpc.WithInsecure())
        if e == nil {
            // Get platform server
            client := axonserver.NewPlatformServiceClient(conn)
            response, e := client.GetPlatformServer(context.Background(), &clientInfo)
            if e == nil {
                log.Printf("Connection: Connected: %v: %v", response.SameConnection, response)
                if !response.SameConnection {
                    panic(fmt.Sprintf("Connection: Need to setup a new connection %v", e))
                }
                return conn, &clientInfo
            }
        }
        time.Sleep(d)
        log.Printf(".")
    }
}
