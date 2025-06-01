package moments

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadAndSaveSnapshotStrategy(t *testing.T) {
	session := createSnapshotSession(t)
	strat := storeStrategies[alwaysSnapshot]
	id := "123"
	(func() {
		calc := newCalculator(id)
		calc.update(5)
		calc.add(2)
		err := strat.save(calc, session)
		assert.Nil(t, err)
	})()

	loadedCalc := newCalculator(id)
	err := strat.load(loadedCalc, session)
	assert.Nil(t, err)
	assert.Equal(t, 7, loadedCalc.State().Value)
	assert.Equal(t, Version(2), loadedCalc.Version())

	defer session.Close()
}
