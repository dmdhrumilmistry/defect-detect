package auth

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/dmdhrumilmistry/defect-detect/pkg/config"
	"github.com/golang-jwt/jwt"
)

type ctxKey string

const UserCtxKey ctxKey = "userId"

func CreateJWT(userId string) (string, error) {
	expiration := time.Second * time.Duration(config.DefaultConfig.JWTExpirationInSeconds)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userId":    userId,
		"expiredAt": time.Now().Add(expiration).Unix(),
	})

	tokenStr, err := token.SignedString([]byte(config.DefaultConfig.JWTSecretKey))
	if err != nil {
		return "", err
	}

	return tokenStr, nil
}

func GetTokenFromRequest(r *http.Request) string {
	token := r.Header.Get("Authorization")
	return token
}

func ValidateJWT(tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(config.DefaultConfig.JWTSecretKey), nil
	})
}

func GetUserIdFromContext(ctx context.Context) int {
	userId, ok := ctx.Value(UserCtxKey).(int)
	if !ok {
		return -1
	}

	return userId
}
