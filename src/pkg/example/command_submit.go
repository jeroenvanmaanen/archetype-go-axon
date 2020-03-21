package example

import (
    context "context"
    errors "errors"
    log "log"
    axonserver "github.com/jeroenvm/archetype-nix-go/src/pkg/grpc/axonserver"
    grpc "google.golang.org/grpc"
    uuid "github.com/google/uuid"
)

func SubmitCommand(message *axonserver.SerializedObject, conn *grpc.ClientConn, clientInfo *axonserver.ClientIdentification) error {
    log.Printf("Submit command: %v: %v", message.Type, conn)
    client := axonserver.NewCommandServiceClient(conn)
    log.Printf("Submit command: Client: %v", client)

    id := uuid.New()
    command := axonserver.Command {
        MessageIdentifier: id.String(),
        Name: (*message).Type,
        Payload: message,
        ClientId: clientInfo.ClientId,
        ComponentName: clientInfo.ComponentName,
    }
    log.Printf("Submit command: Command: %v", command)

    response, e := client.Dispatch(context.Background(), &command)
    log.Printf("Submit command: Response: %v: %v", response, e)
    if e != nil {
        return e
    } else if response.ErrorMessage != nil {
        return errors.New("Command error: " + response.ErrorCode + ": " + response.ErrorMessage.Message)
    }
    return nil
}