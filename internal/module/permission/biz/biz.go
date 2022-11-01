package permissionbiz

import (
	"fmt"

	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	permissionmodel "github.com/xdorro/golang-grpc-base-project/internal/module/permission/model"
	"github.com/xdorro/golang-grpc-base-project/pkg/repo"
	permissionv1 "github.com/xdorro/golang-grpc-base-project/proto/pb/permission/v1"
	"github.com/xdorro/golang-grpc-base-project/utils"
)

var _ IPermissionBiz = &Biz{}

// IPermissionBiz permission service interface.
type IPermissionBiz interface {
	FindAllPermissions(req *permissionv1.FindAllPermissionsRequest) (
		*permissionv1.FindAllPermissionsResponse, error,
	)
	FindPermissionByID(req *permissionv1.CommonUUIDRequest) (
		*permissionv1.Permission, error,
	)
	CreatePermission(req *permissionv1.CreatePermissionRequest) (
		*permissionv1.CommonResponse, error,
	)
	UpdatePermission(req *permissionv1.UpdatePermissionRequest) (
		*permissionv1.CommonResponse, error,
	)
	DeletePermission(req *permissionv1.CommonUUIDRequest) (
		*permissionv1.CommonResponse, error,
	)
}

// Biz struct.
type Biz struct {
	// option
	permissionCollection *mongo.Collection
}

// NewBiz new service.
func NewBiz(repo *repo.Repo) *Biz {
	s := &Biz{
		permissionCollection: repo.CollectionModel(&permissionmodel.Permission{}),
	}

	return s
}

// FindAllPermissions is the permission.v1.PermissionBiz.FindAllPermissions method.
func (s *Biz) FindAllPermissions(req *permissionv1.FindAllPermissionsRequest) (
	*permissionv1.FindAllPermissionsResponse, error,
) {
	// count all permissions with filter
	filter := bson.M{
		"deleted_at": bson.M{
			"$exists": false,
		},
	}
	count, _ := repo.CountDocuments(s.permissionCollection, filter)
	limit := int64(10)
	totalPages := utils.TotalPage(count, limit)
	page := utils.CurrentPage(req.GetPage(), totalPages)

	// find all permissions with filter and option
	opt := options.
		Find().
		SetSort(bson.M{"created_at": -1}).
		SetLimit(limit).
		SetSkip((page - 1) * limit)
	data, err := repo.Find[permissionmodel.Permission](s.permissionCollection, filter, opt)
	if err != nil {
		return nil, err
	}

	res := &permissionv1.FindAllPermissionsResponse{
		TotalPage:   totalPages,
		CurrentPage: page,
		Data:        permissionmodel.PermissionsToProto(data),
	}

	return res, nil
}

// FindPermissionByID is the permission.v1.PermissionBiz.FindPermissionByID method.
func (s *Biz) FindPermissionByID(req *permissionv1.CommonUUIDRequest) (
	*permissionv1.Permission, error,
) {
	id := req.GetId()
	_, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Err(err).Msg("Failed find permission by id")
		return nil, err
	}

	opt := options.
		FindOne().
		SetSort(bson.M{"created_at": -1})
	filter := bson.M{
		"_id": id,
		"deleted_at": bson.M{
			"$exists": false,
		},
	}

	data, err := repo.FindOne[permissionmodel.Permission](s.permissionCollection, filter, opt)
	if err != nil {
		return nil, err
	}

	res := permissionmodel.PermissionToProto(data)
	return res, nil
}

// CreatePermission is the permission.v1.PermissionBiz.CreatePermission method.
func (s *Biz) CreatePermission(req *permissionv1.CreatePermissionRequest) (
	*permissionv1.CommonResponse, error,
) {
	// count all permissions with filter
	countFilter := bson.M{
		"slug": req.GetSlug(),
		"deleted_at": bson.M{
			"$exists": false,
		},
	}
	count, _ := repo.CountDocuments(s.permissionCollection, countFilter)
	if count > 0 {
		return nil, fmt.Errorf("slug already exists")
	}

	data := &permissionmodel.Permission{
		Name:        req.GetName(),
		Slug:        req.GetSlug(),
		RequireAuth: req.GetRequireAuth(),
		RequireHash: req.GetRequireHash(),
	}
	data.PreCreate()

	oid, err := repo.InsertOne(s.permissionCollection, data)
	if err != nil {
		log.Err(err).Msg("Error create permission")
		return nil, err
	}

	resID := oid.InsertedID.(string)
	res := &permissionv1.CommonResponse{
		Data: resID,
	}

	// if err = redis.Del(s.redis, utils.ListAuthPermissionsKey); err != nil {
	// 	return nil, err
	// }

	return res, nil
}

// UpdatePermission is the permission.v1.PermissionBiz.UpdatePermission method.
func (s *Biz) UpdatePermission(req *permissionv1.UpdatePermissionRequest) (
	*permissionv1.CommonResponse, error,
) {
	id := req.GetId()
	_, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Err(err).Msg("Failed find permission by id")
		return nil, err
	}

	filter := bson.M{
		"_id": id,
		"deleted_at": bson.M{
			"$exists": false,
		},
	}
	data, err := repo.FindOne[permissionmodel.Permission](s.permissionCollection, filter)
	if err != nil {
		return nil, err
	}

	// count all permissions with filter
	countFilter := bson.M{
		"_id":  bson.M{"$ne": id},
		"slug": req.GetSlug(),
		"deleted_at": bson.M{
			"$exists": false,
		},
	}
	count, _ := repo.CountDocuments(s.permissionCollection, countFilter)
	if count > 0 {
		return nil, fmt.Errorf("slug already exists")
	}

	data.Name = utils.StringCompareOrPassValue(data.Name, req.GetName())
	data.Slug = utils.StringCompareOrPassValue(data.Slug, req.GetSlug())

	if req.RequireAuth != nil {
		data.RequireAuth = req.GetRequireAuth()
	}

	if req.RequireHash != nil {
		data.RequireHash = req.GetRequireHash()
	}

	data.PreUpdate()

	opt := bson.M{"$set": data}
	if _, err = repo.UpdateOne(s.permissionCollection, filter, opt); err != nil {
		return nil, err
	}

	res := &permissionv1.CommonResponse{
		Data: req.GetId(),
	}

	// if err = redis.Del(s.redis, utils.ListAuthPermissionsKey); err != nil {
	// 	return nil, err
	// }

	return res, nil
}

// DeletePermission is the permission.v1.PermissionBiz.DeletePermission method.
func (s *Biz) DeletePermission(req *permissionv1.CommonUUIDRequest) (
	*permissionv1.CommonResponse, error,
) {
	id := req.GetId()
	_, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Err(err).Msg("Failed find permission by id")
		return nil, err
	}

	filter := bson.M{
		"_id": id,
		"deleted_at": bson.M{
			"$exists": false,
		},
	}
	// count all permissions with filter
	count, _ := repo.CountDocuments(s.permissionCollection, filter)
	if count <= 0 {
		return nil, fmt.Errorf("permission does not exists")
	}

	if _, err = repo.SoftDeleteOne(s.permissionCollection, filter); err != nil {
		return nil, err
	}

	res := &permissionv1.CommonResponse{
		Data: req.GetId(),
	}

	// if err = redis.Del(s.redis, utils.ListAuthPermissionsKey); err != nil {
	// 	return nil, err
	// }

	return res, nil
}
