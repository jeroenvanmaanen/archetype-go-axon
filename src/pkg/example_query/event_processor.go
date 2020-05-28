package example_query

import (
	log "log"
	strings "strings"

	sha256 "crypto/sha256"
	hex "encoding/hex"

	elasticSearch7 "github.com/elastic/go-elasticsearch/v7"

	axon_utils "github.com/jeroenvanmaanen/dendrite/src/pkg/axon_utils"
	elastic_search_utils "github.com/jeroenvm/archetype-go-axon/src/pkg/elastic_search_utils"
	grpc_example "github.com/jeroenvm/archetype-go-axon/src/pkg/grpc/example"
)

// Redeclare event types, so that they can be extended with event handler methods.
type GreetedEvent struct{ grpc_example.GreetedEvent }

func ProcessEvents(host string, port int) *axon_utils.ClientConnection {
	tokenStore := elastic_search_utils.OpenTokenStore("example-processor")
	projection := tokenStore.ES7
	return axon_utils.ProcessEvents("Example", host, port, tokenStore.ProcessorName, projection, prepareUnmarshal, tokenStore)
}

// Map an event name as stored in AxonServer to an object that has two aspects:
// 1. It is a proto.Message, so it can be unmarshalled from the Axon event
// 2. It is an axon_utils.Event, so it can be applied to a projection
func prepareUnmarshal(payloadType string) (event axon_utils.Event) {
	log.Printf("Example event processor: Payload type: %v", payloadType)
	switch payloadType {
	case "GreetedEvent":
		event = &GreetedEvent{}
	default:
		event = nil
	}
	return event
}

// Event Handlers

func (event *GreetedEvent) ApplyTo(projectionWrapper interface{}) {
	es7 := projectionWrapper.(*elasticSearch7.Client)
	if e := addMessageToIndex(event.Message.Message, es7); e != nil {
		log.Printf("Event processor worker: error while indexing message: %v", e)
		panic("Failure in 'example' event processor")
	}
}

func addMessageToIndex(message string, es7 *elasticSearch7.Client) error {
	checksum := sha256.Sum256([]byte(message))
	id := hex.EncodeToString(checksum[:])

	// Build the request body.
	var b strings.Builder
	b.WriteString(`{"message" : "`)
	b.WriteString(message)
	b.WriteString(`"}`)

	return elastic_search_utils.AddToIndex("greetings", id, b.String(), es7)
}
