package types

import (
	"fmt"
	"slices"

	"github.com/mohammadReza-kaspour/HotelReservationBackend/shared/localErrors"
	"github.com/mohammadReza-kaspour/HotelReservationBackend/types/shared"
)

type Role struct {
	Id          string              `bson:"_id,omitempty" json:"id,omitempty"`
	Name        shared.RoleName     `bson:"name" json:"name"`
	Permissions []shared.Permission `bson:"permissions" json:"permissions"`
}

func NewRole(params *CreateRoleParams) *Role {
	initPermissions := []shared.Permission{}
	if len(params.Permissions) != 0 {
		initPermissions = params.Permissions
	}

	return &Role{
		Name:        params.Name,
		Permissions: initPermissions,
	}
}
func (r *Role) HasPermission(p shared.Permission) bool {
	return slices.Contains(r.Permissions, p)
}

type CreateRoleParams struct {
	Name        shared.RoleName     `json:"name"`
	Permissions []shared.Permission `json:"permissions"`
}

func (params *CreateRoleParams) Validate() error {
	if !params.Name.Valid() {
		return localErrors.NewErrorBadRequest(fmt.Sprintf("invalid role name: %s", params.Name))
	}
	for _, p := range params.Permissions {
		if !p.Validate() {
			return localErrors.NewErrorBadRequest(fmt.Sprintf("invalid permission: %s", p))
		}
	}

	return nil
}

type UpdateRoleParams struct {
	Name        *shared.RoleName    `bson:"name" json:"name"`
	Permissions []shared.Permission `bson:"permissions" json:"permissions"`
}

func (params *UpdateRoleParams) Validate() error {
	if !params.Name.Valid() {
		return localErrors.NewErrorBadRequest(fmt.Sprintf("invalid role name: %s", params.Name))
	}
	for _, p := range params.Permissions {
		if !p.Validate() {
			return localErrors.NewErrorBadRequest(fmt.Sprintf("invalid permission: %s", p))
		}
	}

	return nil
}
func (params *UpdateRoleParams) ToUpdateMap() map[string]any {
	updateMap := map[string]any{}

	derefNullable(updateMap, "name", params.Name)

	if params.Permissions != nil {
		updateMap["permissions"] = params.Permissions
	}

	if len(updateMap) == 0 {
		return nil
	}
	return updateMap
}
