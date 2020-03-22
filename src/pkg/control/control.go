package control

import (
    "log"
    axonserver "github.com/jeroenvm/archetype-nix-go/src/pkg/grpc/axonserver"
)

func Listen(streamClient *axonserver.PlatformService_OpenStreamClient) {
    for {
        log.Printf("Listen: Waiting for next message...")
        message, e := (*streamClient).Recv()
        if e != nil {
            log.Printf("Listen: Error while receiving message: %v", e)
            panic("Listen: Panic!")
        }
        log.Printf("Listen: Received message: %v", message)
    }
}
