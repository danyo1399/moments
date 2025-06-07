package moments

import (
	"errors"
	"fmt"
	"sync/atomic"
)

type MemoryStoreProvider struct {
	state  *MemoryStoreState
	config *Config
}

func NewMemoryStoreProvider(config *Config) *MemoryStoreProvider {
	return &MemoryStoreProvider{
		state:  &MemoryStoreState{tenants: map[TenantId]*MemoryStoreTenantState{}},
		config: config,
	}
}

func (p *MemoryStoreProvider) NewTenant(tenant TenantId) error {
	state := p.state
	if _, exists := state.tenants[tenant]; exists {
		return errors.New(fmt.Sprintln("Tenant already exists", tenant))
	}
	state.tenants[tenant] = &MemoryStoreTenantState{
		streams:   map[StreamId]*Stream{},
		eventsMap: map[StreamId][]PersistedEvent{},
		events:    []PersistedEvent{},
		sequence:  atomic.Uint64{},
		snapshots: map[SnapshotId]Snapshot{},
		eventData: make(map[Sequence][]byte),
	}
	return nil
}

func (p *MemoryStoreProvider) TenantExists(id TenantId) (bool, error) {
	_, exists := p.state.tenants[id]
	return exists, nil
}

func (p *MemoryStoreProvider) DeleteTenant(tenant TenantId) error {
	delete(p.state.tenants, tenant)
	return nil
}

func (p *MemoryStoreProvider) NewStore(tenant TenantId) (Store, error) {
	state := p.state
	tenantState, exists := state.tenants[tenant]
	if !exists {
		return nil, errors.New(fmt.Sprintln("Tenant doesnt exist", tenant))
	}
	store := NewMemoryStore(tenantState, p.config)
	return store, nil
}

func (p *MemoryStoreProvider) Close() {
}
