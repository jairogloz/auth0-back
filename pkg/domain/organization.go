package domain

import "fmt"

func ToOrganizationRole(organizationId, role string) *string {
	roleName := fmt.Sprintf("%s_%s", organizationId, role)
	return &roleName
}
