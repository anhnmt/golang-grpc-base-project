package userbiz

import (
	"fmt"
	"strings"

	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	usermodel "github.com/xdorro/golang-grpc-base-project/internal/module/user/model"
	"github.com/xdorro/golang-grpc-base-project/pkg/repo"
	userv1 "github.com/xdorro/golang-grpc-base-project/proto/pb/user/v1"
	"github.com/xdorro/golang-grpc-base-project/utils"
)

var _ IUserBiz = &Biz{}

// IUserBiz user service interface.
type IUserBiz interface {
	FindAllUsers(req *userv1.FindAllUsersRequest) (
		*userv1.FindAllUsersResponse, error,
	)
	FindUserByID(req *userv1.CommonUUIDRequest) (*userv1.User, error)
	CreateUser(req *userv1.CreateUserRequest) (*userv1.CommonResponse, error)
	UpdateUser(req *userv1.UpdateUserRequest) (*userv1.CommonResponse, error)
	DeleteUser(req *userv1.CommonUUIDRequest) (*userv1.CommonResponse, error)
}

// Biz struct.
type Biz struct {
	// option
	userCollection *mongo.Collection
}

// NewBiz new service.
func NewBiz(repo *repo.Repo) IUserBiz {
	s := &Biz{
		userCollection: repo.CollectionModel(&usermodel.User{}),
	}

	return s
}

// FindAllUsers is the user.v1.UserBiz.FindAllUsers method.
func (s *Biz) FindAllUsers(req *userv1.FindAllUsersRequest) (
	*userv1.FindAllUsersResponse, error,
) {
	// count all users with filter
	filter := bson.M{
		"deleted_at": bson.M{
			"$exists": false,
		},
	}
	count, _ := repo.CountDocuments(s.userCollection, filter)
	limit := int64(10)
	totalPages := utils.TotalPage(count, limit)
	page := utils.CurrentPage(req.GetPage(), totalPages)

	// find all genres with filter and option
	opt := options.
		Find().
		SetSort(bson.M{"created_at": -1}).
		SetProjection(bson.M{"password": 0}).
		SetLimit(limit).
		SetSkip((page - 1) * limit)

	data, err := repo.Find[usermodel.User](s.userCollection, filter, opt)
	if err != nil {
		return nil, err
	}

	res := &userv1.FindAllUsersResponse{
		TotalPage:   totalPages,
		CurrentPage: page,
		Data:        usermodel.UsersToProto(data),
	}

	return res, nil
}

// FindUserByID is the user.v1.UserBiz.FindUserByID method.
func (s *Biz) FindUserByID(req *userv1.CommonUUIDRequest) (
	*userv1.User, error,
) {
	id := req.GetId()

	opt := options.
		FindOne().
		SetSort(bson.M{"created_at": -1}).
		SetProjection(bson.M{"password": 0})

	filter := bson.M{
		"_id": id,
		"deleted_at": bson.M{
			"$exists": false,
		},
	}

	data, err := repo.FindOne[usermodel.User](s.userCollection, filter, opt)
	if err != nil {
		return nil, err
	}

	res := usermodel.UserToProto(data)
	return res, nil
}

// CreateUser is the user.v1.UserBiz.CreateUser method.
func (s *Biz) CreateUser(req *userv1.CreateUserRequest) (
	*userv1.CommonResponse, error,
) {
	// count all users with filter
	count, _ := repo.CountDocuments(s.userCollection, bson.M{
		"email": req.GetEmail(),
	})
	if count > 0 {
		return nil, fmt.Errorf("email already exists")
	}

	role := req.GetRole()
	if role == "" {
		role = "user"
	}

	data := &usermodel.User{
		Name:     req.GetName(),
		Email:    req.GetEmail(),
		Password: req.GetPassword(),
		Role:     strings.ToLower(role),
	}
	data.PreCreate()

	// hash password
	err := data.HashPassword()
	if err != nil {
		return nil, err
	}

	result, err := repo.InsertOne(s.userCollection, data)
	if err != nil {
		log.Err(err).Msg("Error create user")
		return nil, err
	}

	res := &userv1.CommonResponse{
		Data: "success",
	}

	switch v := result.InsertedID.(type) {
	case primitive.ObjectID:
		res.Data = v.Hex()
	case string:
		res.Data = v
	}

	return res, nil
}

// UpdateUser is the user.v1.UserBiz.UpdateUser method.
func (s *Biz) UpdateUser(req *userv1.UpdateUserRequest) (
	*userv1.CommonResponse, error,
) {
	id := req.GetId()

	filter := bson.M{
		"_id": id,
		"deleted_at": bson.M{
			"$exists": false,
		},
	}

	data, err := repo.FindOne[usermodel.User](s.userCollection, filter)
	if err != nil {
		return nil, err
	}

	// count all users with filter
	count, _ := repo.CountDocuments(s.userCollection, bson.M{
		"_id":   bson.M{"$ne": id},
		"email": req.GetEmail(),
	})
	if count > 0 {
		return nil, fmt.Errorf("email already exists")
	}

	data.Name = utils.StringCompareOrPassValue(data.Name, req.GetName())
	data.Email = utils.StringCompareOrPassValue(data.Email, req.GetEmail())
	data.Role = utils.StringCompareOrPassValue(data.Role, strings.ToLower(req.GetRole()))
	data.PreUpdate()

	obj := bson.M{"$set": data}
	if _, err = repo.UpdateOne(s.userCollection, filter, obj); err != nil {
		return nil, err
	}

	res := &userv1.CommonResponse{
		Data: req.GetId(),
	}
	return res, nil
}

// DeleteUser is the user.v1.UserBiz.DeleteUser method.
func (s *Biz) DeleteUser(req *userv1.CommonUUIDRequest) (
	*userv1.CommonResponse, error,
) {
	id := req.GetId()

	filter := bson.M{
		"_id": id,
		"deleted_at": bson.M{
			"$exists": false,
		},
	}

	// count all users with filter
	count, _ := repo.CountDocuments(s.userCollection, filter)
	if count <= 0 {
		return nil, fmt.Errorf("user does not exists")
	}

	if _, err := repo.SoftDeleteOne(s.userCollection, filter); err != nil {
		return nil, err
	}

	res := &userv1.CommonResponse{
		Data: req.GetId(),
	}
	return res, nil
}
