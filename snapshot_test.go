package moments

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSnapshotIdString(t *testing.T) {
	id := SnapshotId{
		StreamId:      StreamId{Id: "id", StreamType: "streamType"},
		SchemaVersion: 2,
	}
	assert.Equal(t, "streamType:id:2", id.String())
}
