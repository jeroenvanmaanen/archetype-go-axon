package example

import (
    context "context"
    log "log"
    axonserver "github.com/jeroenvm/archetype-nix-go/src/pkg/grpc/axonserver"
    grpc "google.golang.org/grpc"
    uuid "github.com/google/uuid"
)

func SubmitCommand(message string, conn *grpc.ClientConn, clientInfo *axonserver.ClientIdentification) {
    log.Printf("Submit command: %v: %v", message, conn)
    client := axonserver.NewCommandServiceClient(conn)
    log.Printf("Submit command: Client: %v", client)

    id := uuid.New()
    command := axonserver.Command {
        MessageIdentifier: id.String(),
        Name: "GreetCommand",
        ClientId: clientInfo.ClientId,
        ComponentName: clientInfo.ComponentName,
    }
    log.Printf("Submit command: Command: %v", command)

    response, e := client.Dispatch(context.Background(), &command)
    log.Printf("Submit command: Response: %v: %v", response, e)
}