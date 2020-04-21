package example

import (
    context "context"
    fmt "fmt"
    log "log"

    elasticSearch7 "github.com/elastic/go-elasticsearch/v7"
    grpc "google.golang.org/grpc"
    proto "github.com/golang/protobuf/proto"
    uuid "github.com/google/uuid"

    axonserver "github.com/jeroenvm/archetype-go-axon/src/pkg/grpc/axonserver"
    grpcExample "github.com/jeroenvm/archetype-go-axon/src/pkg/grpc/example"
)

func HandleQueries(host string, port int) (conn *grpc.ClientConn) {
    conn, clientInfo, _ := WaitForServer(host, port, "Query Handler")

    log.Printf("Query handler: Connection: %v", conn)
    client := axonserver.NewQueryServiceClient(conn)
    log.Printf("Query handler: Client: %v", client)

    stream, e := client.OpenStream(context.Background())
    log.Printf("Query handler: Stream: %v: %v", stream, e)

    querySubscribe("SearchQuery", stream, clientInfo)

    go queryWorker(stream, conn, clientInfo.ClientId)

    return conn;
}

func querySubscribe(queryName string, stream axonserver.QueryService_OpenStreamClient, clientInfo *axonserver.ClientIdentification) {
    id := uuid.New()
    subscription := axonserver.QuerySubscription {
        MessageId: id.String(),
        Query: queryName,
        ClientId: clientInfo.ClientId,
        ComponentName: clientInfo.ComponentName,
    }
    log.Printf("Query handler: Subscription: %v", subscription)
    subscriptionRequest := axonserver.QueryProviderOutbound_Subscribe {
        Subscribe: &subscription,
    }

    outbound := axonserver.QueryProviderOutbound {
        Request: &subscriptionRequest,
    }

    e := stream.Send(&outbound)
    if e != nil {
        panic(fmt.Sprintf("Query handler: Error sending subscription", e))
    }
}

func queryWorker(stream axonserver.QueryService_OpenStreamClient, conn *grpc.ClientConn, clientId string) {
    es7 := WaitForElasticSearch();
    log.Printf("Query handler: Elastic Search client: %v", es7)

    for true {
        queryAddPermits(1, stream, clientId)

        log.Printf("Query handler: Waiting for query")
        inbound, e := stream.Recv()
        if (e != nil) {
          log.Printf("Query handler: Error on receive: %v", e)
          break
        }
        log.Printf("Query handler: Inbound: %v", inbound)
        query := inbound.GetQuery()
        if (query != nil) {
            queryName := query.Query
            if (queryName == "SearchQuery") {
                log.Printf("Query handler: Received SearchQuery")
                handleSearchQuery(query, stream, es7)
            } else {
                log.Printf("Query handler: Received unknown query: %v", queryName)
            }
        }
    }
}

func handleSearchQuery(axonQuery *axonserver.QueryRequest, stream axonserver.QueryService_OpenStreamClient, es7 *elasticSearch7.Client) {
    defer queryComplete(stream, axonQuery.MessageIdentifier)
    query := grpcExample.SearchQuery{}
    e := proto.Unmarshal(axonQuery.Payload.Data, &query)
    if (e != nil) {
        log.Printf("Query handler: Could not unmarshal SearchQuery")
    }

    reply := grpcExample.Greeting{
        Message: "Query: '" + query.Query + "'",
    }
    queryRespond(&reply, stream, axonQuery.MessageIdentifier)

    response, e := es7.Search(es7.Search.WithQuery(query.Query))
    if e != nil {
        return
    }
    result, e := unwrapElasticSearchResponse(response)
    log.Printf("Query handler: result: %v", result)
    hitsWrapper := result["hits"].(map[string](interface{}))
    log.Printf("Query handler: hits: %v", hitsWrapper)
    hits := hitsWrapper["hits"].([]interface{})
    log.Printf("Query handler: hits: %v", hits)

    for _, hit := range hits {
        log.Printf("Query handler: hit: %v", hit)
        source := hit.(map[string](interface{}))["_source"]
        log.Printf("Query handler: source: %v", source)
        message := source.(map[string](interface{}))["message"]
        if message == nil {
            continue
        }
        log.Printf("Query handler: message: %v", message.(string))

        reply.Message = message.(string)
        queryRespond(&reply, stream, axonQuery.MessageIdentifier)
    }

    reply.Message = "Over and out."
    queryRespond(&reply, stream, axonQuery.MessageIdentifier)
}

func queryRespond(response *grpcExample.Greeting, stream axonserver.QueryService_OpenStreamClient, requestId string) {
    responseData, e := proto.Marshal(response)
    if e != nil {
        log.Printf("Query handler: Error while marshalling query response: %v", e)
        return
    }

    serializedResponse := axonserver.SerializedObject{
        Type: "Greeting",
        Data: responseData,
    }

    id := uuid.New()
    queryResponse := axonserver.QueryResponse {
        MessageIdentifier: id.String(),
        RequestIdentifier: requestId,
        Payload: &serializedResponse,
    }
    log.Printf("Query handler: Query response: %v", queryResponse)
    queryResponseRequest := axonserver.QueryProviderOutbound_QueryResponse {
        QueryResponse: &queryResponse,
    }

    outbound := axonserver.QueryProviderOutbound {
        Request: &queryResponseRequest,
    }

    e = stream.Send(&outbound)
    if e != nil {
        panic(fmt.Sprintf("Query handler: Error sending command response", e))
    }
}

func queryComplete(stream axonserver.QueryService_OpenStreamClient, requestId string) {
    id := uuid.New()
    queryComplete := axonserver.QueryComplete {
        MessageId: id.String(),
        RequestId: requestId,
    }
    log.Printf("Query handler: Query complete: %v", queryComplete)
    queryCompleteRequest := axonserver.QueryProviderOutbound_QueryComplete {
        QueryComplete: &queryComplete,
    }

    outbound := axonserver.QueryProviderOutbound {
        Request: &queryCompleteRequest,
    }

    e := stream.Send(&outbound)
    if e != nil {
        panic(fmt.Sprintf("Query handler: Error sending command response", e))
    }
}

func queryAddPermits(amount int64, stream axonserver.QueryService_OpenStreamClient, clientId string) {
    flowControl := axonserver.FlowControl {
        ClientId: clientId,
        Permits: amount,
    }
    log.Printf("Query handler: Flow control: %v", flowControl)
    flowControlRequest := axonserver.QueryProviderOutbound_FlowControl {
        FlowControl: &flowControl,
    }

    outbound := axonserver.QueryProviderOutbound {
        Request: &flowControlRequest,
    }

    e := stream.Send(&outbound)
    if e != nil {
        panic(fmt.Sprintf("Query handler: Error sending flow control", e))
    }
}
