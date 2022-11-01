package authservice

import (
	"context"

	authbiz "github.com/xdorro/golang-grpc-base-project/internal/module/auth/biz"
	authv1 "github.com/xdorro/golang-grpc-base-project/proto/pb/auth/v1"
)

var _ authv1.AuthServiceServer = &Service{}

// Service struct.
type Service struct {
	// option
	authBiz *authbiz.Biz

	authv1.UnimplementedAuthServiceServer
}

// NewService new service.
func NewService(authBiz *authbiz.Biz) *Service {
	s := &Service{
		authBiz: authBiz,
	}

	return s
}

// Login is the auth.v1.AuthService.Login method.
func (s *Service) Login(_ context.Context, req *authv1.LoginRequest) (
	*authv1.TokenResponse, error,
) {
	return s.authBiz.Login(req)
}

// RevokeToken is the auth.v1.AuthService.RevokeToken method.
func (s *Service) RevokeToken(_ context.Context, req *authv1.TokenRequest) (
	*authv1.CommonResponse, error,
) {
	return s.authBiz.RevokeToken(req)
}

// RefreshToken is the auth.v1.AuthService.RefreshToken method.
func (s *Service) RefreshToken(_ context.Context, req *authv1.TokenRequest) (
	*authv1.TokenResponse, error,
) {
	return s.authBiz.RefreshToken(req)
}
