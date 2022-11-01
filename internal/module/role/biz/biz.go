package rolebiz

import (
	"fmt"
	"strings"

	"github.com/xdorro/golang-grpc-base-project/pkg/casbin"
	rolev1 "github.com/xdorro/golang-grpc-base-project/proto/pb/role/v1"
)

// Biz struct.
type Biz struct {
	// option
	casbin *casbin.Casbin
}

// NewBiz new service.
func NewBiz(casbin *casbin.Casbin) *Biz {
	b := &Biz{
		casbin: casbin,
	}

	return b
}

// FindAllRoles find all roles
func (b *Biz) FindAllRoles() (
	*rolev1.FindAllRolesResponse, error,
) {
	data := make([]*rolev1.Role, 0)

	roles := b.casbin.Client().GetAllSubjects()
	for _, role := range roles {
		data = append(data, &rolev1.Role{
			Name: role,
		})
	}

	res := &rolev1.FindAllRolesResponse{
		Data: data,
	}

	return res, nil
}

// FindRoleByName find role by name
func (b *Biz) FindRoleByName(req *rolev1.CommonNameRequest) (
	*rolev1.Role, error,
) {
	name := strings.ToLower(req.GetName())

	policies := b.casbin.Client().GetFilteredPolicy(0, name)
	if len(policies) == 0 {
		return nil, fmt.Errorf("role does not exists")
	}

	permissions := make([]string, 0)
	for _, policy := range policies {
		permissions = append(permissions, policy[1])
	}

	res := &rolev1.Role{
		Name:        name,
		Permissions: permissions,
	}

	return res, nil
}

// CreateRole create role
func (b *Biz) CreateRole(req *rolev1.CreateRoleRequest) (
	*rolev1.CommonResponse, error,
) {
	name := strings.ToLower(req.GetName())

	policies := b.casbin.Client().GetFilteredPolicy(0, name)
	if len(policies) > 0 {
		return nil, fmt.Errorf("role already exists")
	}

	for _, per := range req.GetPermissions() {
		policies = append(policies, []string{name, per})
	}

	// add policies to casbin
	_, err := b.casbin.Client().AddPolicies(policies)
	if err != nil {
		return nil, err
	}

	res := &rolev1.CommonResponse{
		Data: "success",
	}

	return res, nil
}

// UpdateRole update role
func (b *Biz) UpdateRole(req *rolev1.UpdateRoleRequest) (
	*rolev1.CommonResponse, error,
) {
	name := strings.ToLower(req.GetName())

	oldPolicies := b.casbin.Client().GetFilteredPolicy(0, name)
	if len(oldPolicies) == 0 {
		return nil, fmt.Errorf("role does not exists")
	}

	_, err := b.casbin.Client().RemovePolicies(oldPolicies)
	if err != nil {
		return nil, err
	}

	policies := make([][]string, 0)
	for _, per := range req.GetPermissions() {
		policies = append(policies, []string{name, per})
	}

	// update policies to casbin
	_, err = b.casbin.Client().AddPolicies(policies)
	if err != nil {
		return nil, err
	}

	res := &rolev1.CommonResponse{
		Data: "success",
	}

	return res, nil
}

// DeleteRole delete role
func (b *Biz) DeleteRole(req *rolev1.CommonNameRequest) (
	*rolev1.CommonResponse, error,
) {
	name := strings.ToLower(req.GetName())

	policies := b.casbin.Client().GetFilteredPolicy(0, name)
	if len(policies) == 0 {
		return nil, fmt.Errorf("role does not exists")
	}

	_, err := b.casbin.Client().RemovePolicies(policies)
	if err != nil {
		return nil, err
	}

	res := &rolev1.CommonResponse{
		Data: "success",
	}

	return res, nil
}
