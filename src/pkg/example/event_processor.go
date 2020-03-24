package example

import (
    context "context"
    errors "errors"
    log "log"
    strings "strings"
    time "time"

    hex "encoding/hex"
    json "encoding/json"
    sha256 "crypto/sha256"

    axonserver "github.com/jeroenvm/archetype-nix-go/src/pkg/grpc/axonserver"
    elasticSearch7 "github.com/elastic/go-elasticsearch/v7"
    esapi "github.com/elastic/go-elasticsearch/v7/esapi"
    grpc "google.golang.org/grpc"
    grpcExample "github.com/jeroenvm/archetype-nix-go/src/pkg/grpc/example"
    proto "github.com/golang/protobuf/proto"
    uuid "github.com/google/uuid"
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
    log.Printf("Event processor worker: Get events request: %v", getEventsRequest)


    log.Printf("Event processor worker: Ready to process events")
    greetedEvent := grpcExample.GreetedEvent{}
    for true {
        e = client.Send(&getEventsRequest)
        if e != nil {
            log.Printf("Event processor worker: Error while sending get events request: %v", e)
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

        e = proto.Unmarshal(event.Event.Payload.Data, &greetedEvent)
        log.Printf("Event processor worker: Payload of greeted event: %v", greetedEvent)

        addToIndex(greetedEvent.Message.Message, es7)
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

func WaitForElasticSearch() *elasticSearch7.Client {
    cfg := elasticSearch7.Config{
        Addresses: []string{
            "http://elastic-search:9200",
        },
    }
    d, _ := time.ParseDuration("3s")
    log.Printf(".(es)")
    for true {
        es7, e := elasticSearch7.NewClient(cfg)
        if e == nil {
            if e = getElasticSearchInfo(es7); e == nil {
                return es7
            }
        }
        time.Sleep(d)
        log.Printf(".(es) %v", e)
    }
    return nil
}

func getElasticSearchInfo(es7 *elasticSearch7.Client) error  {
    info, e := es7.Info(es7.Info.WithContext(context.Background()), es7.Info.WithHuman())
    if e != nil {
        log.Printf("Error while requesting Elastic Search info: %v", e)
        return e
    }
    log.Printf("Elastic Search: info: %v", info)

    // Check response status
    if info.IsError() {
        log.Printf("Error: %s", info.String())
        return errors.New("Elastic Search error: " + info.String())
    }

    if info.Body == nil {
        log.Printf("Missing body")
        return errors.New("Elastic Search error: info has no body")
    }

    defer info.Body.Close()
    // Deserialize the response into a map.
    var r map[string]interface{}
    if err := json.NewDecoder(info.Body).Decode(&r); err != nil {
        log.Printf("Error parsing the response body: %s", err)
        return err
    }
    // Print client and server version numbers.
    log.Printf("Client: %s", elasticSearch7.Version)
    log.Printf("Server: %s", r["version"].(map[string]interface{})["number"])
    log.Println(strings.Repeat("~", 37))
    return nil
}

func addToIndex(message string, es7 *elasticSearch7.Client) {
    checksum := sha256.Sum256([]byte(message))
    id := hex.EncodeToString(checksum[:])
    log.Printf("Add to index: Document ID: %v", id)

    // Build the request body.
    var b strings.Builder
    b.WriteString(`{"message" : "`)
    b.WriteString(message)
    b.WriteString(`"}`)

    // Set up the request object.
    req := esapi.IndexRequest{
        Index:      "greetings",
        DocumentID: id,
        Body:       strings.NewReader(b.String()),
        Refresh:    "true",
    }

    // Perform the request with the client.
    res, err := req.Do(context.Background(), es7)
    if err != nil {
        log.Fatalf("Error getting response: %s", err)
    }
    defer res.Body.Close()

    if res.IsError() {
        log.Printf("[%s] Error indexing document ID=%d", res.Status(), id)
    } else {
        // Deserialize the response into a map.
        var r map[string]interface{}
        if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
            log.Printf("Error parsing the response body: %s", err)
        } else {
            // Print the response status and indexed document version.
            log.Printf("[%s] %s; version=%d", res.Status(), r["result"], int(r["_version"].(float64)))
        }
    }
}