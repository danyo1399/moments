package moments

import (
	"fmt"
)

type (
	Reducer[T any] func(state T, events ...any) T
	Version        uint64
	Sequence       uint64

	Aggregate[TState any] struct {
		// used for snapshots
		version       Version
		state         TState
		reducer       Reducer[TState]
		unsavedEvents []Event
		aggregateType string
		id            string
	}
)

type AggregateEvents interface {
	UnsavedEvents() []Event
	Version() Version
	StreamId() StreamId
	ClearUnsavedEvents()
	AggregateType() string
}

func (a *Aggregate[TState]) Id() string {
	return a.id
}

func (a *Aggregate[TState]) Version() Version {
	return a.version
}

func (a *Aggregate[TState]) AggregateType() string {
	return a.aggregateType
}

func (a *Aggregate[TState]) State() TState {
	return a.state
}

func (a *Aggregate[TState]) UnsavedEvents() []Event {
	clone := make([]Event, len(a.unsavedEvents))
	copy(clone, a.unsavedEvents)
	return clone
}

func (a *Aggregate[TState]) ClearUnsavedEvents() {
	a.unsavedEvents = []Event{}
}

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

func (a *Aggregate[T]) load(events []Event) {
	ed := MapSlice(events, func(e Event) any {
		return e.Data
	})

	a.state = a.reducer(a.state, ed...)
	a.version += Version(len(ed))
}

func (a *Aggregate[T]) Apply(data any, args *ApplyArgs) T {
	evt := NewEvent(data, args)
	a.state = a.reducer(a.state, data)
	a.version++
	a.unsavedEvents = append(a.unsavedEvents, evt)
	return a.state
}

func (a *Aggregate[T]) StreamId() StreamId {
	return StreamId{Id: a.id, StreamType: a.aggregateType}
}

type (
	NewOption[T any]        func(args *Aggregate[T])
	InitialStateFunc[T any] func() T
	newAggregateFunc[T any] func(options ...NewOption[T]) *Aggregate[T]
)

func WithSnapshot[T any](id string, snapshot T, version Version) NewOption[T] {
	return func(a *Aggregate[T]) {
		a.state = snapshot
		a.version = version
		a.id = id
	}
}

func WithEvents[T any](id string, events []any) NewOption[T] {
	return func(a *Aggregate[T]) {
		a.id = id
		a.Load(events)
	}
}

func WithId[T any](id string) NewOption[T] {
	return func(a *Aggregate[T]) {
		if id != "" {
			a.id = id
		}
	}
}

func NewAggregateFactory[T any](
	aggregateType string, initial InitialStateFunc[T], reducer Reducer[T],
) newAggregateFunc[T] {
	return func(options ...NewOption[T]) *Aggregate[T] {
		return newAggregate(aggregateType, initial(), reducer, options...)
	}
}

func newAggregate[T any](aggregateType string, initial T, reducer Reducer[T], opts ...NewOption[T]) *Aggregate[T] {
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

func (a *Aggregate[T]) CreateSnapshot() Snapshot[T] {
	clonedState := DeepCopyJson(a.state)
	return Snapshot[T]{
		StreamId: a.StreamId(),
		Version:  a.version,
		State:    clonedState,
	}
}
