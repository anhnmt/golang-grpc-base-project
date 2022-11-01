package model

import (
	"go.mongodb.org/mongo-driver/mongo"

	permissionv1 "github.com/xdorro/golang-grpc-base-project/proto/pb/permission/v1"
	"github.com/xdorro/golang-grpc-base-project/utils"
)

var _ IPermission = &Permission{}

// IPermission is the interface for a user
type IPermission interface {
	utils.IBaseModel
}

// Permission is a user struct.
type Permission struct {
	utils.BaseModel `bson:",inline"`

	Name        string `json:"name,omitempty" bson:"name,omitempty"`
	Slug        string `json:"slug,omitempty" bson:"slug,omitempty"`
	RequireAuth bool   `json:"require_auth,omitempty" bson:"require_auth,omitempty"`
	RequireHash bool   `json:"require_hash,omitempty" bson:"require_hash,omitempty"`
}

// CollectionName returns the name of the collection from struct name
func (m *Permission) CollectionName() string {
	return utils.CollectionName(m)
}

// GetIndexModels returns the index models
func (m *Permission) GetIndexModels() []mongo.IndexModel {
	return []mongo.IndexModel{}
}

// PreCreate is a callback that gets called before creating a models.
func (m *Permission) PreCreate() {
	m.BaseModel.PreCreate()
}

// PreUpdate is a callback that gets called before updating a models.
func (m *Permission) PreUpdate() {
	m.BaseModel.PreUpdate()
}

// PermissionToProto converts a user to a proto
func PermissionToProto(m *Permission) *permissionv1.Permission {
	return &permissionv1.Permission{
		Id:          m.Id,
		Name:        m.Name,
		Slug:        m.Slug,
		RequireAuth: m.RequireAuth,
		RequireHash: m.RequireHash,
	}
}

// PermissionsToProto converts a slice of users to a slice of proto
func PermissionsToProto(list []*Permission) []*permissionv1.Permission {
	return utils.ToProto[Permission, permissionv1.Permission](list, PermissionToProto)
}
