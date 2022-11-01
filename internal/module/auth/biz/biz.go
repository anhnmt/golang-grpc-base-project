package authbiz

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/sync/errgroup"

	usermodel "github.com/xdorro/golang-grpc-base-project/internal/module/user/model"
	"github.com/xdorro/golang-grpc-base-project/pkg/repo"
	authv1 "github.com/xdorro/golang-grpc-base-project/proto/pb/auth/v1"
	"github.com/xdorro/golang-grpc-base-project/utils"
)

var _ IAuthBiz = &Biz{}

// IAuthBiz auth service interface.
type IAuthBiz interface {
	Login(req *authv1.LoginRequest) (*authv1.TokenResponse, error)
	RevokeToken(req *authv1.TokenRequest) (*authv1.CommonResponse, error)
	RefreshToken(req *authv1.TokenRequest) (*authv1.TokenResponse, error)
}

// Biz struct.
type Biz struct {
	// option
	userCollection *mongo.Collection
}

// NewBiz new service.
func NewBiz(repo *repo.Repo) *Biz {
	s := &Biz{
		userCollection: repo.CollectionModel(&usermodel.User{}),
	}

	return s
}

func (s *Biz) Login(req *authv1.LoginRequest) (
	*authv1.TokenResponse, error,
) {
	filter := bson.M{
		"email": req.GetEmail(),
		"deleted_at": bson.M{
			"$exists": false,
		},
	}
	data, err := repo.FindOne[usermodel.User](s.userCollection, filter)
	if err != nil {
		return nil, err
	}

	// verify password
	if !data.ComparePassword(req.GetPassword()) {
		return nil, fmt.Errorf("password is incorrect")
	}

	// generate a new auth token
	res, err := s.generateAuthToken(data)
	if err != nil {
		return nil, err
	}

	return res, nil
}

// RevokeToken is the auth.v1.AuthBiz.RevokeToken method.
func (s *Biz) RevokeToken(req *authv1.TokenRequest) (
	*authv1.CommonResponse, error,
) {
	token := req.GetToken()

	// verify & remove old token
	_, err := s.removeAuthToken(token)
	if err != nil {
		return nil, err
	}

	res := &authv1.CommonResponse{
		Token: token,
	}

	return res, nil
}

// RefreshToken is the auth.v1.AuthBiz.RefreshToken method.
func (s *Biz) RefreshToken(req *authv1.TokenRequest) (
	*authv1.TokenResponse, error,
) {
	// verify & remove old token
	claims, err := s.removeAuthToken(req.GetToken())
	if err != nil {
		return nil, err
	}

	id := claims.Subject

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

	// generate a new auth token
	res, err := s.generateAuthToken(data)
	if err != nil {
		return nil, err
	}

	return res, nil
}

// generateAuthToken generates a new auth token for the user.
func (s *Biz) generateAuthToken(data *usermodel.User) (
	*authv1.TokenResponse, error,
) {
	uid := data.Id
	sessionID := uuid.NewString()
	now := time.Now()
	refreshExpire := now.Add(utils.RefreshExpire)
	accessExpire := now.Add(utils.AccessExpire)

	result := &authv1.TokenResponse{
		TokenType:   utils.TokenType,
		TokenExpire: accessExpire.Unix(),
	}

	var eg errgroup.Group

	// Create a new refreshToken
	eg.Go(func() error {
		var err error
		result.RefreshToken, err = utils.EncryptToken(&jwt.RegisteredClaims{
			Subject:   uid,
			ExpiresAt: jwt.NewNumericDate(refreshExpire),
			ID:        sessionID,
		})
		if err != nil {
			return err
		}

		// key := fmt.Sprintf(utils.AuthSessionKey, uid, sessionID)
		// err = redis.Set(s.redis, key, result.RefreshToken, utils.RefreshExpire)
		// if err != nil {
		// 	log.Err(err).Msg("Failed to set auth session")
		// 	return err
		// }

		return nil
	})

	// Create a new accessToken
	eg.Go(func() error {
		var err error
		result.AccessToken, err = utils.EncryptToken(&jwt.RegisteredClaims{
			Subject:   uid,
			ExpiresAt: jwt.NewNumericDate(accessExpire),
			ID:        sessionID,
			Audience:  []string{data.Role},
		})
		if err != nil {
			return err
		}

		return nil
	})

	if err := eg.Wait(); err != nil {

		return nil, err
	}

	return result, nil
}

// removeAuthToken removes the auth token from the redis.
func (s *Biz) removeAuthToken(token string) (*jwt.RegisteredClaims, error) {
	// verify refresh token
	claims, err := utils.DecryptToken(token)
	if err != nil {
		return nil, err
	}

	log.Info().
		Interface("claims", claims).
		Msg("Token decrypted")

	// // check if the refresh token is existed
	// key := fmt.Sprintf(utils.AuthSessionKey, claims.Subject, claims.ID)
	// if check := redis.Exists(s.redis, key); !check {
	// 	return nil, fmt.Errorf("token is not found")
	// }
	//
	// if err = redis.Del(s.redis, key); err != nil {
	// 	return nil, err
	// }

	return claims, nil
}
