package example

import (
    context "context"
    log "log"
    axonserver "github.com/jeroenvm/archetype-nix-go/src/pkg/grpc/axonserver"
    grpc "google.golang.org/grpc"
//    uuid "github.com/google/uuid"
)

func HandleCommands(conn *grpc.ClientConn) {
    log.Printf("Connection: %v", conn)
    client := axonserver.NewCommandServiceClient(conn)
    log.Printf("Client: %v", client)

    stream, e := client.OpenStream(context.Background())
    log.Printf("Stream: %v: %v", stream, e)

    subscription := axonserver.CommandSubscription {
        MessageId: "54321", // uuid.String(),
        Command: "GreetCommand",
        ClientId: "12345",
        ComponentName: "GoClient",
    }
    log.Printf("Subscription: %v", subscription)
    subscriptionRequest := axonserver.CommandProviderOutbound_Subscribe {
        Subscribe: &subscription,
    }

    outbound := axonserver.CommandProviderOutbound {
        Request: &subscriptionRequest,
    }

    stream.Send(&outbound)
    for true {
        inbound, _ := stream.Recv()
        log.Printf("Inbound: %v", inbound)
// TODO: Ack
//         command = inbound.GetCommand()
//         if (command != null) {
//
//         }
    }
}