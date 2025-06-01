package moments

import (
	"encoding/json"
	"fmt"
	"path"
	"reflect"
)


type EventDeserialiser func(data []byte) (any, error)
type EventDeserialiserConfig map[EventType]EventDeserialiser
func (c *EventDeserialiserConfig) Deserialise(eventType EventType, data []byte) (any, error) {
	fn, ok := (*c)[eventType]
	if !ok {
		return nil, fmt.Errorf("no deserialiser for event type %v", eventType)
	}
	return fn(data)
}

// AddEventDeSerialiser creates a function that can be used to deserialise events.
func AddEventDeSerialiser[T any](config EventDeserialiserConfig) {
 fn := func(data []byte) (any, error) {
		var val T;
		err := json.Unmarshal(data, &val);
		if err != nil {
			return nil, err;
		}
		return val, nil;
	}
	var zero T
	config[GetEventType(zero)] = fn
}

func GetEventType(value any) EventType {
	ty := reflect.TypeOf(value)
	pkg := path.Base(ty.PkgPath())
	if pkg == "." {
		return EventType(ty.Name())
	}
	return EventType(pkg + "." + ty.Name())
}
