package moments

type Config struct {
	Aggregates map[AggregateType]AggregateConfig
	SnapshotSerialiser *SnapshotSerialiser
	EventDeserialiser *EventDeserialiserConfig

}
type AggregateConfig struct {
	StoreStrategy     storeStrategyType
	SnapshotFrequency int
}
