package moments

import (
	"fmt"
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestCreateSnapshot(t *testing.T) {
	calc := NewCalculator("")
	calc.add(10)

	snap := calc.CreateSnapshot(JsonSerialiser)

	var snapState CalculatorState 
	err := JsonSerialiser.Unmarshal(snap.State, &snapState)
	assert.Nil(t, err)

	assert.Equal(t, calc.Id(), snap.StreamId.Id)
	assert.Equal(t, calc.Version(), snap.Version)
	assert.Equal(t, calc.State(), snapState)
}

func TestLoadFromEvents(t *testing.T) {
	calc := NewCalculator("")
	calc.update(5)
	calc.add(10)
	events := calc.UnsavedEvents()

	calc2 := NewCalculatorFromEvents("2", anySlice(events))
	assert.Equal(t, 15, calc2.State().Value)
}

func TestAppendEvents(t *testing.T) {
	calc := NewCalculator("")
	assert.Equal(t, 0, calc.State().Value)
	calc.add(10)
	calc.subtract(3)

	evt1 := calc.UnsavedEvents()[0]
	unsavedEvents := calc.UnsavedEvents()
	assert.Equal(t, 7, calc.State().Value)
	assert.Equal(t, Version(2), calc.Version())
	assert.Equal(t, 2, len(unsavedEvents))
	assert.NotEmpty(t, evt1.EventId)
	assert.NotEmpty(t, evt1.Timestamp)
	assert.Equal(t, typeName(evt1.Data), "moments.added")
	// assert.NotEmpty(t, evt1.CorrelationId)
}

func TestStreamIdFormatted(t *testing.T) {
	calc := NewCalculator("")
	assert.Equal(t, "Calculator__"+calc.Id(), calc.StreamId().String())
}

func TestIdNotSame(t *testing.T) {
	calc := NewCalculator("")
	calc2 := NewCalculator("")
	assert.NotEqual(t, calc.Id(), calc2.Id())
}

func TestLoadAggregate(t *testing.T) {
	calc := NewCalculator("")
	calc.update(10)
	calc.subtract(3)

	calc2 := NewCalculator("")
	fmt.Println(calc.UnsavedEvents())

	calc2.Load(anySlice(calc.UnsavedEvents()))

	assert.NotEqual(t, calc.Id(), calc2.Id())
	assert.Equal(t, Version(2), calc2.Version())
	assert.Equal(t, 7, calc2.State().Value)
	assert.Empty(t, calc2.UnsavedEvents())
}
