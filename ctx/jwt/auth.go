// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package jwt

import (
	"context"
	"errors"
	"net/http"
	"slices"
	"strings"
	"time"

	"github.com/dedyf5/resik/ctx/lang"
	"github.com/dedyf5/resik/ctx/lang/term"
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
	UserID      uint64   `json:"user_id"`
	Username    string   `json:"username"`
	MerchantIDs []uint64 `json:"merchant_ids"`
	OutletIDs   []uint64 `json:"outlet_ids"`
}

func (a *AuthClaims) Valid() error {
	return nil
}

func (a *AuthClaims) MerchantIDIsAccessible(merchantID uint64) (ok bool, err *resPkg.Status) {
	if a == nil {
		return a.statusUnauthorized()
	}
	return a.checkAccess(a.MerchantIDs, merchantID)
}

func (a *AuthClaims) OutletIDIsAccessible(outletID uint64) (ok bool, err *resPkg.Status) {
	if a == nil {
		return a.statusUnauthorized()
	}
	return a.checkAccess(a.OutletIDs, outletID)
}

func (a *AuthClaims) checkAccess(ids []uint64, id uint64) (ok bool, err *resPkg.Status) {
	if !slices.Contains(ids, id) {
		return a.statusUnauthorized()
	}
	return true, nil
}

func (a *AuthClaims) statusUnauthorized() (bool, *resPkg.Status) {
	return false, resPkg.NewStatusCode(http.StatusUnauthorized)
}

func AuthTokenGenerate(moduleConfig config.Module, authConfig config.Auth, userID uint64, username string, merchantIDs []uint64, outletIDs []uint64) (token string, err *resPkg.Status) {
	claims := AuthClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    moduleConfig.Name,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(authConfig.Expires)),
		},
		UserID:      userID,
		Username:    username,
		MerchantIDs: merchantIDs,
		OutletIDs:   outletIDs,
	}
	tokenGen := jwt.NewWithClaims(AUTH_SIGNING_METHOD, claims)
	token, errToken := tokenGen.SignedString([]byte(authConfig.SignatureKey))
	if errToken != nil {
		return "", resPkg.NewStatusError(http.StatusInternalServerError, errToken)
	}
	return
}

func AuthClaimsFromString(tokenString string, signatureKey string, lang *lang.Lang) (claim *AuthClaims, err *resPkg.Status) {
	if tokenString == "" {
		return nil, resPkg.NewStatusMessage(
			http.StatusUnauthorized,
			term.Unauthorized.Localize(lang.Localizer),
			errors.New("missing value in request header"),
		)
	}
	token, errParse := jwt.ParseWithClaims(tokenString, &AuthClaims{}, func(t *jwt.Token) (any, error) {
		return []byte(signatureKey), nil
	})
	if errParse != nil {
		return nil, resPkg.NewStatusMessage(
			http.StatusUnauthorized,
			term.InvalidOrExpiredSessionLoginAgain.Localize(lang.Localizer),
			errParse,
		)
	}
	if claims, ok := token.Claims.(*AuthClaims); ok {
		return claims, nil
	}
	return nil, resPkg.NewStatusMessage(
		http.StatusUnauthorized,
		term.InvalidOrExpiredSessionLoginAgain.Localize(lang.Localizer),
		errors.New("error while casting AuthClaims"),
	)
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

func HTTPStatusError(err error, lang *lang.Lang) *resPkg.Status {
	if strings.Contains(err.Error(), "invalid") {
		return resPkg.NewStatusMessage(
			http.StatusUnauthorized,
			term.InvalidOrExpiredSessionLoginAgain.Localize(lang.Localizer),
			err,
		)
	}
	return resPkg.NewStatusMessage(
		http.StatusUnauthorized,
		term.Unauthorized.Localize(lang.Localizer),
		err,
	)
}
