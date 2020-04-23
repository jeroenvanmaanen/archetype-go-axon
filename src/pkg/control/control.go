package control

import (
    "log"
    axon_server "github.com/jeroenvm/archetype-go-axon/src/pkg/grpc/axon_server"
)

func Listen(streamClient *axon_server.PlatformService_OpenStreamClient) {
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
