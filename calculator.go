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

type CalculatorState struct {
	Value int
}
type calculator struct {
	Aggregate[CalculatorState]
}

const calculatorType AggregateType = "Calculator"
var initStateFunc = func() CalculatorState {
	return CalculatorState{0}
}
var newCalculatorAggregate = NewAggregateFactory(calculatorType, initStateFunc, reducer)

func reducer(state CalculatorState, events ...any) CalculatorState {
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

func NewCalculatorFromEvents(id string, events []any) *calculator {
	agg := newCalculatorAggregate(WithEvents[CalculatorState](id, events))
	calculator := calculator{
		Aggregate: *agg,
	}
	return &calculator
}

func NewCalculator(id string) *calculator {

	agg := newCalculatorAggregate(WithId[CalculatorState](id))
	calculator := calculator{
		Aggregate: *agg,
	}
	return &calculator
}
