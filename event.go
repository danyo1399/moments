package moments

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/danyo1399/gotils"
)

type EventType struct {
	SchemaVersion SchemaVersion
	AggregateType string
	Name          string
}

type (
	CorrelationId string
	CausationId   string
	EventId       string
	Metadata      map[string]any
)

type PersistedEvent struct {
	Event
	StreamId       StreamId
	Sequence       Sequence
	GlobalSequence Sequence
	CausationId    CausationId
	CorrelationId  CorrelationId
	Metadata       Metadata
	EventType      EventType
	Version        Version
	Timestamp      time.Time
}

type Event struct {
	EventId EventId
	Data    any
}

type ApplyArgs struct {
	EventId   EventId
	Timestamp time.Time
}

func getEventTypeFromName(name string) (*EventType, error) {
	parts := strings.Split(name, "_")
	if len(parts) != 3 {
		return nil, fmt.Errorf("invalid event type name %v", name)
	}
	version, err := strconv.Atoi(parts[2][1:])
	if err != nil {
		return nil, fmt.Errorf("invalid event type version %v", parts[2])
	}
	return &EventType{
		SchemaVersion: SchemaVersion(version),
		AggregateType: gotils.ToSnakeCase(parts[0]),
		Name:          gotils.ToSnakeCase(parts[1]),
	}, nil
}

func GetEventType(value any) (*EventType, error) {
	ty := reflect.TypeOf(value)
	name := ty.Name()
	return getEventTypeFromName(name)
}

func NewEvent(data any, args *ApplyArgs) Event {
	if args == nil {
		args = &ApplyArgs{}
	}
	eventId := newSquentialString()
	evt := Event{
		EventId: defaultIfEmpty(&args.EventId, EventId(eventId)),
		Data:    data,
	}
	return evt
}

func (e *PersistedEvent) ToEvent() Event {
	return Event{Data: e.Data, EventId: e.EventId}
}

func (e *Event) ToPersistedEvent(
	streamId StreamId, sequence Sequence, globalSequence Sequence,
	version Version, correlationId CorrelationId,
	causationId CausationId, metadata Metadata,
) PersistedEvent {
	eventType, err := e.EventType()
	if err != nil {
		panic(err)
	}
	r := PersistedEvent{
		StreamId:       streamId,
		Sequence:       sequence,
		GlobalSequence: globalSequence,
		EventType:      *eventType,
		Version:        version,
		CorrelationId:  correlationId,
		CausationId:    causationId,
		Metadata:       metadata,
		Timestamp:      time.Now(),
		Event: Event{
			EventId: e.EventId,
			Data:    e.Data,
		},
	}
	return r
}

func (e *Event) EventType() (*EventType, error) {
	return GetEventType(e.Data)
}
