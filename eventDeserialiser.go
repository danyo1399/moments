package moments

import (
	"encoding/json"
	"fmt"
)

type (
	EventDeserialiserFunc func(data []byte) (any, error)
	EventDeserialiser     map[EventType]EventDeserialiserFunc
)

func NewEventDeserialiser() EventDeserialiser {
	return make(EventDeserialiser)
}

func (c *EventDeserialiser) Deserialise(eventType EventType, data []byte) (any, error) {
	fn, ok := (*c)[eventType]
	if !ok {
		return nil, fmt.Errorf("no deserialiser for event type %v", eventType)
	}
	return fn(data)
}

// AddJsonEventDeserialiser creates a function that can be used to deserialise events.
func AddJsonEventDeserialiser[T any](deserialiser EventDeserialiser) {
	fn := func(data []byte) (any, error) {
		var val T
		err := json.Unmarshal(data, &val)
		if err != nil {
			return nil, err
		}
		return val, nil
	}
	var zero T
	deserialiser[GetEventType(zero)] = fn
}
