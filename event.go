package moments

import (
	"time"
)

type (
	CorrelationId string
	CausationId   string
	EventType     string
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
	r := PersistedEvent{
		StreamId:       streamId,
		Sequence:       sequence,
		GlobalSequence: globalSequence,
		EventType:      e.EventType(),
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

func (e *Event) EventType() EventType {
	return GetEventType(e.Data)
}
