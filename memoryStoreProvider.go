package moments

import (
	"errors"
	"fmt"
	"sync/atomic"
)

type MemoryStoreProvider struct {
	state *MemoryStoreState
}

func NewMemoryStoreProvider() MemoryStoreProvider {
	return MemoryStoreProvider{state: &MemoryStoreState{tenants: map[string]*MemoryStoreTenantState{}}}
}

func (p MemoryStoreProvider) NewTenant(tenant string) error {
	state := p.state
	if _, exists := state.tenants[tenant]; exists {
		return errors.New(fmt.Sprintln("Tenant already exists", tenant))
	}
	state.tenants[tenant] = &MemoryStoreTenantState{
		streams:   map[StreamId]*Stream{},
		eventsMap: map[StreamId][]PersistedEvent{},
		events:    []PersistedEvent{},
		sequence:  atomic.Uint64{},
		snapshots: map[StreamId]Snapshot{},
	}
	return nil
}

func (p MemoryStoreProvider) DeleteTenant(tenant string) error {
	delete(p.state.tenants, tenant)
	return nil
}

func (p MemoryStoreProvider) GetStore(tenant string) (Store, error) {
	state := p.state
	tenantState, exists := state.tenants[tenant]
	if !exists {
		return nil, errors.New(fmt.Sprintln("Tenant doesnt exist", tenant))
	}
	store := NewMemoryStore(tenantState)
	return store, nil
}

func (p MemoryStoreProvider) Close() {
}
