// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package jwt

import (
	"context"
	"net/http"
	"time"

	"github.com/dedyf5/resik/entities/config"
	"github.com/dedyf5/resik/pkg/array"
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
	UserID      uint64   `json:"user_id"`
	Username    string   `json:"username"`
	MerchantIDs []uint64 `json:"merchant_ids"`
	OutletIDs   []uint64 `json:"outlet_ids"`
}

func (a *AuthClaims) Valid() error {
	return nil
}

func (a *AuthClaims) MerchantIDIsAccessible(merchantID uint64) (ok bool, status *resPkg.Status) {
	if a == nil {
		return a.statusUnauthorized()
	}
	if array.InArray(merchantID, a.MerchantIDs) < 0 {
		return a.statusUnauthorized()
	}
	return true, nil
}

func (a *AuthClaims) OutletIDIsAccessible(outletID uint64) (ok bool, status *resPkg.Status) {
	if a == nil {
		return a.statusUnauthorized()
	}
	if array.InArray(outletID, a.OutletIDs) < 0 {
		return a.statusUnauthorized()
	}
	return true, nil
}

func (a *AuthClaims) statusUnauthorized() (bool, *resPkg.Status) {
	return false, &resPkg.Status{
		Code: http.StatusUnauthorized,
	}
}

func AuthTokenGenerate(appConfig config.App, authConfig config.Auth, userID uint64, username string, merchantIDs []uint64, outletIDs []uint64) (token string, status *resPkg.Status) {
	duration := time.Duration(authConfig.Expires) * time.Second
	claims := AuthClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    appConfig.Name,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(duration)),
		},
		UserID:      userID,
		Username:    username,
		MerchantIDs: merchantIDs,
		OutletIDs:   outletIDs,
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

func AuthClaimsFromString(tokenString string, signatureKey string) (claim *AuthClaims, status *resPkg.Status) {
	token, err := jwt.ParseWithClaims(tokenString, &AuthClaims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(signatureKey), nil
	})
	if err != nil {
		return nil, &resPkg.Status{
			Code:       http.StatusBadRequest,
			Message:    "invalid jwt",
			CauseError: err,
		}
	}
	if claims, ok := token.Claims.(*AuthClaims); ok {
		return claims, nil
	}
	return nil, &resPkg.Status{
		Code:    http.StatusBadRequest,
		Message: "invalid or expired jwt",
	}
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
