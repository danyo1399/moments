package test

import (
	"testing"

	"github.com/danyo1399/moments"
	"github.com/stretchr/testify/assert"
)


func TestLoadSaveDeleteSnapshot(t *testing.T, store moments.Store) {

	snapshot := moments.Snapshot{
		StreamId:      moments.StreamId{Id: "id", StreamType: "streamType"},
		Version:       1,
		SchemaVersion: 2,
		State:         []byte("test"),
	}
	store.SaveSnapshot(&snapshot)
	loadedSnapshot, err := store.LoadSnapshot(snapshot.StreamId)
	assert.NoError(t, err)
	err = store.DeleteSnapshot(snapshot.StreamId)
	assert.NoError(t, err)
	deletedSnapshot, err := store.LoadSnapshot(snapshot.StreamId)
	assert.NoError(t, err)

	assert.Nil(t, deletedSnapshot)
	assert.Equal(t, snapshot, *loadedSnapshot)
}
