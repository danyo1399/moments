package moments

import (
	"encoding/json"
	"errors"
	"fmt"
)

type MemoryStore struct {
	state  *MemoryStoreTenantState
	config *Config
}

func NewMemoryStore(state *MemoryStoreTenantState, config *Config) *MemoryStore {
	s := MemoryStore{state: state, config: config}
	return &s
}

func (s *MemoryStore) Close() {
}

func (s *MemoryStore) SaveSnapshot(snapshot *Snapshot) error {
	id := snapshot.Id
	s.state.snapshots[id] = *snapshot
	return nil
}

func (s *MemoryStore) LoadSnapshot(id SnapshotId) (*Snapshot, error) {
	ss, ok := s.state.snapshots[id]
	if !ok {
		return nil, nil
	}
	return &ss, nil
}

func (s *MemoryStore) DeleteSnapshot(id SnapshotId) error {
	delete(s.state.snapshots, id)
	return nil
}

func (s MemoryStore) SaveEvents(args SaveEventArgs) error {
	streamId := args.StreamId
	events := args.Events
	expectedVersion := args.ExpectedVersion
	correlationId := args.CorrelationId
	causationId := args.CausationId
	metadata := args.Metadata
	snapshot := args.Snapshot

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
	if expectedVersion != Version(endVersion) {
		return errors.New(fmt.Sprintln("Unexpected version. expected", expectedVersion, "actual", endVersion))
	}
	for _, evt := range events {
		seq := Sequence(state.sequence.Add(1))
		pe := evt.ToPersistedEvent(stream.StreamId, seq, seq,
			stream.Version+1, correlationId, causationId, metadata)
		streamEvents = append(streamEvents, pe)
		state.events = append(state.events, pe)
		stream.Version++
		data, err := json.Marshal(pe.Data)
		if err != nil {
			return err
		}
		pe.Data = nil
		state.eventData[seq] = data
	}
	if snapshot != nil {
		state.snapshots[snapshot.Id] = *snapshot
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
	options LoadEventArgs,
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
	re := filterSlice(events, func(evt PersistedEvent) bool {
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
		if streamId.Id != "" && evt.StreamId != streamId {
			return false
		}
		return true
	})
	if count != 0 {
		re = re[:count]
	}
	for _, evt := range re {
		data, ok := state.eventData[evt.Sequence]
		if !ok {
			return nil, fmt.Errorf("missing event data for sequence %v", evt.Sequence)
		}
		dataValue, err := s.config.EventDeserialiser.Deserialise(evt.EventType, data)
		if err != nil {
			return nil, err
		}
		evt.Data = dataValue
	}
	return re, nil
}
