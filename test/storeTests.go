package test

import (
	"testing"

	"github.com/danyo1399/moments"
	"github.com/stretchr/testify/assert"
)

func TestLoadSaveDeleteSnapshot(t *testing.T, store moments.Store) {
	streamId := moments.StreamId{Id: "id", StreamType: "streamType"}
	snapshot := moments.Snapshot{
		Id:      moments.NewSnapshotId(streamId, 2),
		Version: 1,
		State:   []byte("test"),
	}
	store.SaveSnapshot(&snapshot)
	loadedSnapshot, err := store.LoadSnapshot(snapshot.Id)
	assert.NoError(t, err)
	err = store.DeleteSnapshot(snapshot.Id)
	assert.NoError(t, err)
	deletedSnapshot, err := store.LoadSnapshot(snapshot.Id)
	assert.NoError(t, err)

	assert.Nil(t, deletedSnapshot)
	assert.Equal(t, snapshot, *loadedSnapshot)
}
