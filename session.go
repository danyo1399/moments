package moments

import (
	"errors"
	"fmt"
)

type SessionProvider struct {
	StoreProvider StoreProvider
	Config        Config
}
type Session struct {
	Store         Store
	CorrelationId CorrelationId
	CausationId   CausationId
	Metadata      Metadata
	tenant        string
	config        Config
	persister     storeStrategy
}

func NewSessionProvider(storeProvider StoreProvider, config Config) SessionProvider {
	if config.SnapshotSerialiser == nil {
		config.SnapshotSerialiser = &JsonSnapshotSerialiser
	}
	return SessionProvider{
		StoreProvider: storeProvider,
		Config:        config,
	}
}

func (sp *SessionProvider) NewSession(tenant string) (*Session, error) {
	store, err := sp.StoreProvider.GetStore(tenant)
	if err != nil {
		return nil, err
	}

	session := newSession(tenant, store, sp.Config)
	return &session, nil
}

func newSession(tenant string, store Store, config Config) Session {
	return Session{
		Store:         store,
		tenant:        tenant,
		CorrelationId: "",
		CausationId:   "",
		Metadata:      map[string]any{},
		config:        config,
	}
}

func (s *Session) LoadAggregate(aggregate IAggregate) error {
	aggregateConfig, ok := s.config.Aggregates[aggregate.AggregateType()]
	if !ok {
		return errors.New(fmt.Sprintln("Unknown aggregate type", aggregate.AggregateType()))
	}
	aggregateStrategy := aggregateConfig.StoreStrategy

	storeStrategy, ok := storeStrategies[aggregateStrategy]
	if !ok {
		return errors.New(fmt.Sprintln("Unknown store strategy", aggregateStrategy))
	}
	return storeStrategy.load(aggregate, s)
}

func (s *Session) Save(aggregate IAggregate) error {
	aggregateConfig, ok := s.config.Aggregates[aggregate.AggregateType()]
	if !ok {
		return errors.New(fmt.Sprintln("Unknown aggregate type", aggregate.AggregateType()))
	}
	aggregateStrategy := aggregateConfig.StoreStrategy

	storeStrategy, ok := storeStrategies[aggregateStrategy]
	if !ok {
		return errors.New(fmt.Sprintln("Unknown store strategy", aggregateStrategy))
	}
	return storeStrategy.save(aggregate, s)
}

func (s *Session) newSaveEventArgs(streamId StreamId, events []Event, expectedVersion Version) SaveEventArgs {
	args := SaveEventArgs{
		StreamId:        streamId,
		Events:          events,
		ExpectedVersion: expectedVersion,
		CorrelationId:   s.CorrelationId,
		CausationId:     s.CausationId,
		Metadata:        s.Metadata,
	}
	return args
}

func (s *Session) saveEvents(streamId StreamId, events []Event, expectedVersion Version) error {
	if expectedVersion == 0 {
		return errors.New("cannot save stream with no events")
	}

	a := s.newSaveEventArgs(streamId, events, expectedVersion)
	return s.Store.SaveEvents(a)
}

func (s *Session) saveEventsWithSnapshot(
	streamId StreamId, events []Event, expectedVersion Version, snapshot *Snapshot,
) error {
	if expectedVersion == 0 {
		return errors.New("cannot save stream with no events")
	}

	a := s.newSaveEventArgs(streamId, events, expectedVersion)
	a.Snapshot = snapshot
	return s.Store.SaveEvents(a)
}

func (s *Session) LoadStream(streamId StreamId) ([]PersistedEvent, error) {
	events, err := s.Store.LoadEvents(LoadEventsArgs{StreamId: streamId})
	if err != nil {
		return nil, err
	}
	return events, nil
}

func (s *Session) LoadEvents(
	options LoadEventsArgs,
) ([]PersistedEvent, error) {
	return s.Store.LoadEvents(options)
}

func (s *Session) Close() {
	s.Store.Close()
}
