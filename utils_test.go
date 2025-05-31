package moments

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMapSlice(t *testing.T) {
	input := []int{1, 2, 3}
	expected := []string{"1", "2", "3"}
	actual := mapSlice(input, func(i int) string {
		return strconv.Itoa(i)
	})
	assert.Equal(t, expected, actual)
}

func TestNewUUIDString(t *testing.T) {
	id := newSequentialUUIDString()
	assert.NotEmpty(t, id)
	assert.Equal(t, 36, len(id))
}

func TestNewString(t *testing.T) {
	id := newSquentialString()
	assert.NotEmpty(t, id)
	assert.Equal(t, 32, len(id))
}
