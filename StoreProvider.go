package moments

type StoreProvider interface {
	CreateTenant(tenant string) error
	DeleteTenant(tenant string) error
	GetStore(tenant string) (Store, error)
	Close()
}
