package event_broker

type EventBroker interface {
	Emit(event Event) error
	Subscribe(eventName string, handler func(message []byte) error) error
}

type Event struct {
	Id      string       `json:"id"`
	Name    string       `json:"name"`
	Payload EventPayload `json:"payload"`
}

type EventPayload = interface{}
