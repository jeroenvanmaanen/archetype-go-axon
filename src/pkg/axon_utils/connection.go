package axon_utils

import (
    context "context"
    fmt "fmt"
    log "log"
    time "time"
    axonserver "github.com/jeroenvm/archetype-go-axon/src/pkg/grpc/axonserver"
    grpc "google.golang.org/grpc"
    grpcKeepAlive "google.golang.org/grpc/keepalive"
    uuid "github.com/google/uuid"
)

func WaitForServer(host string, port int, qualifier string) (*grpc.ClientConn, *axonserver.ClientIdentification, *axonserver.PlatformService_OpenStreamClient) {
    id := uuid.New()
    clientInfo := axonserver.ClientIdentification {
        ClientId: id.String(),
        ComponentName: "GoClient " + qualifier,
        Version: "0.0.1",
    }

    serverAddress := fmt.Sprintf("%s:%d", host, port)
    log.Printf("Connection: Client identification: %v", clientInfo)
    d, _ := time.ParseDuration("3s")
    first := true
    for {
        if first {
            first = false
        } else {
            time.Sleep(d)
            log.Printf(".")
        }
        keepAlive := grpcKeepAlive.ClientParameters{
            Time:                20 * time.Second,
            Timeout:             10 * time.Second,
            PermitWithoutStream: true,
        }
        conn, e := grpc.Dial(serverAddress, grpc.WithInsecure(), grpc.WithKeepaliveParams(keepAlive), )
        if e != nil {
            continue
        }
        // Get platform server
        client := axonserver.NewPlatformServiceClient(conn)
        response, e := client.GetPlatformServer(context.Background(), &clientInfo)
        if e != nil {
            continue
        }
        log.Printf("Connection: Connected: %v: %v", response.SameConnection, response)
        if !response.SameConnection {
            panic(fmt.Sprintf("Connection: Need to setup a new connection %v", e))
        }
        streamClient, e := registerClient(&clientInfo, &client)
        if e != nil {
            continue
        }
        return conn, &clientInfo, streamClient
    }
}

func registerClient(clientInfo *axonserver.ClientIdentification, client *axonserver.PlatformServiceClient) (*axonserver.PlatformService_OpenStreamClient, error) {

    // Open stream
    streamClient, e := (*client).OpenStream(context.Background())
    if e != nil {
        panic(fmt.Sprintf("Connection: Could not open stream %v", e))
    }

    // Send client info
    var instruction axonserver.PlatformInboundInstruction
    registrationRequest := axonserver.PlatformInboundInstruction_Register{
        Register: clientInfo,
    }
    id := uuid.New()
    instruction.Request = &registrationRequest
    instruction.InstructionId = id.String()
    if e = streamClient.Send(&instruction); e != nil {
        panic(fmt.Sprintf("Connection: Error sending clientInfo %v", e))
    }

    log.Printf("Connection receive platform instruction: Waiting for outbound")
    outbound, e := streamClient.Recv()
    if (e != nil) {
      log.Printf("Connection receive platform instruction: Error on receive, %v", e)
      return nil, e
    }
    log.Printf("Connection receive platform instruction: Outbound: %v", outbound)

    return &streamClient, nil
}
