package moments

type (
	// Reducer is a function that applies events to a state to produce a new state.
	// It's the core of event sourcing, determining how events modify the aggregate state.
	Reducer[T any] func(state T, events ...any) T

	// Version represents the version number of an aggregate.
	// It increases with each event applied to the aggregate.
	Version       uint64
	AggregateType string

	// Sequence represents a sequence number in an event stream.
	Sequence      uint64
	SchemaVersion uint
)
