package moments

import (
	"fmt"
)

type added struct {
	value int
}

type updated struct {
	value int
}
type subtracted struct {
	value int
}

type calculatorState struct {
	Value int
}
type calculator struct {
	Aggregate[calculatorState]
}

const calculatorType AggregateType = "Calculator"

var initStateFunc = func() calculatorState {
	return calculatorState{0}
}
var newCalculatorAggregate = NewAggregateFactory(calculatorType, initStateFunc, reducer)

func reducer(state calculatorState, events ...any) calculatorState {
	for _, event := range events {
		switch e := event.(type) {
		case added:
			state.Value += e.value
		case subtracted:
			state.Value -= e.value
		case updated:
			state.Value = e.value
		default:
			panic(fmt.Sprintln("unknown event type", e, event))
		}
	}
	return state
}

func (c *calculator) subtract(val int) {
	evt := subtracted{val}
	c.Apply(evt, nil)
}

func (c *calculator) update(val int) {
	evt := updated{val}
	c.Apply(evt, nil)
}

func (c *calculator) add(val int) {
	evt := added{val}
	c.Apply(evt, nil)
}

func newCalculatorFromEvents(id string, events []any) *calculator {
	agg := newCalculatorAggregate(WithEvents[calculatorState](id, events))
	calculator := calculator{
		Aggregate: *agg,
	}
	return &calculator
}

func newCalculator(id string) *calculator {
	agg := newCalculatorAggregate(WithId[calculatorState](id))
	calculator := calculator{
		Aggregate: *agg,
	}
	return &calculator
}
