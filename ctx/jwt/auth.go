// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package jwt

import (
	"context"
	"net/http"
	"time"

	"github.com/dedyf5/resik/entities/config"
	resPkg "github.com/dedyf5/resik/pkg/response"
	"github.com/golang-jwt/jwt/v5"
)

var (
	AUTH_SIGNING_METHOD = jwt.SigningMethodHS256
)

type authClaimsKey string

const (
	AuthClaimsKey authClaimsKey = "auth_claims"
)

type AuthClaims struct {
	jwt.RegisteredClaims
	UserID      int64   `json:"user_id"`
	Username    string  `json:"username"`
	MerchantIDs []int64 `json:"merchant_ids"`
}

func (a AuthClaims) Valid() error {
	return nil
}

func AuthTokenGenerate(appConfig config.App, authConfig config.Auth, userID int64, username string, merchantIDs []int64) (token string, status *resPkg.Status) {
	duration := time.Duration(authConfig.Expires) * time.Second
	claims := AuthClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    appConfig.Name,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(duration)),
		},
		UserID:      userID,
		Username:    username,
		MerchantIDs: merchantIDs,
	}
	tokenGen := jwt.NewWithClaims(AUTH_SIGNING_METHOD, claims)
	token, err := tokenGen.SignedString([]byte(authConfig.SignatureKey))
	if err != nil {
		return "", &resPkg.Status{
			Code:       http.StatusInternalServerError,
			CauseError: err,
		}
	}
	return
}

func AuthClaimsFromString(tokenString string, signatureKey string) *AuthClaims {
	token, err := jwt.ParseWithClaims(tokenString, &AuthClaims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(signatureKey), nil
	})
	if err != nil {
		return nil
	}
	if claims, ok := token.Claims.(*AuthClaims); ok {
		return claims
	}
	return nil
}

func AuthClaimsFromContext(ctx context.Context) *AuthClaims {
	value := ctx.Value(AuthClaimsKey)
	if value == nil {
		return nil
	}
	if claims, ok := value.(*AuthClaims); ok {
		return claims
	}
	return nil
}
