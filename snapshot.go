package moments

import "encoding/json"

// Snapshot is a snapshot of an aggregate's state at a given point in time.
type Snapshot struct {
	SchemaVersion SchemaVersion
	StreamId      StreamId
	Version       Version
	State         []byte
}

type SnapshotSerialiser struct {
	Marshal func(v any) ([]byte, error)
	Unmarshal func(data []byte, v any) error
}

var JsonSerialiser SnapshotSerialiser = SnapshotSerialiser{
	Marshal: func(v any) ([]byte, error) {
		return json.Marshal(v)
	},
	Unmarshal: func(data []byte, v any) error {
		return json.Unmarshal(data, v)
	},
}
