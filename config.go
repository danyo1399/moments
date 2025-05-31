package moments

type Config struct {
	Aggregates map[AggregateType]AggregateConfig
	Serialiser *SnapshotSerialiser
}
type AggregateConfig struct {
	StoreStrategy     storeStrategyType
	SnapshotFrequency int
}
