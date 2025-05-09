package moments

import (
	"time"
)

type (
	CorrelationId string
	CausationId   string
	EventType     string
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

type EventDataVersion struct {
	Version Version
	Data    any
}

type Event struct {
	EventId   string
	Data      any
	Timestamp time.Time
}

type EventSlice struct {
	Events []Event
}

type PersistedEventSlice struct {
	Events []PersistedEvent
}

type ApplyArgs struct {
	EventId   string
	Timestamp time.Time
}

func NewEvent(data any, args *ApplyArgs) Event {
	if args == nil {
		args = &ApplyArgs{}
	}
	eventId := NewSquentialString()
	evt := Event{
		EventId:   DefaultIfEmpty(&args.EventId, eventId),
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
		EventType:      EventType(TypeName(e.Data)),
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

type EventHeader struct {
	EventId       string
	EventType     string
	CausationId   string
	CorrelationId string
	Metadata      map[string]any
}
