package axon_utils

import (
    fmt "fmt"
    log "log"
    net "net"

    grpc "google.golang.org/grpc"
    reflection "google.golang.org/grpc/reflection"
)

func Serve(clientConnection *ClientConnection, registerWithServer func(*grpc.Server, *ClientConnection)) {
    port := 8181
    lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
    if err != nil {
        log.Fatalf("Server: Failed to listen: %v", err)
        panic("Server: Panic!")
    }
    log.Printf("Server: Listening on port: %d", port)
    grpcServer := grpc.NewServer()
    registerWithServer(grpcServer, clientConnection)
    reflection.Register(grpcServer)
    // ... // determine whether to use TLS
    grpcServer.Serve(lis)
}
