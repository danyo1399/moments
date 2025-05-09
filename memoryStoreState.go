package moments

import (
	"sync/atomic"
)

type MemoryStoreState struct {
	tenants map[string]*MemoryStoreTenantState
}

type MemoryStoreTenantState struct {
	streams   map[StreamId]*Stream
	eventsMap map[StreamId][]PersistedEvent
	events    []PersistedEvent
	sequence  atomic.Uint64
}
