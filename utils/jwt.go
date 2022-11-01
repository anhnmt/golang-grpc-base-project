package utils

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

const (
	// TokenType is the type of the token.
	TokenType = "Bearer"
	// AccessExpire is the duration of the access token.
	AccessExpire = 1 * time.Hour // 1 hour
	// RefreshExpire is the duration of the refresh token.
	RefreshExpire = 1 * 24 * time.Hour // 1 day
)

// EncryptToken encrypt token
func EncryptToken(claims *jwt.RegisteredClaims) (string, error) {
	signKey, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(viper.GetString("jwt.signKey")))
	if err != nil {
		log.Err(err).Msg("Error ParseRSAPrivateKeyFromPEM")
		return "", err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS512, claims)
	tokenString, err := token.SignedString(signKey)
	if err != nil {
		log.Err(err).Msg("Error encrypt token")
		return "", err
	}

	return tokenString, nil
}

// DecryptToken decrypt token
func DecryptToken(tokenString string) (*jwt.RegisteredClaims, error) {
	verifyKey, err := jwt.ParseRSAPublicKeyFromPEM([]byte(viper.GetString("jwt.verifyKey")))
	if err != nil {
		log.Err(err).Msg("Error ParseRSAPublicKeyFromPEM")
		return nil, err
	}

	token, err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{}, func(token *jwt.Token) (
		interface{}, error,
	) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return verifyKey, nil
	})
	if err != nil {
		log.Err(err).Msg("Error decrypt token")
		return nil, err
	}

	claims, ok := token.Claims.(*jwt.RegisteredClaims)
	if ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("unexpected claims type: %T", token.Claims)
}
