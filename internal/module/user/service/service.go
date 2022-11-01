package userservice

import (
	"context"

	userbiz "github.com/xdorro/golang-grpc-base-project/internal/module/user/biz"
	userv1 "github.com/xdorro/golang-grpc-base-project/proto/pb/user/v1"
)

var _ userv1.UserServiceServer = &Service{}

// Service struct.
type Service struct {
	// option
	userBiz *userbiz.Biz

	userv1.UnimplementedUserServiceServer
}

// NewService new service.
func NewService(userBiz *userbiz.Biz) *Service {
	s := &Service{
		userBiz: userBiz,
	}

	return s
}

// FindAllUsers is the user.v1.UserService.FindAllUsers method.
func (s *Service) FindAllUsers(_ context.Context, req *userv1.FindAllUsersRequest) (
	*userv1.FindAllUsersResponse, error,
) {
	return s.userBiz.FindAllUsers(req)
}

// FindUserByID is the user.v1.UserService.FindUserByID method.
func (s *Service) FindUserByID(_ context.Context, req *userv1.CommonUUIDRequest) (
	*userv1.User, error,
) {
	return s.userBiz.FindUserByID(req)
}

// CreateUser is the user.v1.UserService.CreateUser method.
func (s *Service) CreateUser(_ context.Context, req *userv1.CreateUserRequest) (
	*userv1.CommonResponse, error,
) {
	return s.userBiz.CreateUser(req)
}

// UpdateUser is the user.v1.UserService.UpdateUser method.
func (s *Service) UpdateUser(_ context.Context, req *userv1.UpdateUserRequest) (
	*userv1.CommonResponse, error,
) {
	return s.userBiz.UpdateUser(req)
}

// DeleteUser is the user.v1.UserService.DeleteUser method.
func (s *Service) DeleteUser(_ context.Context, req *userv1.CommonUUIDRequest) (
	*userv1.CommonResponse, error,
) {
	return s.userBiz.DeleteUser(req)
}
