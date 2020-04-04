package example

import (
    context "context"
    log "log"
    strconv "strconv"
    strings "strings"

    hex "encoding/hex"
    ioutil "io/ioutil"
    json "encoding/json"
    sha256 "crypto/sha256"

    elasticSearch7 "github.com/elastic/go-elasticsearch/v7"
    grpc "google.golang.org/grpc"
    proto "github.com/golang/protobuf/proto"
    uuid "github.com/google/uuid"

    axonserver "github.com/jeroenvm/archetype-go-axon/src/pkg/grpc/axonserver"
    grpcExample "github.com/jeroenvm/archetype-go-axon/src/pkg/grpc/example"
)

func ProcessEvents(host string, port int) *grpc.ClientConn {
    conn, clientInfo, stream := WaitForServer(host, port, "Event processor")
    log.Printf("Connection and client info: %v: %v: %v", conn, clientInfo, stream)

    processorName := "example-processor"

    if e := registerProcessor(processorName, stream, clientInfo); e != nil {
        return conn
    }

    go eventProcessorWorker(stream, conn, clientInfo, processorName)

    return conn;
}

func registerProcessor(processorName string, stream *axonserver.PlatformService_OpenStreamClient, clientInfo *axonserver.ClientIdentification) error {
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
        log.Printf("Event processor: Error sending registration", e)
        return e
    }

    e = eventProcessorReceivePlatformInstruction(stream)
    if e != nil {
        log.Printf("Event processor: Error while waiting for acknowledgement of registration")
        return e
    }
    return nil
}

func eventProcessorWorker(stream *axonserver.PlatformService_OpenStreamClient, conn *grpc.ClientConn, clientInfo *axonserver.ClientIdentification, processorName string) {
    es7 := WaitForElasticSearch();
    log.Printf("Elastic Search client: %v", es7)

    token, _ := readToken(processorName, es7)

    eventStoreClient := axonserver.NewEventStoreClient(conn)
    log.Printf("Event processor worker: Event store client: %v", eventStoreClient)
    client, e := eventStoreClient.ListEvents(context.Background())
    if e != nil {
        log.Printf("Event processor worker: Error while opening ListEvents stream: %v", e)
        return
    }
    log.Printf("Event processor worker: List events client: %v", client)

    getEventsRequest := axonserver.GetEventsRequest{
        NumberOfPermits: 1,
        ClientId: clientInfo.ClientId,
        ComponentName: clientInfo.ComponentName,
        Processor: processorName,
    }
    if token != nil {
        getEventsRequest.TrackingToken = *token + 1
    }
    log.Printf("Event processor worker: Get events request: %v", getEventsRequest)


    log.Printf("Event processor worker: Ready to process events")
    greetedEvent := grpcExample.GreetedEvent{}
    defer func() {
        log.Printf("Event processor worker stopped")
    }()
    first := true
    for true {
        if first {
            first = false
        } else {
            var b strings.Builder
            b.WriteString(`{"token" : "`)
            b.WriteString(strconv.FormatInt(getEventsRequest.TrackingToken, 16))
            b.WriteString(`"}`)
            if e = addToIndex("tracking-token", processorName, b.String(), es7); e != nil {
                log.Printf("Event processor worker: Error while storing tracking token: %v", e)
                return
            }
        }
        e = client.Send(&getEventsRequest)
        if e != nil {
            log.Printf("Event processor worker: Error while sending get events request: %v", e)
            return
        }

        event, e := client.Recv()
        if e != nil {
            log.Printf("Event processor worker: Error while receiving next event: %v", e)
            return
        }
        log.Printf("Event processor worker: Next event: %v", event)
        getEventsRequest.TrackingToken = event.Token

        if event.Event == nil || event.Event.Payload == nil || event.Event.Payload.Type != "GreetedEvent" {
            continue
        }

        if e = proto.Unmarshal(event.Event.Payload.Data, &greetedEvent); e != nil {
            log.Printf("Event processor worker: Unmarshalling of GreetedEvent failed: %v", e)
            return
        }
        log.Printf("Event processor worker: Payload of greeted event: %v", greetedEvent)

        if e = addMessageToIndex(greetedEvent.Message.Message, es7); e != nil {
            log.Printf("Event processor worker: error while indexing message: %v", e)
            return
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

func addMessageToIndex(message string, es7 *elasticSearch7.Client) error {
    checksum := sha256.Sum256([]byte(message))
    id := hex.EncodeToString(checksum[:])

    // Build the request body.
    var b strings.Builder
    b.WriteString(`{"message" : "`)
    b.WriteString(message)
    b.WriteString(`"}`)

    return addToIndex("greetings", id, b.String(), es7)
}

func readToken(processorName string, es7 *elasticSearch7.Client) (*int64, error) {
    response, e := es7.Get("tracking-token", processorName)
    if e != nil {
        log.Printf("Elastic search: Error while reading token: %v", e)
        return nil, e
    }
    log.Printf("Elastic search: token document: %v", response)

    if response.StatusCode == 404 {
        return nil, nil
    }

    responseJson, e := ioutil.ReadAll(response.Body)
    if e != nil {
        log.Printf("Elastic search: Error while reading response body: %v", e)
        return nil, e
    }

    jsonMap := make(map[string](interface{}))
    e = json.Unmarshal(responseJson, &jsonMap)
    if e != nil {
        log.Printf("Elastic search: Error while unmarshalling JSON, %v", e)
        return nil, e
    }

    hexToken := jsonMap["_source"].(map[string](interface{}))["token"].(string)
    token, e := strconv.ParseInt(hexToken, 16, 64)
    if e != nil {
        log.Printf("Elastic search: Error while parsing hex token, %v", e)
        return nil, e
    }
    log.Printf("Elastic search: token: %v", token)

    return &token, nil
}