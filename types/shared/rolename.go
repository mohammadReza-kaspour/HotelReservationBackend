package shared

type RoleName string

const (
	RoleCustomer RoleName = "customer"
	RoleOwner    RoleName = "owner"
	RoleStaff    RoleName = "staff"
	RoleAdmin    RoleName = "admin"
)

func (r RoleName) Valid() bool {
	switch r {
	case RoleCustomer, RoleOwner, RoleStaff, RoleAdmin:
		return true
	}
	return false
}
