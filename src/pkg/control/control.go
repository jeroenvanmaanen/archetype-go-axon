package control

import (
    "fmt"
    axonserver "github.com/jeroenvm/archetype-nix-go/src/pkg/grpc/axonserver"
)

func Listen(streamClient *axonserver.PlatformService_OpenStreamClient) {
    for {
        fmt.Println("Waiting for next message...")
        message, e := (*streamClient).Recv()
        if e != nil {
            panic(fmt.Sprintf("Error while receiving message %v", e))
        }
        fmt.Println(fmt.Sprintf("Received message: %v", message))
    }
}
