// Package moments provides an event sourcing framework for Go applications.
package moments

import (
	"fmt"
	"log/slog"
)

type (
	// Aggregate is the core structure for event sourcing.
	// It maintains state, tracks events, and provides methods for event application.
	// TState is the type of the state maintained by this aggregate.
	Aggregate[TState any] struct {
		// schemaVersion tracks the current version of the aggregate schema. Used to verify snapshot state
		schemaVersion SchemaVersion
		// version tracks the current version of the aggregate
		version Version
		// state holds the current state of the aggregate
		state TState
		// reducer is the function that applies events to the state
		reducer Reducer[TState]
		// unsavedEvents contains events that have been applied but not yet persisted
		unsavedEvents []Event
		// aggregateType is the type name of this aggregate
		aggregateType AggregateType
		// id is the unique identifier for this aggregate instance
		id string
	}
)

type IAggregate interface {
	Load(events []any)
	// Version returns the current version of the aggregate
	Version() Version
	// UnsavedEvents returns events that have been applied but not yet persisted
	UnsavedEvents() []Event

	// StreamId returns the unique identifier for the event stream
	StreamId() StreamId

	// ClearUnsavedEvents removes all unsaved events from the aggregate
	ClearUnsavedEvents()

	// AggregateType returns the type name of the aggregate
	AggregateType() AggregateType

	loadSnapshot(snapshot *Snapshot, serialiser *SnapshotSerialiser)
	Snapshot(serialiser *SnapshotSerialiser) Snapshot

	HasUnsavedChanges() bool
	SchemaVersion() SchemaVersion
}

func (a *Aggregate[TState]) Snapshot(serialiser *SnapshotSerialiser) Snapshot {
	state, err := serialiser.Marshal(a.state)
	if err != nil {
		slog.Error("failed to marshal state", "err", err)
	}

	return Snapshot{
		StreamId:      a.StreamId(),
		Version:       a.version,
		State:         state,
		SchemaVersion: a.schemaVersion,
	}
}

func (a *Aggregate[TState]) SchemaVersion() SchemaVersion {
	return a.schemaVersion
}

func (a *Aggregate[TState]) loadSnapshot(snapshot *Snapshot, serialiser *SnapshotSerialiser) {
	if snapshot.SchemaVersion != a.schemaVersion {
		panic(fmt.Sprintf("snapshot schema version %v does not match aggregate schema version %v",
			snapshot.SchemaVersion, a.schemaVersion))
	}
	if snapshot.StreamId != a.StreamId() {
		panic(fmt.Sprintf("snapshot stream id %v does not match aggregate stream id %v", snapshot.StreamId, a.StreamId()))
	}
	if a.HasUnsavedChanges() {
		panic("cannot load snapshot into aggregate with unsaved changes")
	}

	if a.version > 0 {
		panic("cannot load snapshot into an already loaded aggregate")
	}
	err := serialiser.Unmarshal(snapshot.State, &a.state)
	if err != nil {
		slog.Error("failed to unmarshal state", "err", err)
		panic(err)
	}
	a.version = snapshot.Version
}

// Id returns the unique identifier of the aggregate.
func (a *Aggregate[TState]) Id() string {
	return a.id
}

// Version returns the current version of the aggregate.
func (a *Aggregate[TState]) Version() Version {
	return a.version
}

// AggregateType returns the type name of the aggregate.
func (a *Aggregate[TState]) AggregateType() AggregateType {
	return a.aggregateType
}

// State returns the current state of the aggregate.
func (a *Aggregate[TState]) State() TState {
	return a.state
}

func (a *Aggregate[TState]) HasUnsavedChanges() bool {
	return len(a.unsavedEvents) > 0
}

// UnsavedEvents returns a copy of all events that have been applied but not yet persisted.
// Returns a clone to prevent external modification of the internal events slice.
func (a *Aggregate[TState]) UnsavedEvents() []Event {
	clone := make([]Event, len(a.unsavedEvents))
	copy(clone, a.unsavedEvents)
	return clone
}

// ClearUnsavedEvents removes all unsaved events from the aggregate.
// This is typically called after events have been successfully persisted.
func (a *Aggregate[TState]) ClearUnsavedEvents() {
	a.unsavedEvents = []Event{}
}

// Load applies a slice of events to the aggregate.
// It accepts both Event and PersistedEvent types, converting them as needed.
// This method is typically used to reconstruct an aggregate from its event history.
func (a *Aggregate[TState]) Load(events []any) {
	ed := []Event{}
	for _, evt := range events {
		switch evt := evt.(type) {
		case PersistedEvent:
			ed = append(ed, evt.ToEvent())
		case Event:
			ed = append(ed, evt)
		default:
			panic(fmt.Sprintln("unknown event type", evt))
		}
	}
	a.load(ed)
}

// load is an internal method that applies a slice of Event objects to the aggregate.
// It extracts the data from each event, applies it using the reducer, and updates the version.
func (a *Aggregate[T]) load(events []Event) {
	ed := MapSlice(events, func(e Event) any {
		return e.Data
	})

	a.state = a.reducer(a.state, ed...)
	a.version += Version(len(ed))
}

// Apply creates a new event from the provided data and applies it to the aggregate.
// It updates the state using the reducer, increments the version, and adds the event to unsavedEvents.
// Returns the new state after applying the event.
func (a *Aggregate[T]) Apply(data any, args *ApplyArgs) T {
	evt := NewEvent(data, args)
	a.state = a.reducer(a.state, data)
	a.version++
	a.unsavedEvents = append(a.unsavedEvents, evt)
	return a.state
}

// StreamId returns the unique identifier for this aggregate's event stream.
// The StreamId combines the aggregate's id and type.
func (a *Aggregate[T]) StreamId() StreamId {
	return StreamId{Id: a.id, StreamType: a.aggregateType}
}

type (
	// NewOption is a function that configures an Aggregate during creation.
	// It follows the functional options pattern for flexible aggregate configuration.
	NewOption[T any] func(args *Aggregate[T])

	// InitialStateFunc returns the initial state for a new aggregate.
	// This allows different aggregate types to have different initial states.
	InitialStateFunc[T any] func() T

	// newAggregateFunc is a factory function type that creates new Aggregate instances.
	// It accepts configuration options and returns a configured Aggregate.
	newAggregateFunc[T any] func(options ...NewOption[T]) *Aggregate[T]
)

// WithEvents creates an option to initialize an aggregate by loading a sequence of events.
// It sets the aggregate's id and applies the provided events to build the state.
func WithEvents[T any](id string, events []any) NewOption[T] {
	return func(a *Aggregate[T]) {
		a.id = id
		a.Load(events)
	}
}

// WithId creates an option to set the id of an aggregate.
// If the provided id is empty, the aggregate's existing id is retained.
func WithId[T any](id string) NewOption[T] {
	return func(a *Aggregate[T]) {
		if id != "" {
			a.id = id
		}
	}
}

// NewAggregateFactory creates a factory function for producing aggregates of a specific type.
// It encapsulates the aggregate type, initial state function, and reducer function.
// The returned factory function can be used to create new aggregate instances with various options.
func NewAggregateFactory[T any](
	aggregateType AggregateType, initial InitialStateFunc[T], reducer Reducer[T],
) newAggregateFunc[T] {
	return func(options ...NewOption[T]) *Aggregate[T] {
		return newAggregate(aggregateType, initial(), reducer, options...)
	}
}

// newAggregate creates and configures a new Aggregate instance.
// It initializes the aggregate with the provided type, state, and reducer,
// then applies any configuration options.
func newAggregate[T any](aggregateType AggregateType, initial T, reducer Reducer[T], opts ...NewOption[T]) *Aggregate[T] {
	a := Aggregate[T]{
		aggregateType: aggregateType,
		state:         initial,
		reducer:       reducer,
		version:       0,
		id:            NewSquentialString(),
	}
	for _, opt := range opts {
		opt(&a)
	}
	return &a
}

// CreateSnapshot generates a snapshot of the aggregate's current state.
// It creates a deep copy of the state to ensure the snapshot is immutable.
// The snapshot includes the stream ID, current version, and state.
func (a *Aggregate[T]) CreateSnapshot(serialiser SnapshotSerialiser) Snapshot {
	state, err := serialiser.Marshal(a.state)
	if err != nil {
		slog.Error("failed to marshal state", "err", err)
		panic(err)
	}

	return Snapshot{
		StreamId: a.StreamId(),
		Version:  a.version,
		State:    state,
	}
}
