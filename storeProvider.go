package moments

type StoreProvider interface {
	TenantProvider
	NewStore(tenant TenantId) (Store, error)
	Close()
}
