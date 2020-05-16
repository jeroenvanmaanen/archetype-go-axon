package example_query

import (
	context "context"
	fmt "fmt"
	log "log"

	elasticSearch7 "github.com/elastic/go-elasticsearch/v7"
	proto "github.com/golang/protobuf/proto"
	uuid "github.com/google/uuid"

	axon_utils "github.com/jeroenvm/archetype-go-axon/src/pkg/axon_utils"
	elastic_search_utils "github.com/jeroenvm/archetype-go-axon/src/pkg/elastic_search_utils"
	axon_server "github.com/jeroenvm/archetype-go-axon/src/pkg/grpc/axon_server"
	grpc_example "github.com/jeroenvm/archetype-go-axon/src/pkg/grpc/example"
)

func HandleQueries(host string, port int) (clientConnection *axon_utils.ClientConnection) {
	clientConnection, _ = axon_utils.WaitForServer(host, port, "Query Handler")
	conn := clientConnection.Connection

	log.Printf("Query handler: Connection: %v", conn)
	client := axon_server.NewQueryServiceClient(conn)
	log.Printf("Query handler: Client: %v", client)

	stream, e := client.OpenStream(context.Background())
	log.Printf("Query handler: Stream: %v: %v", stream, e)

	querySubscribe("SearchQuery", stream, clientConnection)

	go queryWorker(stream, clientConnection)

	return clientConnection
}

func querySubscribe(queryName string, stream axon_server.QueryService_OpenStreamClient, clientConnection *axon_utils.ClientConnection) {
	clientInfo := clientConnection.ClientInfo
	id := uuid.New()
	subscription := axon_server.QuerySubscription{
		MessageId:     id.String(),
		Query:         queryName,
		ClientId:      clientInfo.ClientId,
		ComponentName: clientInfo.ComponentName,
	}
	log.Printf("Query handler: Subscription: %v", subscription)
	subscriptionRequest := axon_server.QueryProviderOutbound_Subscribe{
		Subscribe: &subscription,
	}

	outbound := axon_server.QueryProviderOutbound{
		Request: &subscriptionRequest,
	}

	e := stream.Send(&outbound)
	if e != nil {
		panic(fmt.Sprintf("Query handler: Error sending subscription: %v", e))
	}
}

func queryWorker(stream axon_server.QueryService_OpenStreamClient, clientConnection *axon_utils.ClientConnection) {
	clientId := clientConnection.ClientInfo.ClientId
	es7 := elastic_search_utils.WaitForElasticSearch()
	log.Printf("Query handler: Elastic Search client: %v", es7)

	for true {
		queryAddPermits(1, stream, clientId)

		log.Printf("Query handler: Waiting for query")
		inbound, e := stream.Recv()
		if e != nil {
			log.Printf("Query handler: Error on receive: %v", e)
			break
		}
		log.Printf("Query handler: Inbound: %v", inbound)
		query := inbound.GetQuery()
		if query != nil {
			queryName := query.Query
			if queryName == "SearchQuery" {
				log.Printf("Query handler: Received SearchQuery")
				handleSearchQuery(query, stream, es7)
			} else {
				log.Printf("Query handler: Received unknown query: %v", queryName)
			}
		}
	}
}

func handleSearchQuery(axonQuery *axon_server.QueryRequest, stream axon_server.QueryService_OpenStreamClient, es7 *elasticSearch7.Client) {
	defer queryComplete(stream, axonQuery.MessageIdentifier)
	query := grpc_example.SearchQuery{}
	e := proto.Unmarshal(axonQuery.Payload.Data, &query)
	if e != nil {
		log.Printf("Query handler: Could not unmarshal SearchQuery")
	}

	reply := grpc_example.Greeting{
		Message: "Query: '" + query.Query + "'",
	}
	queryRespond(&reply, stream, axonQuery.MessageIdentifier)

	response, e := es7.Search(es7.Search.WithQuery(query.Query))
	if e != nil {
		return
	}
	result, e := elastic_search_utils.UnwrapElasticSearchResponse(response)
	log.Printf("Query handler: result: %v", result)
	hitsWrapper := result["hits"].(map[string]interface{})
	log.Printf("Query handler: hits: %v", hitsWrapper)
	hits := hitsWrapper["hits"].([]interface{})
	log.Printf("Query handler: hits: %v", hits)

	for _, hit := range hits {
		log.Printf("Query handler: hit: %v", hit)
		source := hit.(map[string]interface{})["_source"]
		log.Printf("Query handler: source: %v", source)
		message := source.(map[string]interface{})["message"]
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

func queryRespond(response *grpc_example.Greeting, stream axon_server.QueryService_OpenStreamClient, requestId string) {
	responseData, e := proto.Marshal(response)
	if e != nil {
		log.Printf("Query handler: Error while marshalling query response: %v", e)
		return
	}

	serializedResponse := axon_server.SerializedObject{
		Type: "Greeting",
		Data: responseData,
	}

	id := uuid.New()
	queryResponse := axon_server.QueryResponse{
		MessageIdentifier: id.String(),
		RequestIdentifier: requestId,
		Payload:           &serializedResponse,
	}
	log.Printf("Query handler: Query response: %v", queryResponse)
	queryResponseRequest := axon_server.QueryProviderOutbound_QueryResponse{
		QueryResponse: &queryResponse,
	}

	outbound := axon_server.QueryProviderOutbound{
		Request: &queryResponseRequest,
	}

	e = stream.Send(&outbound)
	if e != nil {
		panic(fmt.Sprintf("Query handler: Error sending command response: %v", e))
	}
}

func queryComplete(stream axon_server.QueryService_OpenStreamClient, requestId string) {
	id := uuid.New()
	queryComplete := axon_server.QueryComplete{
		MessageId: id.String(),
		RequestId: requestId,
	}
	log.Printf("Query handler: Query complete: %v", queryComplete)
	queryCompleteRequest := axon_server.QueryProviderOutbound_QueryComplete{
		QueryComplete: &queryComplete,
	}

	outbound := axon_server.QueryProviderOutbound{
		Request: &queryCompleteRequest,
	}

	e := stream.Send(&outbound)
	if e != nil {
		panic(fmt.Sprintf("Query handler: Error sending command response: %v", e))
	}
}

func queryAddPermits(amount int64, stream axon_server.QueryService_OpenStreamClient, clientId string) {
	flowControl := axon_server.FlowControl{
		ClientId: clientId,
		Permits:  amount,
	}
	log.Printf("Query handler: Flow control: %v", flowControl)
	flowControlRequest := axon_server.QueryProviderOutbound_FlowControl{
		FlowControl: &flowControl,
	}

	outbound := axon_server.QueryProviderOutbound{
		Request: &flowControlRequest,
	}

	e := stream.Send(&outbound)
	if e != nil {
		panic(fmt.Sprintf("Query handler: Error sending flow control: %v", e))
	}
}
