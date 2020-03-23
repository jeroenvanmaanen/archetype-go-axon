package example

import (
    context "context"
    errors "errors"
    fmt "fmt"
    log "log"
    strings "strings"
    time "time"
    json "encoding/json"
    axonserver "github.com/jeroenvm/archetype-nix-go/src/pkg/grpc/axonserver"
    elasticSearch7 "github.com/elastic/go-elasticsearch/v7"
    grpc "google.golang.org/grpc"
    uuid "github.com/google/uuid"
)

func ProcessEvents(host string, port int) *grpc.ClientConn {
    conn, clientInfo, stream := WaitForServer(host, port, "Event processor")
    log.Printf("Connection and client info: %v: %v: %v", conn, clientInfo, stream)

    registerProcessor("example-processor", stream, clientInfo)

    go eventProcessorWorker(stream, conn, clientInfo.ClientId)

    go tryElasticSearch()

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

func tryElasticSearch() {
    log.Printf("Try Elastic Search")
    es7 := WaitForElasticSearch();
    log.Printf("Elastic Search client: %v", es7)
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