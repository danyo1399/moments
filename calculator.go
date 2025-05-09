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
type Calculator struct {
	Aggregate[CalculatorState]
}

var initStateFunc = func() CalculatorState {
	return CalculatorState{0}
}
var newCalculatorAggregate = NewAggregateFactory("Calculator", initStateFunc, reducer)

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

func (c *Calculator) Subtract(val int) {
	evt := subtracted{val}
	c.Apply(evt, nil)
}

func (c *Calculator) Update(val int) {
	evt := updated{val}
	c.Apply(evt, nil)
}

func (c *Calculator) Add(val int) {
	evt := added{val}
	c.Apply(evt, nil)
}

func NewCalculatorFromEvents(id string, events []any) *Calculator {
	agg := newCalculatorAggregate(WithEvents[CalculatorState](id, events))
	calculator := Calculator{
		Aggregate: *agg,
	}
	return &calculator
}

func NewCalculator(id string) *Calculator {

	agg := newCalculatorAggregate(WithId[CalculatorState](id))
	calculator := Calculator{
		Aggregate: *agg,
	}
	return &calculator
}
