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
}

type Event struct {
	EventId   EventId
	Data      any
	Timestamp time.Time
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
		EventId:   defaultIfEmpty(&args.EventId, EventId(eventId)),
		Timestamp: time.Now(),
		Data:      data,
	}
	return evt
}

func (e *PersistedEvent) ToEvent() Event {
	return Event{Data: e.Data, EventId: e.EventId, Timestamp: e.Timestamp}
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
		EventType:      EventType(typeName(e.Data)),
		Version:        version,
		CorrelationId:  correlationId,
		CausationId:    causationId,
		Metadata:       metadata,
		Event: Event{
			EventId:   e.EventId,
			Data:      e.Data,
			Timestamp: e.Timestamp,
		},
	}
	return r
}

