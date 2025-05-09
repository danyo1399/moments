package moments

type LoadEventsOptions struct {
	StreamId     StreamId
	Count        uint
	FromVersion  Version
	ToVersion    Version
	FromSequence Sequence
	ToSequence   Sequence
}

type SaveEventArgs struct {
	StreamId        StreamId
	Events          []Event
	CorrelationId   CorrelationId
	CausationId     CausationId
	Metadata        Metadata
	ExpectedVersion *Version
}

type Store interface {
	SaveEvents(
		args SaveEventArgs,
	) error
	LoadEvents(options LoadEventsOptions) ([]PersistedEvent, error)
	SaveSnapshot(snapshot Snapshot[any]) error
	LoadSnpashot(streamId StreamId) (Snapshot[any], error)
	Close()
}
