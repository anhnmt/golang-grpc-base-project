// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package main

import (
	"github.com/xdorro/golang-grpc-base-project/internal/module/auth/biz"
	"github.com/xdorro/golang-grpc-base-project/internal/module/auth/service"
	"github.com/xdorro/golang-grpc-base-project/internal/module/user/biz"
	"github.com/xdorro/golang-grpc-base-project/internal/module/user/service"
	"github.com/xdorro/golang-grpc-base-project/internal/server"
	"github.com/xdorro/golang-grpc-base-project/internal/server/gateway"
	"github.com/xdorro/golang-grpc-base-project/internal/server/grpc"
	"github.com/xdorro/golang-grpc-base-project/internal/service"
	"github.com/xdorro/golang-grpc-base-project/pkg/repo"
)

// Injectors from wire.go:

func initServer() *server.Server {
	repoRepo := repo.NewRepo()
	iUserBiz := userbiz.NewBiz(repoRepo)
	userserviceService := userservice.NewService(iUserBiz)
	biz := authbiz.NewBiz(repoRepo)
	authserviceService := authservice.NewService(biz)
	serviceService := service.NewService(repoRepo, userserviceService, authserviceService)
	grpcServer := grpc.NewGrpcServer(serviceService)
	serveMux := gateway.NewGatewayServer(serviceService)
	serverServer := server.NewServer(serviceService, grpcServer, serveMux)
	return serverServer
}
