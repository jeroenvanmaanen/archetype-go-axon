package example

import (
    context "context"
    log "log"
    axonserver "github.com/jeroenvm/archetype-nix-go/src/pkg/grpc/axonserver"
    grpc "google.golang.org/grpc"
    uuid "github.com/google/uuid"
)

func SubmitCommand(message string, conn *grpc.ClientConn) {
    log.Printf("Submit command: %v: %v", message, conn)
    client := axonserver.NewCommandServiceClient(conn)
    log.Printf("Client: %v", client)

    uuid := uuid.New()
    command := axonserver.Command {
        MessageIdentifier: uuid.String(),
        Name: "GreetCommand",
        ClientId: "12345",
        ComponentName: "GoClient",
    }
    log.Printf("Command: %v", command)

    response, e := client.Dispatch(context.Background(), &command)
    log.Printf("Response: %v: %v", response, e)
}