package permissionservice

import (
	"context"

	permissionbiz "github.com/xdorro/golang-grpc-base-project/internal/module/permission/biz"
	permissionv1 "github.com/xdorro/golang-grpc-base-project/proto/pb/permission/v1"
)

var _ permissionv1.PermissionServiceServer = &Service{}

// Service struct.
type Service struct {
	// option
	permissionBiz *permissionbiz.Biz

	permissionv1.UnimplementedPermissionServiceServer
}

// NewService new service.
func NewService(permissionBiz *permissionbiz.Biz) *Service {
	s := &Service{
		permissionBiz: permissionBiz,
	}

	return s
}

// FindAllPermissions find all permissions
func (s *Service) FindAllPermissions(_ context.Context, req *permissionv1.FindAllPermissionsRequest) (
	*permissionv1.FindAllPermissionsResponse, error,
) {
	return s.permissionBiz.FindAllPermissions(req)
}

// FindPermissionByID find permission by id
func (s *Service) FindPermissionByID(_ context.Context, req *permissionv1.CommonUUIDRequest) (
	*permissionv1.Permission, error,
) {
	return s.permissionBiz.FindPermissionByID(req)
}

// CreatePermission create permission
func (s *Service) CreatePermission(_ context.Context, req *permissionv1.CreatePermissionRequest) (
	*permissionv1.CommonResponse, error,
) {
	return s.permissionBiz.CreatePermission(req)
}

// UpdatePermission update permission by id
func (s *Service) UpdatePermission(_ context.Context, req *permissionv1.UpdatePermissionRequest) (
	*permissionv1.CommonResponse, error,
) {
	return s.permissionBiz.UpdatePermission(req)
}

// DeletePermission delete permission by id
func (s *Service) DeletePermission(_ context.Context, req *permissionv1.CommonUUIDRequest) (
	*permissionv1.CommonResponse, error,
) {
	return s.permissionBiz.DeletePermission(req)
}
