package moments

type storeStrategyType int

const (
	eventSourced storeStrategyType = iota
	alwaysSnapshot
)

func (a storeStrategyType) String() string {
	switch a {
	case eventSourced:
		return "EventSourced"
	case alwaysSnapshot:
		return "AlwaysSnapshot"
	default:
		return "Unknown"
	}
}

type storeStrategy interface {
	load(aggregate IAggregate, session *Session) error
	save(aggregate IAggregate, session *Session) error
}

// storeStrategies is a map of StoreStrategyType to IStoreStrategy
var storeStrategies = map[storeStrategyType]storeStrategy{
	eventSourced:   &eventSourcedPersistenceStrategy{},
	alwaysSnapshot: &snapshotStoreStrategy{},
}

type eventSourcedPersistenceStrategy struct{}

func (s *eventSourcedPersistenceStrategy) load(aggregate IAggregate, session *Session) error {
	streamId := aggregate.StreamId()
	fromVersion := aggregate.Version() + Version(1)
	events, err := session.LoadEvents(LoadEventsArgs{
		StreamId:    streamId,
		FromVersion: fromVersion,
	})
	if err != nil {
		return err
	}
	aggregate.Load(anySlice(events))
	return nil
}

func (s *eventSourcedPersistenceStrategy) save(agg IAggregate, session *Session) error {
	events := agg.UnsavedEvents()

	version := agg.Version()
	err := session.saveEvents(agg.StreamId(), events, version)
	if err != nil {
		return err
	}
	agg.ClearUnsavedEvents()
	return nil
}

type snapshotStoreStrategy struct{}

func (s *snapshotStoreStrategy) load(aggregate IAggregate, session *Session) error {
	streamId := aggregate.StreamId()
	state, err := session.Store.LoadSnapshot(streamId)
	if err != nil {
		return err
	}
	if state != nil {
		aggregate.loadSnapshot(state, session.config.Serialiser)
	}
	fromVersion := aggregate.Version() + Version(1)
	events, err := session.LoadEvents(LoadEventsArgs{
		StreamId:    streamId,
		FromVersion: fromVersion,
	})
	if err != nil {
		return err
	}
	aggregate.Load(anySlice(events))
	return nil
}

func (s *snapshotStoreStrategy) save(agg IAggregate, session *Session) error {
	events := agg.UnsavedEvents()
	if len(events) == 0 {
		return nil
	}
	snapshot := agg.Snapshot(session.config.Serialiser)

	return session.saveEventsWithSnapshot(
		agg.StreamId(), events, agg.Version(), &snapshot,
	)
}
