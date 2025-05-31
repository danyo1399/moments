package moments

type StoreProvider interface {
	NewTenant(tenant string) error
	DeleteTenant(tenant string) error
	GetStore(tenant string) (Store, error)
	Close()
}
