// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package jwt

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/dedyf5/resik/ctx/lang"
	"github.com/dedyf5/resik/ctx/lang/term"
	"github.com/dedyf5/resik/entities/config"
	merchantEntity "github.com/dedyf5/resik/entities/merchant"
	outletEntity "github.com/dedyf5/resik/entities/outlet"
	userEntity "github.com/dedyf5/resik/entities/user"
	"github.com/dedyf5/resik/internal/identity"
	resPkg "github.com/dedyf5/resik/pkg/response"
	uuidPkg "github.com/dedyf5/resik/pkg/uuid"
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
	User        User     `json:"user"`
	MerchantIDs []uint64 `json:"-"`
	OutletsIDs  []uint64 `json:"-"`
}

type User struct {
	Base
	Username string `json:"username"`
}

type Base struct {
	ID       uint64         `json:"-"`
	PublicID uuidPkg.UUIDV7 `json:"id"`
}

func (a *AuthClaims) UserID() uint64 {
	if a == nil {
		return 0
	}
	return a.User.ID
}

func (a *AuthClaims) UserPublicID() uuidPkg.UUIDV7 {
	if a == nil {
		return uuidPkg.UUIDV7{}
	}
	return a.User.PublicID
}

func (a *AuthClaims) Username() string {
	if a == nil {
		return ""
	}
	return a.User.Username
}

// getID gets the ID for a specific table by public ID, and check if user has access to it for the given permission code
func (a *AuthClaims) getID(c context.Context, resolver identity.IdentityResolver, lang *lang.Lang, permissionCode string, tableName string, publicID string) (id uint64, err *resPkg.Status) {
	if a == nil {
		return 0, resPkg.NewStatusCode(http.StatusUnauthorized)
	}

	uuidV7, errParse := uuidPkg.ParseUUIDV7(publicID)
	if errParse != nil {
		return 0, resPkg.NewStatusMessage(
			http.StatusBadRequest,
			term.InvalidID.Localize(lang.Localizer),
			errParse,
		)
	}

	hasAccess, errResolver := resolver.HasAccessByPublicID(c, a.User.ID, permissionCode, tableName, uuidV7)
	if errResolver != nil {
		return 0, HTTPStatusError(errResolver, lang)
	}

	if hasAccess {
		id, errResolve := resolver.Resolve(c, tableName, uuidV7)
		if errResolve != nil {
			return 0, HTTPStatusError(errResolve, lang)
		}
		return id, nil
	}

	return 0, resPkg.NewStatusCode(http.StatusUnauthorized)
}

// GetMerchantID gets the merchant ID by merchant public ID, and check if user has access to it for permission code
func (a *AuthClaims) GetMerchantID(c context.Context, resolver identity.IdentityResolver, lang *lang.Lang, merchantPublicID string, permissionCode string) (merchantID uint64, err *resPkg.Status) {
	return a.getID(c, resolver, lang, permissionCode, merchantEntity.TABLE_NAME, merchantPublicID)
}

// GetOutletID gets the outlet ID by outlet public ID, and check if user has access to it for permission code
func (a *AuthClaims) GetOutletID(c context.Context, resolver identity.IdentityResolver, lang *lang.Lang, outletPublicID string, permissionCode string) (outletID uint64, err *resPkg.Status) {
	return a.getID(c, resolver, lang, permissionCode, outletEntity.TABLE_NAME, outletPublicID)
}

func AuthTokenGenerate(moduleConfig config.Module, authConfig config.Auth, user User, merchantIDs, outletIDs []uint64) (token string, err *resPkg.Status) {
	claims := AuthClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    moduleConfig.Name,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(authConfig.Expires)),
		},
		User:        user,
		MerchantIDs: merchantIDs,
		OutletsIDs:  outletIDs,
	}

	tokenGen := jwt.NewWithClaims(AUTH_SIGNING_METHOD, claims)
	token, errToken := tokenGen.SignedString([]byte(authConfig.SignatureKey))
	if errToken != nil {
		return "", resPkg.NewStatusError(http.StatusInternalServerError, errToken)
	}
	return
}

func AuthClaimsFromString(tokenString string, signatureKey string, c context.Context, resolver identity.IdentityResolver, lang *lang.Lang) (claim *AuthClaims, err *resPkg.Status) {
	if tokenString == "" {
		return nil, statusInvalid(errors.New("missing value in request header"), lang)
	}

	token, errParse := jwt.ParseWithClaims(tokenString, &AuthClaims{}, func(t *jwt.Token) (any, error) {
		return []byte(signatureKey), nil
	})
	if errParse != nil {
		return nil, statusInvalid(errParse, lang)
	}

	if claims, ok := token.Claims.(*AuthClaims); ok {
		if !claims.User.PublicID.IsEmpty() {
			userID, errResolver := resolver.Resolve(c, userEntity.TABLE_NAME, claims.User.PublicID)
			if errResolver != nil {
				return nil, HTTPStatusError(errResolver, lang)
			}
			claims.User.ID = userID
		}

		merchantIDs, errMerchants := resolver.GetTenantMerchantIDs(c, claims.User.ID)
		if errMerchants != nil {
			return nil, HTTPStatusError(errMerchants, lang)
		}
		claims.MerchantIDs = merchantIDs

		outletIDs, errOutlets := resolver.GetTenantOutletIDs(c, claims.User.ID)
		if errOutlets != nil {
			return nil, HTTPStatusError(errOutlets, lang)
		}
		claims.OutletsIDs = outletIDs

		return claims, nil
	}

	return nil, statusInvalid(errors.New("error while casting AuthClaims"), lang)
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
		return statusInvalid(err, lang)
	}

	return resPkg.NewStatusError(
		http.StatusInternalServerError,
		err,
	)
}

func statusInvalid(err error, lang *lang.Lang) *resPkg.Status {
	return resPkg.NewStatusMessage(
		http.StatusUnauthorized,
		term.InvalidOrExpiredSessionLoginAgain.Localize(lang.Localizer),
		err,
	)
}
