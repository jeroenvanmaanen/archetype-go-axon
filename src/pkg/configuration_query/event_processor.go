package configuration_query

import (
	log "log"

	authentication "github.com/jeroenvanmaanen/archetype-go-axon/src/pkg/authentication"
	grpc_example "github.com/jeroenvanmaanen/archetype-go-axon/src/pkg/grpc/example"
	trusted "github.com/jeroenvanmaanen/archetype-go-axon/src/pkg/trusted"
	axon_utils "github.com/jeroenvanmaanen/dendrite/src/pkg/axon_utils"
)

// Redeclare event types, so that they can be extended with event handler methods.
type PropertyChangedEvent struct {
	grpc_example.PropertyChangedEvent
}
type TrustedKeyAddedEvent struct {
	grpc_example.TrustedKeyAddedEvent
}
type TrustedKeyRemovedEvent struct {
	grpc_example.TrustedKeyRemovedEvent
}
type KeyManagerAddedEvent struct {
	grpc_example.KeyManagerAddedEvent
}
type KeyManagerRemovedEvent struct {
	grpc_example.KeyManagerRemovedEvent
}
type CredentialsAddedEvent struct {
	grpc_example.CredentialsAddedEvent
}
type CredentialsRemovedEvent struct {
	grpc_example.CredentialsRemovedEvent
}

// Projection

type Projection struct {
	Configuration map[string]string
}

var projection Projection

func ProcessEvents(host string, port int) *axon_utils.ClientConnection {
	projection = Projection{
		Configuration: make(map[string]string),
	}
	tokenStore := &axon_utils.NullTokenStore{}
	return axon_utils.ProcessEvents("Configuration", host, port, "configuration-event-processor", &projection, prepareUnmarshal, tokenStore)
}

// Map an event name as stored in AxonServer to an object that has two aspects:
// 1. It is a proto.Message, so it can be unmarshalled from the Axon event
// 2. It is an axon_utils.Event, so it can be applied to a projection
func prepareUnmarshal(payloadType string) (event axon_utils.Event) {
	log.Printf("Configuration event processor: Payload type: %v", payloadType)
	switch payloadType {
	case "PropertyChangedEvent":
		event = &PropertyChangedEvent{}
	case "TrustedKeyAddedEvent":
		event = &TrustedKeyAddedEvent{}
	case "TrustedKeyRemovedEvent":
		event = &TrustedKeyRemovedEvent{}
	case "KeyManagerAddedEvent":
		event = &KeyManagerAddedEvent{}
	case "KeyManagerRemovedEvent":
		event = &KeyManagerRemovedEvent{}
	case "CredentialsAddedEvent":
		event = &CredentialsAddedEvent{}
	case "CredentialsRemovedEvent":
		event = &CredentialsRemovedEvent{}
	default:
		event = nil
	}
	return event
}

// Event Handlers

func (event *PropertyChangedEvent) ApplyTo(projectionWrapper interface{}) {
	projection := projectionWrapper.(*Projection)
	key := event.Property.Key
	value := event.Property.Value
	log.Printf("Configuration event handler: Set property: %v: %v", key, value)
	projection.Configuration[key] = value
}

func (event *TrustedKeyAddedEvent) ApplyTo(_ interface{}) {
	trusted.UnsafeSetTrustedKey(event.PublicKey)
}

func (event *TrustedKeyRemovedEvent) ApplyTo(_ interface{}) {
	trusted.UnsafeSetTrustedKey(getEmptyPublicKey(event.Name))
}

func (event *KeyManagerAddedEvent) ApplyTo(_ interface{}) {
	trusted.UnsafeSetKeyManager(event.PublicKey)
}

func (event *KeyManagerRemovedEvent) ApplyTo(_ interface{}) {
	trusted.UnsafeSetKeyManager(getEmptyPublicKey(event.Name))
}

func (event *CredentialsAddedEvent) ApplyTo(_ interface{}) {
	authentication.UnsafeSetCredentials(event.Credentials)
}

func (event *CredentialsRemovedEvent) ApplyTo(_ interface{}) {
	emptyCredentials := grpc_example.Credentials{
		Identifier: event.Identifier,
		Secret:     "",
	}
	authentication.UnsafeSetCredentials(&emptyCredentials)
}

// Public accessor

//noinspection GoUnusedExportedFunction
func GetProperty(key string) string {
	return projection.Configuration[key]
}

// Helper functions

func getEmptyPublicKey(name string) *grpc_example.PublicKey {
	return &grpc_example.PublicKey{
		Name:      name,
		PublicKey: "",
	}
}
