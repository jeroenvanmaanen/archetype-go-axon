package axon_utils

import (
    context "context"
    errors "errors"
    log "log"

    grpc "google.golang.org/grpc"
    proto "github.com/golang/protobuf/proto"
    uuid "github.com/google/uuid"

    axonserver "github.com/jeroenvm/archetype-go-axon/src/pkg/grpc/axonserver"
)

func SendCommand(commandType string, command proto.Message, conn *grpc.ClientConn, clientInfo *axonserver.ClientIdentification) error {
    data, err := proto.Marshal(command)
    if err != nil {
        log.Printf("Server: Error while marshalling command: %v", commandType)
        return errors.New("Marshalling error")
    }
    serializedCommand := axonserver.SerializedObject{
        Type: commandType,
        Data: data,
    }

    return submitCommand(&serializedCommand, conn, clientInfo)
}

func submitCommand(message *axonserver.SerializedObject, conn *grpc.ClientConn, clientInfo *axonserver.ClientIdentification) error {
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