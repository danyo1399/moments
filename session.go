package moments

import (
	"errors"
)

type SessionProvider struct {
	StoreProvider StoreProvider
}
type Session struct {
	Store         Store
	CorrelationId CorrelationId
	CausationId   CausationId
	Metadata      Metadata
	tenant        string
}

func NewSessionProvider(storeProvider StoreProvider) SessionProvider {
	return SessionProvider{
		StoreProvider: storeProvider,
	}
}

func (sp *SessionProvider) NewSession(tenant string) (*Session, error) {
	store, err := sp.StoreProvider.GetStore(tenant)
	if err != nil {
		return nil, err
	}

	session := NewSession(tenant, store)
	return &session, nil
}

func NewSession(tenant string, store Store) Session {
	return Session{
		Store:         store,
		tenant:        tenant,
		CorrelationId: "",
		CausationId:   "",
		Metadata:      map[string]any{},
	}
}

type LoadAggregateArgs interface {
	StreamId() StreamId
	Version() Version
	Load(events []any)
}

func (s *Session) LoadAggregate(aggregate LoadAggregateArgs) error {
	streamId := aggregate.StreamId()
	fromVersion := aggregate.Version() + Version(1)
	events, err := s.Store.LoadEvents(LoadEventsOptions{
		StreamId:    streamId,
		FromVersion: fromVersion,
	})
	if err != nil {
		return err
	}
	aggregate.Load(AnySlice(events))
	return nil
}

func (s *Session) Save(agg AggregateEvents) error {
	events := agg.UnsavedEvents()

	version := agg.Version()
	err := s.SaveEvents(agg.StreamId(), events, &version)
	if err != nil {
		return err
	}
	agg.ClearUnsavedEvents()
	return nil
}

func (s *Session) SaveEvents(streamId StreamId, events []Event, expectedVersion *Version) error {
	if expectedVersion != nil && *expectedVersion == 0 {
		return errors.New("cannot save stream with no events")
	}
	args := SaveEventArgs{
		StreamId:        streamId,
		Events:          events,
		ExpectedVersion: expectedVersion,
		CorrelationId:   s.CorrelationId,
		CausationId:     s.CausationId,
		Metadata:        s.Metadata,
	}
	return s.Store.SaveEvents(args)
}

func (s *Session) LoadStream(streamId StreamId) ([]PersistedEvent, error) {
	events, err := s.Store.LoadEvents(LoadEventsOptions{StreamId: streamId})
	if err != nil {
		return nil, err
	}
	return events, nil
}

func (s *Session) LoadEvents(
	options LoadEventsOptions,
) ([]PersistedEvent, error) {
	return s.Store.LoadEvents(options)
}

func (s *Session) Close() {
	s.Store.Close()
}
