package main

import (
    context "context"
    "fmt"
    submathpackage "github.com/jeroenvm/archetype-nix-go/src/pkg/submathpackage"
    axonserver "github.com/jeroenvm/archetype-nix-go/src/pkg/grpc/axonserver"
    grpc "google.golang.org/grpc"
)

var _ axonserver.PlatformOutboundInstruction

func main() {
    fmt.Println(submathpackage.Add(1,2))
    serverAddress := "axon-server:8124"
    conn, e := grpc.Dial(serverAddress, grpc.WithInsecure())

    if e != nil {
        panic(e)
    }
    defer conn.Close()

    client := axonserver.NewPlatformServiceClient(conn)
    var clientInfo axonserver.ClientIdentification
    if response, e := client.GetPlatformServer(context.Background(), &clientInfo); e != nil {
        panic(fmt.Sprintf("Was not able to get Axon platform sever %v", e))
    } else {
        fmt.Println(response)
    }
    fmt.Println("xxx")
}