package moments

import (
	"errors"
	"fmt"
)

type MemoryStore struct {
	state *MemoryStoreTenantState
}

var globalSequence = 0

func NewMemoryStore(state *MemoryStoreTenantState) MemoryStore {
	s := MemoryStore{state: state}
	return s
}

func (s MemoryStore) Close() {
}

func (s MemoryStore) SaveSnapshot(snapshot Snapshot[any]) error {
	panic("implement me")
}

func (s MemoryStore) LoadSnpashot(streamId StreamId) (Snapshot[any], error) {
	panic("implement me")
}

func (s MemoryStore) SaveEvents(args SaveEventArgs) error {
	streamId := args.StreamId
	events := args.Events
	expectedVersion := args.ExpectedVersion
	correlationId := args.CorrelationId
	causationId := args.CausationId
	metadata := args.Metadata

	state := s.state
	stream, streamExists := state.streams[streamId]

	if !streamExists {
		stream = &Stream{StreamId: streamId}
	}
	streamEvents, eventsExist := state.eventsMap[streamId]
	if !eventsExist {
		streamEvents = []PersistedEvent{}
	}
	endVersion := stream.Version + Version(len(events))
	if expectedVersion != nil && *expectedVersion != Version(endVersion) {
		return errors.New(fmt.Sprintln("Unexpected version. expected", expectedVersion, "actual", endVersion))
	}
	for _, evt := range events {
		seq := Sequence(state.sequence.Add(1))
		pe := evt.ToPersistedEvent(stream.StreamId, seq, seq, 
			stream.Version + 1, correlationId, causationId, metadata)
		streamEvents = append(streamEvents, pe)
		state.events = append(state.events, pe)
		stream.Version++
	}
	if !streamExists {
		state.streams[streamId] = stream
	}
	if !eventsExist {
		state.eventsMap[streamId] = streamEvents
	}
	return nil
}

func (s MemoryStore) LoadEvents(
	options LoadEventsOptions,
) ([]PersistedEvent, error) {
	state := s.state
	events := state.events
	fromVersion := options.FromVersion
	toVersion := options.ToVersion
	fromSequence := options.FromSequence
	toSequence := options.ToSequence
	streamId := options.StreamId
	count := options.Count

	if events == nil {
		return []PersistedEvent{}, nil
	}
	re := FilterSlice(events, func(evt PersistedEvent) bool {
		if fromVersion != 0 && evt.Version < fromVersion {
			return false
		}
		if toVersion != 0 && evt.Version > toVersion {
			return false
		}
		if fromSequence != 0 && evt.Sequence < fromSequence {
			return false
		}
		if toSequence != 0 && evt.Sequence > toSequence {
			return false
		}
		if streamId.Id != ""  && evt.StreamId != streamId {
			return false
		}
		return true
	})
	if count != 0 {
		re = re[:count]
	}
	return re, nil
}
