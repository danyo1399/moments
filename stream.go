package moments

import "fmt"

type StreamId struct {
	Id         string
	StreamType AggregateType
}

type Stream struct {
	StreamId StreamId
	Version  Version
	Deleted  bool
}

func (s StreamId) String() string {
	return fmt.Sprintf("%v:%v", s.StreamType, s.Id)
}
