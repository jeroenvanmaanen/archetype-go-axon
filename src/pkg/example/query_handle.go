package example

import (
    context "context"
    fmt "fmt"
    log "log"
    axonserver "github.com/jeroenvm/archetype-nix-go/src/pkg/grpc/axonserver"
    grpcExample "github.com/jeroenvm/archetype-nix-go/src/pkg/grpc/example"
    grpc "google.golang.org/grpc"
    proto "github.com/golang/protobuf/proto"
    uuid "github.com/google/uuid"
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
                log.Printf("Received SearchQuery")
                handleSearchQuery(query, stream, conn)
            } else {
                log.Printf("Received unknown query: %v", queryName)
            }
        }
    }
}

func handleSearchQuery(query *axonserver.QueryRequest, stream axonserver.QueryService_OpenStreamClient, conn *grpc.ClientConn) {
    deserializedQuery := grpcExample.SearchQuery{}
    e := proto.Unmarshal(query.Payload.Data, &deserializedQuery)
    if (e != nil) {
        log.Printf("Could not unmarshall SearchQuery")
    }

    response := grpcExample.Greeting{
        Message: query.Query,
    }
    queryRespond(&response, stream, query.MessageIdentifier)

    response.Message = "Over and out."
    queryRespond(&response, stream, query.MessageIdentifier)

    queryComplete(stream, query.MessageIdentifier)
}

func queryRespond(response *grpcExample.Greeting, stream axonserver.QueryService_OpenStreamClient, requestId string) {
    responseData, e := proto.Marshal(response)
    if e != nil {
        log.Printf("Server: Error while marshalling query response: %v", e)
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
