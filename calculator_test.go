package moments

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestCalculator(t *testing.T) {
	calculator := newCalculator("")
	assert.Equal(t, 0, calculator.State().Value)
}

func TestCalculatorCalculate(t *testing.T) {
	calculator := newCalculator("")
	calculator.update(5)
	calculator.add(10)
	calculator.subtract(3)
	assert.Equal(t, 12, calculator.State().Value)
}
