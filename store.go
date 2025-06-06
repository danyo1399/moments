package moments

type LoadEventArgs struct {
	StreamId     StreamId
	Count        uint
	FromVersion  Version
	ToVersion    Version
	FromSequence Sequence
	ToSequence   Sequence
	Descending   bool
}

type SaveEventArgs struct {
	StreamId        StreamId
	Events          []Event
	CorrelationId   CorrelationId
	CausationId     CausationId
	Metadata        Metadata
	ExpectedVersion Version
	Snapshot        *Snapshot
}
type Store interface {
	SnapshotStore
	SaveEvents(
		args SaveEventArgs,
	) error
	LoadEvents(options LoadEventArgs) ([]PersistedEvent, error)
	Close()
}
