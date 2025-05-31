package moments

type StoreStrategyType int

const (
	EventSourced StoreStrategyType = iota
	AlwaysSnapshot
)

func (a StoreStrategyType) String() string {
	switch a {
	case EventSourced:
		return "EventSourced"
	case AlwaysSnapshot:
		return "AlwaysSnapshot"
	default:
		return "Unknown"
	}
}

type IStoreStrategy interface {
	Load(aggregate IAggregate, session *Session) error
	Save(aggregate IAggregate, session *Session) error
}

// StoreStrategies is a map of StoreStrategyType to IStoreStrategy
var StoreStrategies = map[StoreStrategyType]IStoreStrategy{
	EventSourced:   &EventSourcedPersistenceStrategy{},
	AlwaysSnapshot: &SnapshotStoreStrategy{},
}

type EventSourcedPersistenceStrategy struct{}

func (s *EventSourcedPersistenceStrategy) Load(aggregate IAggregate, session *Session) error {
	streamId := aggregate.StreamId()
	fromVersion := aggregate.Version() + Version(1)
	events, err := session.LoadEvents(LoadEventsOptions{
		StreamId:    streamId,
		FromVersion: fromVersion,
	})
	if err != nil {
		return err
	}
	aggregate.Load(AnySlice(events))
	return nil
}

func (s *EventSourcedPersistenceStrategy) Save(agg IAggregate, session *Session) error {
	events := agg.UnsavedEvents()

	version := agg.Version()
	err := session.saveEvents(agg.StreamId(), events, version)
	if err != nil {
		return err
	}
	agg.ClearUnsavedEvents()
	return nil
}

type SnapshotStoreStrategy struct{}

func (s *SnapshotStoreStrategy) Load(aggregate IAggregate, session *Session) error {
	streamId := aggregate.StreamId()
	state, err := session.Store.LoadSnapshot(streamId)
	if err != nil {
		return err
	}
	if state != nil {
		aggregate.loadSnapshot(state, session.config.Serialiser)
	}
	fromVersion := aggregate.Version() + Version(1)
	events, err := session.LoadEvents(LoadEventsOptions{
		StreamId:    streamId,
		FromVersion: fromVersion,
	})
	if err != nil {
		return err
	}
	aggregate.Load(AnySlice(events))
	return nil
}

func (s *SnapshotStoreStrategy) Save(agg IAggregate, session *Session) error {
	events := agg.UnsavedEvents()
	if len(events) == 0 {
		return nil
	}
	snapshot := agg.Snapshot(session.config.Serialiser)

	return session.saveEventsWithSnapshot(
		agg.StreamId(), events, agg.Version(), &snapshot,
	)
}
