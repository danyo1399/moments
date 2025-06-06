package moments

type Config struct {
	Aggregates map[AggregateType]AggregateConfig
	SnapshotSerialiser *SnapshotSerialiser
	EventDeserialiser *EventDeserialiser

}
type AggregateConfig struct {
	StoreStrategy     storeStrategyType
	SnapshotFrequency int
}
