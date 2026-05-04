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
	merchantEntity "github.com/dedyf5/resik/entities/merchant"
	outletEntity "github.com/dedyf5/resik/entities/outlet"
	userEntity "github.com/dedyf5/resik/entities/user"
	"github.com/dedyf5/resik/internal/identity"
	"github.com/dedyf5/resik/pkg/collection"
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
	User      User   `json:"user"`
	Merchants []Base `json:"merchants"`
	Outlets   []Base `json:"outlets"`
}

type User struct {
	Base
	Username string `json:"username"`
}

type Base struct {
	ID       uint64         `json:"-"`
	PublicID uuidPkg.UUIDV7 `json:"id"`
}

func (a *AuthClaims) Valid() error {
	return nil
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

func (a *AuthClaims) GetMerchantID(merchantPublicID string) (merchantID uint64, err *resPkg.Status) {
	if a == nil {
		return 0, resPkg.NewStatusCode(http.StatusUnauthorized)
	}

	for _, v := range a.Merchants {
		if v.PublicID.Equal(merchantPublicID) {
			return v.ID, nil
		}
	}

	return 0, resPkg.NewStatusCode(http.StatusUnauthorized)
}

func (a *AuthClaims) MerchantIDIsAccessible(merchantID uint64) (ok bool, err *resPkg.Status) {
	if a == nil {
		return statusUnauthorized()
	}

	merchantIDs := collection.Map(a.Merchants, func(n Base) uint64 {
		return n.ID
	})

	return checkAccess(merchantIDs, merchantID)
}

func (a *AuthClaims) GetOutletID(outletPublicID string) (outletID uint64, err *resPkg.Status) {
	if a == nil {
		return 0, resPkg.NewStatusCode(http.StatusUnauthorized)
	}

	for _, v := range a.Outlets {
		if v.PublicID.Equal(outletPublicID) {
			return v.ID, nil
		}
	}

	return 0, resPkg.NewStatusCode(http.StatusUnauthorized)
}

func (a *AuthClaims) OutletIDIsAccessible(outletID uint64) (ok bool, err *resPkg.Status) {
	if a == nil {
		return statusUnauthorized()
	}

	outletIDs := collection.Map(a.Outlets, func(n Base) uint64 {
		return n.ID
	})

	return checkAccess(outletIDs, outletID)
}

func AuthTokenGenerate(moduleConfig config.Module, authConfig config.Auth, user User, merchants, outlets []Base) (token string, err *resPkg.Status) {
	claims := AuthClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    moduleConfig.Name,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(authConfig.Expires)),
		},
		User:      user,
		Merchants: merchants,
		Outlets:   outlets,
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

		merchantPublicIDs := collection.Map(claims.Merchants, func(n Base) uuidPkg.UUIDV7 {
			return n.PublicID
		})

		if len(merchantPublicIDs) > 0 {
			merchantIDs, errResolver := resolver.ResolveBatch(c, merchantEntity.TABLE_NAME, merchantPublicIDs)
			if errResolver != nil {
				return nil, HTTPStatusError(errResolver, lang)
			}

			for i, id := range merchantIDs {
				claims.Merchants[i].ID = id
			}
		}

		outletPublicIDs := collection.Map(claims.Outlets, func(n Base) uuidPkg.UUIDV7 {
			return n.PublicID
		})

		if len(outletPublicIDs) > 0 {
			outletIDs, errResolver := resolver.ResolveBatch(c, outletEntity.TABLE_NAME, outletPublicIDs)
			if errResolver != nil {
				return nil, HTTPStatusError(errResolver, lang)
			}

			for i, id := range outletIDs {
				claims.Outlets[i].ID = id
			}
		}

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

func statusUnauthorized() (bool, *resPkg.Status) {
	return false, resPkg.NewStatusCode(http.StatusUnauthorized)
}

func checkAccess[T comparable](ids []T, id T) (ok bool, err *resPkg.Status) {
	if !slices.Contains(ids, id) {
		return statusUnauthorized()
	}
	return true, nil
}
