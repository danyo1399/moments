package test

import (
	"fmt"
	m "github.com/danyo1399/moments"
)

type Added struct {
	value int
}

type Updated struct {
	value int
}
type Subtracted struct {
	value int
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
			state.Value += e.value
		case Subtracted:
			state.Value -= e.value
		case Updated:
			state.Value = e.value
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
