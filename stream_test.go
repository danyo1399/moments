package moments

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateStreamId(t *testing.T) {
	streamId := StreamId{Id: "1", StreamType: "Calculator"}
	assert.Equal(t, "Calculator__1", streamId.String())
}
