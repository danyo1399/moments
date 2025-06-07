package moments

type (
	TenantId       string
	TenantProvider interface {
		NewTenant(id TenantId) error
		DeleteTenant(id TenantId) error
		TenantExists(id TenantId) (bool, error)
	}
)
