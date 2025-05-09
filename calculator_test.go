package moments

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestCalculator(t *testing.T) {
	calculator := NewCalculator("")
	assert.Equal(t, 0, calculator.State().Value)
}

func TestCalculatorCalculate(t *testing.T) {
	calculator := NewCalculator("")
	calculator.Update(5)
	calculator.Add(10)
	calculator.Subtract(3)
	assert.Equal(t, 12, calculator.State().Value)
}
