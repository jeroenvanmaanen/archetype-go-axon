package example

import (
    fmt "fmt"
    log "log"
    axonserver "github.com/jeroenvm/archetype-nix-go/src/pkg/grpc/axonserver"
    grpc "google.golang.org/grpc"
    uuid "github.com/google/uuid"
)

func ProcessEvents(host string, port int) *grpc.ClientConn {
    conn, clientInfo, stream := WaitForServer(host, port, "Event processor")
    log.Printf("Connection and client info: %v: %v: %v", conn, clientInfo, stream)

    registerProcessor("example-processor", stream, clientInfo)

    go eventProcessorWorker(stream, conn, clientInfo.ClientId)

    return conn;
}

func registerProcessor(processorName string, stream *axonserver.PlatformService_OpenStreamClient, clientInfo *axonserver.ClientIdentification) {
    processorInfo := axonserver.EventProcessorInfo {
        ProcessorName: processorName,
        Mode: "Tracking",
        ActiveThreads: 0,
        Running: true,
        Error: false,
        SegmentStatus: nil,
        AvailableThreads: 1,
    }
    log.Printf("Event processor: event processor info: %v", processorInfo)
    subscriptionRequest := axonserver.PlatformInboundInstruction_EventProcessorInfo {
        EventProcessorInfo: &processorInfo,
    }

    id := uuid.New()
    inbound := axonserver.PlatformInboundInstruction {
        Request: &subscriptionRequest,
        InstructionId: id.String(),
    }
    log.Printf("Event processor: event processor info: instruction ID: %v", inbound.InstructionId)

    e := (*stream).Send(&inbound)
    if e != nil {
        panic(fmt.Sprintf("Event processor: Error sending subscription", e))
    }
}

func eventProcessorWorker(stream *axonserver.PlatformService_OpenStreamClient, conn *grpc.ClientConn, clientId string) {
    for true {
        // addPermits(1, stream, clientId)

        e := eventProcessorReceivePlatformInstruction(stream)
        if e != nil {
            break
        }
    }
}

func eventProcessorReceivePlatformInstruction(stream *axonserver.PlatformService_OpenStreamClient) error {
    log.Printf("Event processor receive platform instruction: Waiting for outbound")
    outbound, e := (*stream).Recv()
    if (e != nil) {
      log.Printf("Event processor receive platform instruction: Error on receive, %v", e)
      return e
    }
    log.Printf("Event processor receive platform instruction: Outbound: %v", outbound)
    return nil
}