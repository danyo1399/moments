package moments

import (
	"fmt"
	"testing"
	"github.com/stretchr/testify/assert"
)

func createSession(t *testing.T) *Session {
	var provider StoreProvider = NewMemoryStoreProvider()
	provider.CreateTenant("default")
	sessionProvider := NewSessionProvider(provider)
	session, err := sessionProvider.NewSession("default")
	if err != nil {
		t.Error(err)
	}
	return session
}

func TestShouldNotSaveWhenNoEvents(t *testing.T) {
	session := createSession(t)
	calc := NewCalculator("")
	err := session.Save(calc)
	assert.NotNil(t, err)

	defer session.Close()
}

func TestSaveAggregate(t *testing.T) {
	session := createSession(t)
	calc := NewCalculator("")
	calc.Update(5)
	calc.Add(10)

	err := session.Save(calc)
	assert.Nil(t, err)
	loadedCalc := NewCalculator(calc.id)
	err = session.LoadAggregate(&loadedCalc.Aggregate)
	assert.Nil(t, err)
	assert.Equal(t, 15, loadedCalc.State().Value)
	assert.Len(t, loadedCalc.UnsavedEvents(), 0)
	session.Close()
}

func TestUpdateMetadata(t *testing.T) {
	session := createSession(t)
	session.CorrelationId = "corr"
	session.CausationId = "caus"
	session.Metadata["key"] = "value"
	calc := NewCalculator("")
	calc.Update(5)
	calc.Add(10)

	err := session.Save(calc)
	assert.Nil(t, err)

	fmt.Printf("%+v\n", calc.StreamId())

	loadedCalc := NewCalculator(calc.Id())
	err = session.LoadAggregate(&loadedCalc.Aggregate)
	assert.Nil(t, err)
	events, err := session.LoadEvents(LoadEventsOptions{StreamId: calc.StreamId()})
	assert.Nil(t, err)

	assert.Len(t, events, 2)
	for _, evt := range events {
		assert.Equal(t, CausationId("caus"), evt.CausationId)
		assert.Equal(t, CorrelationId("corr"), evt.CorrelationId)
		assert.Equal(t, "value", evt.Metadata["key"])
	}
}

func TestLoadAndSave(t *testing.T) {
	session := createSession(t)
	calc := NewCalculator("")
	calc.Update(5)
	calc.Add(10)

	err := session.Save(calc)
	assert.Nil(t, err)

	fmt.Printf("%+v\n", calc.StreamId())

	loadedEvents, err := session.LoadStream(calc.StreamId())
	if err != nil {
		t.Log("Should load events")
		t.Fail()
	}
	loadedCalc := NewCalculatorFromEvents(
		calc.StreamId().String(), AnySlice(loadedEvents))
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
