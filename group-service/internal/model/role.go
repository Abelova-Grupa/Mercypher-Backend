package model

type GroupRole int16

const (
	RoleUnspecified GroupRole = 0
	RoleOwner       GroupRole = 1
	RoleMember      GroupRole = 2
)

func (r GroupRole) IsAdmin() bool {
	return r == RoleOwner
}

func (r GroupRole) IsOwner() bool {
	return r == RoleOwner
}