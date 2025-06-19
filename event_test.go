package moments

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGetEventTypeName(t *testing.T) {
	evtType, err := getEventTypeFromName("OrderItem_PartiallyFulfilled_V1")
	assert.NoError(t, err)
	assert.Equal(t, SchemaVersion(1), evtType.SchemaVersion)
	assert.Equal(t, "order_item", evtType.AggregateType)
	assert.Equal(t, "partially_fulfilled", evtType.Name)
	assert.Equal(t, "OrderItem_PartiallyFulfilled_V1", evtType.Id)
}

func TestNewEventWithJustData(t *testing.T) {
	evt := NewEvent("8", nil)

	assert.Equal(t, evt.EventId, evt.EventId)
	assert.NotEmpty(t, evt.Data)
	assert.NotEmpty(t, evt.EventId)
}

func TestToPersistedEvent(t *testing.T) {
	evt := NewEvent(calculator_added_v1{value: 8}, &ApplyArgs{
		EventId:   "1",
		Timestamp: time.Now(),
	})

	seq := Sequence(1)
	globalSeq := Sequence(2)
	streamId := StreamId{Id: "1", StreamType: "Calculator"}
	pe := evt.ToPersistedEvent(streamId, seq, globalSeq, Version(1), CorrelationId("co"), CausationId("ca"), Metadata{})
	assert.Equal(t, evt.EventId, pe.EventId)
	assert.Equal(t, EventType{
		SchemaVersion: SchemaVersion(1), AggregateType: "calculator", Name: "added",
		Id: "calculator_added_v1",
	},
		pe.EventType)
	assert.Equal(t, Version(1), pe.Version)
	assert.Equal(t, evt.Data, pe.Data)
	assert.Equal(t, seq, pe.Sequence)
	assert.Equal(t, streamId, pe.StreamId)
	assert.Equal(t, globalSeq, pe.GlobalSequence)
}
