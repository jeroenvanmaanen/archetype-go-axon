package trusted

import (
    log "log"

    axon_utils "github.com/jeroenvm/archetype-go-axon/src/pkg/axon_utils"
    grpc_example "github.com/jeroenvm/archetype-go-axon/src/pkg/grpc/example"
)

// Redeclare event types, so that they can be extended with event handler methods.
type TrustedKeyAddedSourceEvent   struct { grpc_example.TrustedKeyAddedEvent   }
type TrustedKeyRemovedSourceEvent struct { grpc_example.TrustedKeyRemovedEvent }
type KeyManagerAddedSourceEvent   struct { grpc_example.KeyManagerAddedEvent   }
type KeyManagerRemovedSourceEvent struct { grpc_example.KeyManagerRemovedEvent }

// Projection

type Projection struct {
    TrustedKeys map[string]string
    KeyManagers map[string]string
}

func RestoreProjection(aggregateIdentifier string, clientConnection *axon_utils.ClientConnection) *Projection {
    projection := &Projection{
        TrustedKeys: make(map[string]string),
        KeyManagers: make(map[string]string),
    }
    axon_utils.RestoreProjection("Trusted Keys", aggregateIdentifier, projection, clientConnection, prepareUnmarshal)
    return projection
}

func (projection *Projection) Apply(event axon_utils.SourceEvent) {
    event.ApplyTo(projection)
}

// Map an event name as stored in AxonServer to an object that has two aspects:
// 1. It is a proto.Message, so it can be unmarshalled from the Axon event
// 2. It is an axon_utils.SourceEvent, so it can be applied to a projection
func prepareUnmarshal(payloadType string) (sourceEvent axon_utils.SourceEvent) {
    log.Printf("Configuration Projection: Payload type: %v", payloadType)
    switch payloadType {
        case "TrustedKeyAddedEvent":   sourceEvent = &TrustedKeyAddedSourceEvent{}
        case "TrustedKeyRemovedEvent": sourceEvent = &TrustedKeyRemovedSourceEvent{}
        case "KeyManagerAddedEvent":   sourceEvent = &KeyManagerAddedSourceEvent{}
        case "KeyManagerRemovedEvent": sourceEvent = &KeyManagerRemovedSourceEvent{}
        default: sourceEvent = nil
    }
    return sourceEvent
}

// Event Handlers

func (sourceEvent *TrustedKeyAddedSourceEvent) ApplyTo(projectionWrapper interface{}) {
    projection := projectionWrapper.(*Projection)
    projection.TrustedKeys[sourceEvent.PublicKey.Name] = sourceEvent.PublicKey.PublicKey
}

func (sourceEvent *TrustedKeyRemovedSourceEvent) ApplyTo(projectionWrapper interface{}) {
    projection := projectionWrapper.(*Projection)
    projection.TrustedKeys[sourceEvent.Name] = ""
}

func (sourceEvent *KeyManagerAddedSourceEvent) ApplyTo(projectionWrapper interface{}) {
    projection := projectionWrapper.(*Projection)
    projection.KeyManagers[sourceEvent.PublicKey.Name] = sourceEvent.PublicKey.PublicKey
}

func (sourceEvent *KeyManagerRemovedSourceEvent) ApplyTo(projectionWrapper interface{}) {
    projection := projectionWrapper.(*Projection)
    projection.KeyManagers[sourceEvent.Name] = ""
}
