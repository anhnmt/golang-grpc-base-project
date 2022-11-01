package roleservice

import (
	"context"

	rolebiz "github.com/xdorro/golang-grpc-base-project/internal/module/role/biz"
	rolev1 "github.com/xdorro/golang-grpc-base-project/proto/pb/role/v1"
)

var _ rolev1.RoleServiceServer = &Service{}

// Service struct.
type Service struct {
	// option
	roleBiz *rolebiz.Biz

	rolev1.UnimplementedRoleServiceServer
}

// NewService new service.
func NewService(roleBiz *rolebiz.Biz) *Service {
	s := &Service{
		roleBiz: roleBiz,
	}

	return s
}

// FindAllRoles is the role.v1.RoleService.FindAllRoles method.
func (s *Service) FindAllRoles(_ context.Context, _ *rolev1.FindAllRolesRequest) (
	*rolev1.FindAllRolesResponse, error,
) {
	return s.roleBiz.FindAllRoles()
}

// FindRoleByName is the role.v1.RoleService.FindRoleByName method.
func (s *Service) FindRoleByName(_ context.Context, req *rolev1.CommonNameRequest) (
	*rolev1.Role, error,
) {
	return s.roleBiz.FindRoleByName(req)
}

// CreateRole is the role.v1.RoleService.CreateRole method.
func (s *Service) CreateRole(_ context.Context, req *rolev1.CreateRoleRequest) (
	*rolev1.CommonResponse, error,
) {
	return s.roleBiz.CreateRole(req)
}

// UpdateRole is the role.v1.RoleService.UpdateRole method.
func (s *Service) UpdateRole(_ context.Context, req *rolev1.UpdateRoleRequest) (
	*rolev1.CommonResponse, error,
) {
	return s.roleBiz.UpdateRole(req)
}

// DeleteRole is the role.v1.RoleService.DeleteRole method.
func (s *Service) DeleteRole(_ context.Context, req *rolev1.CommonNameRequest) (
	*rolev1.CommonResponse, error,
) {
	return s.roleBiz.DeleteRole(req)
}
