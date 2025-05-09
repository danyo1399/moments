package moments

type Snapshot[TState any] struct {
  StreamId StreamId
  Version Version
  State TState
}
