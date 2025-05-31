package moments

type Config struct {
	Aggregates map[AggregateType]AggregateConfig
	Serialiser *SnapshotSerialiser
}
type AggregateConfig struct {
	StoreStrategy     StoreStrategyType
	SnapshotFrequency int
}
