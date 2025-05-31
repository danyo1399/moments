package moments

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadAndSaveSnapStrategy(t *testing.T) {
	session := createSnapshotSession(t)
	strat := StoreStrategies[AlwaysSnapshot]
	id := "123"
	(func() {
		calc := NewCalculator(id)
		calc.Update(5)
		calc.Add(2)
		err := strat.Save(calc, session)
		assert.Nil(t, err)
	})()

	loadedCalc := NewCalculator(id)
	err := strat.Load(loadedCalc, session)
	assert.Nil(t, err)
	assert.Equal(t, 7, loadedCalc.State().Value)
	assert.Equal(t, Version(2), loadedCalc.Version())

	defer session.Close()
}
