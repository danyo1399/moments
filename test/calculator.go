package test

import (
	"fmt"
	m "github.com/danyo1399/moments"
)

type Added struct {
	Value int
}

type Updated struct {
	Value int
}
type Subtracted struct {
	Value int
}

type CalculatorState struct {
	Value int
}
type Calculator struct {
	m.Aggregate[CalculatorState]
}

const CalculatorType m.AggregateType = "Calculator"
var initStateFunc = func() CalculatorState {
	return CalculatorState{0}
}
var newCalculatorAggregate = m.NewAggregateFactory(CalculatorType, initStateFunc, reducer)

func reducer(state CalculatorState, events ...any) CalculatorState {
	for _, event := range events {
		switch e := event.(type) {
		case Added:
			state.Value += e.Value
		case Subtracted:
			state.Value -= e.Value
		case Updated:
			state.Value = e.Value
		default:
			panic(fmt.Sprintln("unknown event type", e, event))
		}
	}
	return state
}

func (c *Calculator) subtract(val int) {
	evt := Subtracted{val}
	c.Apply(evt, nil)
}

func (c *Calculator) update(val int) {
	evt := Updated{val}
	c.Apply(evt, nil)
}

func (c *Calculator) add(val int) {
	evt := Added{val}
	c.Apply(evt, nil)
}

func NewCalculatorFromEvents(id string, events []any) *Calculator {
	agg := newCalculatorAggregate(m.WithEvents[CalculatorState](id, events))
	calculator := Calculator{
		Aggregate: *agg,
	}
	return &calculator
}

func NewCalculator(id string) *Calculator {

	agg := newCalculatorAggregate(m.WithId[CalculatorState](id))
	calculator := Calculator{
		Aggregate: *agg,
	}
	return &calculator
}
