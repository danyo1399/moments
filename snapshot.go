package moments

import "encoding/json"

// Snapshot is a snapshot of an aggregate's state at a given point in time.
type Snapshot struct {
	Id      SnapshotId
	Version Version
	State   []byte
}
type SnapshotId struct {
	StreamId      StreamId
	SchemaVersion SchemaVersion
}

func NewSnapshotId(streamId StreamId, schemaVersion SchemaVersion) SnapshotId {
	return SnapshotId{
		StreamId:      streamId,
		SchemaVersion: schemaVersion,
	}
}

type SnapshotSerialiser struct {
	Marshal   func(v any) ([]byte, error)
	Unmarshal func(data []byte, v any) error
}

var JsonSnapshotSerialiser SnapshotSerialiser = SnapshotSerialiser{
	Marshal: func(v any) ([]byte, error) {
		return json.Marshal(v)
	},
	Unmarshal: func(data []byte, v any) error {
		return json.Unmarshal(data, v)
	},
}

type SnapshotStore interface {
	SaveSnapshot(snapshot *Snapshot) error
	LoadSnapshot(id SnapshotId) (*Snapshot, error)
	DeleteSnapshot(id SnapshotId) error
}
