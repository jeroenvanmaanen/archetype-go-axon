package configuration_command

import (
    log "log"

    axon_utils "github.com/jeroenvm/archetype-go-axon/src/pkg/axon_utils"
    grpc_example "github.com/jeroenvm/archetype-go-axon/src/pkg/grpc/example"
)

// Redeclare event types, so that they can be extended with event handler methods.
type PropertyChangedEvent struct { grpc_example.PropertyChangedEvent }

// Projection

type Projection struct {
    Configuration map[string]string
}

func RestoreProjection(aggregateIdentifier string, clientConnection *axon_utils.ClientConnection) *Projection {
    projection := &Projection{
        Configuration: make(map[string]string),
    }
    axon_utils.RestoreProjection("Configuration", aggregateIdentifier, projection, clientConnection, prepareUnmarshal)
    return projection
}

func (projection *Projection) Apply(event axon_utils.Event) {
    event.ApplyTo(projection)
}

// Map an event name as stored in AxonServer to an object that has two aspects:
// 1. It is a proto.Message, so it can be unmarshalled from the Axon event
// 2. It is an axon_utils.Event, so it can be applied to a projection
func prepareUnmarshal(payloadType string) (event axon_utils.Event) {
    log.Printf("Configuration Projection: Payload type: %v", payloadType)
    switch payloadType {
        case "PropertyChangedEvent": event = &PropertyChangedEvent{}
        default: event = nil
    }
    return event
}

// Event Handlers

func (event *PropertyChangedEvent) ApplyTo(projectionWrapper interface{}) {
    projection := projectionWrapper.(*Projection)
    projection.Configuration[event.Property.Key] = event.Property.Value
}
