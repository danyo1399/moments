package moments

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func createEventDeserialiser() *EventDeserialiserConfig {
	var cfg EventDeserialiserConfig = make(map[EventType]EventDeserialiser)
	AddEventDeSerialiser[added](cfg)
	AddEventDeSerialiser[subtracted](cfg)
	AddEventDeSerialiser[updated](cfg)
	return &cfg
}

func createEventSourcedSession(t *testing.T) *Session {
	config := Config{
		Aggregates: map[AggregateType]AggregateConfig{
			"Calculator": {StoreStrategy: eventSourced},
		},
		EventDeserialiser: createEventDeserialiser(),
		SnapshotSerialiser: &JsonSnapshotSerialiser,
	}
	var provider StoreProvider = NewMemoryStoreProvider(&config)
	provider.NewTenant("default")
	sessionProvider := NewSessionProvider(provider, config)
	session, err := sessionProvider.NewSession("default")
	if err != nil {
		t.Error(err)
	}
	return session
}


func createSnapshotSession(t *testing.T) *Session {
	config := Config{
		Aggregates: map[AggregateType]AggregateConfig{
			"Calculator": {StoreStrategy: alwaysSnapshot},
		},
		EventDeserialiser: createEventDeserialiser(),
		SnapshotSerialiser: &JsonSnapshotSerialiser,
	}
	var provider StoreProvider = NewMemoryStoreProvider(&config)
	provider.NewTenant("default")
	sessionProvider := NewSessionProvider(provider, config)
	session, err := sessionProvider.NewSession("default")
	if err != nil {
		t.Error(err)
	}
	return session
}

func TestLoadAndSaveSnapshot(t *testing.T) {
	session := createSnapshotSession(t)
	id := "123"
	(func() {
		calc := newCalculator(id)
		calc.update(5)
		calc.add(2)
		err := session.Save(calc)
		assert.Nil(t, err)
	})()

	loadedCalc := newCalculator(id)
	err := session.LoadAggregate(loadedCalc)
	assert.Nil(t, err)
	assert.Equal(t, 7, loadedCalc.State().Value)
	assert.Equal(t, Version(2), loadedCalc.Version())

	defer session.Close()
}

func TestAggregateHasNoUpdates_LoadAggregate_NoChanges(t *testing.T) {
	session := createEventSourcedSession(t)
	calc := newCalculator("")
	calc.update(5)
	err := session.Save(calc)
	assert.Nil(t, err)

	err = session.LoadAggregate(calc)
	assert.Nil(t, err)
	assert.Equal(t, 5, calc.State().Value)

	defer session.Close()
}

func TestShouldNotSaveWhenNoEvents(t *testing.T) {
	session := createEventSourcedSession(t)
	calc := newCalculator("")
	err := session.Save(calc)
	assert.NotNil(t, err)

	defer session.Close()
}

func TestSaveAggregate(t *testing.T) {
	session := createEventSourcedSession(t)
	calc := newCalculator("")
	calc.update(5)
	calc.add(10)

	err := session.Save(calc)
	assert.Nil(t, err)
	loadedCalc := newCalculator(calc.id)
	err = session.LoadAggregate(&loadedCalc.Aggregate)
	assert.Nil(t, err)
	assert.Equal(t, 15, loadedCalc.State().Value)
	assert.Len(t, loadedCalc.UnsavedEvents(), 0)
	assert.Equal(t, calc.version, loadedCalc.version)
	assert.Equal(t, calc.schemaVersion, loadedCalc.schemaVersion)
	assert.Equal(t, calc.id, loadedCalc.id)
	assert.Equal(t, calc.aggregateType, loadedCalc.aggregateType)
	session.Close()
}

func TestUpdateMetadata(t *testing.T) {
	session := createEventSourcedSession(t)
	session.CorrelationId = "corr"
	session.CausationId = "caus"
	session.Metadata["key"] = "value"
	calc := newCalculator("")
	calc.update(5)
	calc.add(10)

	err := session.Save(calc)
	assert.Nil(t, err)

	fmt.Printf("%+v\n", calc.StreamId())

	loadedCalc := newCalculator(calc.Id())
	err = session.LoadAggregate(&loadedCalc.Aggregate)
	assert.Nil(t, err)
	events, err := session.LoadEvents(LoadEventArgs{StreamId: calc.StreamId()})
	assert.Nil(t, err)

	assert.Len(t, events, 2)
	for _, evt := range events {
		assert.Equal(t, CausationId("caus"), evt.CausationId)
		assert.Equal(t, CorrelationId("corr"), evt.CorrelationId)
		assert.Equal(t, "value", evt.Metadata["key"])
	}
}

func TestLoadAndSave(t *testing.T) {
	session := createEventSourcedSession(t)
	calc := newCalculator("")
	calc.update(5)
	calc.add(10)

	err := session.Save(calc)
	assert.Nil(t, err)

	fmt.Printf("%+v\n", calc.StreamId())

	loadedEvents, err := session.LoadStream(calc.StreamId())
	if err != nil {
		t.Log("Should load events")
		t.Fail()
	}
	loadedCalc := newCalculatorFromEvents(
		calc.StreamId().String(), anySlice(loadedEvents))
	if err != nil {
		t.Log("Should load events", err)
		t.Fail()
	}
	assert.Equal(t, Version(2), loadedCalc.Version())
	assert.Equal(t, 15, loadedCalc.State().Value)
	assert.Len(t, loadedEvents, 2)
	assert.Equal(t, Version(1), loadedEvents[0].Version)
	session.Close()
}
