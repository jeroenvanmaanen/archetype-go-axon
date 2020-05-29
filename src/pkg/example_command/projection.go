package example_command

import (
	log "log"

	grpc_example "github.com/jeroenvanmaanen/archetype-go-axon/src/pkg/grpc/example"
	axon_utils "github.com/jeroenvanmaanen/dendrite/src/pkg/axon_utils"
)

// Redeclare event types, so that they can be extended with event handler methods.
type GreetedEvent struct{ grpc_example.GreetedEvent }
type StartedRecordingEvent struct {
	grpc_example.StartedRecordingEvent
}
type StoppedRecordingEvent struct {
	grpc_example.StoppedRecordingEvent
}

// Projection

type Projection struct {
	Recording      bool
	AggregateState axon_utils.AggregateState
}

func (projection *Projection) GetAggregateState() axon_utils.AggregateState {
	return projection.AggregateState
}

func RestoreProjection(aggregateIdentifier string, clientConnection *axon_utils.ClientConnection) *Projection {
	projection := axon_utils.RestoreProjection("Example", aggregateIdentifier, createInitialProjection, clientConnection, prepareUnmarshal).(*Projection)
	log.Printf("Example is recording: %v", projection.Recording)
	return projection
}

func createInitialProjection() interface{} {
	return &Projection{
		Recording:      true,
		AggregateState: axon_utils.NewAggregateState(),
	}
}

func (projection *Projection) Apply(event axon_utils.Event) {
	event.ApplyTo(projection)
}

// Map an event name as stored in AxonServer to an object that has two aspects:
// 1. It is a proto.Message, so it can be unmarshalled from the Axon event
// 2. It is an axon_utils.Event, so it can be applied to a projection
func prepareUnmarshal(payloadType string) (event axon_utils.Event) {
	log.Printf("Example Projection: Payload type: %v", payloadType)
	switch payloadType {
	case "StartedRecordingEvent":
		event = &StartedRecordingEvent{}
	case "StoppedRecordingEvent":
		event = &StoppedRecordingEvent{}
	default:
		event = nil
	}
	return event
}

// Event Handlers

func (event *GreetedEvent) ApplyTo(_ interface{}) {}

func (event *StartedRecordingEvent) ApplyTo(projectionWrapper interface{}) {
	projection := projectionWrapper.(*Projection)
	projection.Recording = true
}

func (event *StoppedRecordingEvent) ApplyTo(projectionWrapper interface{}) {
	projection := projectionWrapper.(*Projection)
	projection.Recording = false
}
