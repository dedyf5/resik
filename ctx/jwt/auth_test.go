// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package jwt

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"testing"
	"time"

	langCtx "github.com/dedyf5/resik/ctx/lang"
	branchEntity "github.com/dedyf5/resik/entities/branch"
	configEntity "github.com/dedyf5/resik/entities/config"
	orgEntity "github.com/dedyf5/resik/entities/organization"
	permisEntity "github.com/dedyf5/resik/entities/permission"
	userEntity "github.com/dedyf5/resik/entities/user"
	mockIdentity "github.com/dedyf5/resik/internal/identity/mock"
	responsePkg "github.com/dedyf5/resik/pkg/response"
	uuidPkg "github.com/dedyf5/resik/pkg/uuid"
	"go.uber.org/mock/gomock"
	"golang.org/x/text/language"
)

func getTestLang(t *testing.T) *langCtx.Lang {
	t.Helper()
	reqLang := language.English
	return langCtx.NewLangLocaleDir(
		langCtx.Default, &reqLang, "",
		fmt.Sprintf("%s%s", "../../", langCtx.LocaleDir),
	)
}

func TestAuthTokenGenerateAndClaimsFromString(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	modConfig := configEntity.Module{Name: "TestIssuer"}
	authConfig := configEntity.Auth{
		Expires:      time.Hour,
		SignatureKey: "test-signature-secret-key-very-secure",
	}

	pubID := uuidPkg.NewUUIDV7()
	user := User{
		ID:       1,
		PublicID: pubID,
		Username: "budi",
	}

	token, errStatus := AuthTokenGenerate(modConfig, authConfig, user)
	if errStatus != nil {
		t.Fatalf("expected no error generating token, got %v", errStatus)
	}
	if token == "" {
		t.Fatal("expected generated token string to be non-empty")
	}

	resolver := mockIdentity.NewMockIdentityResolver(ctrl)

	resolver.EXPECT().Resolve(gomock.Any(), userEntity.TABLE_NAME, pubID).Return(uint64(1), nil)
	resolver.EXPECT().GetTenantOrganizationIDs(gomock.Any(), uint64(1)).Return([]uint64{10, 20}, nil)
	resolver.EXPECT().GetTenantBranchIDs(gomock.Any(), uint64(1)).Return([]uint64{100, 200}, nil)

	c := context.Background()
	l := getTestLang(t)

	claims, errStatus := AuthClaimsFromString(token, authConfig.SignatureKey, c, resolver, l)
	if errStatus != nil {
		t.Fatalf("expected no error parsing token, got %v", errStatus)
	}
	if claims == nil {
		t.Fatalf("expected claims to be non-nil")
	}

	if claims.Username() != "budi" {
		t.Fatalf("expected claims.Username() = %q, got %q", "budi", claims.Username())
	}
	if claims.UserID() != 1 {
		t.Fatalf("expected claims.UserID() = %d, got %d", 1, claims.UserID())
	}
	if claims.UserPublicID() != pubID {
		t.Fatalf("expected claims.UserPublicID() = %v, got %v", pubID, claims.UserPublicID())
	}
}

func TestAuthClaimsHelperOnNil(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	var claims *AuthClaims

	if id := claims.UserID(); id != 0 {
		t.Errorf("expected claims.UserID() = %d, got %d", 0, id)
	}
	if pubID := claims.UserPublicID(); !pubID.IsEmpty() {
		t.Errorf("expected empty claims.UserPublicID(), got %v", pubID)
	}
	if username := claims.Username(); username != "" {
		t.Errorf("expected empty claims.Username(), got %q", username)
	}

	c := context.Background()
	resolver := mockIdentity.NewMockIdentityResolver(ctrl)
	l := getTestLang(t)
	orgPubID := uuidPkg.NewUUIDV7().String()

	_, _, errStatus := claims.GetOrganizationID(c, resolver, l, orgPubID, permisEntity.PERMISSION_CODE_ORG_READ)
	if errStatus == nil || errStatus.Code != http.StatusUnauthorized {
		t.Errorf("expected 401 Unauthorized for nil claims GetOrganizationID, got %v", errStatus)
	}
}

func TestAuthClaimsGetOrganizationAndBranchID(t *testing.T) {
	claims := &AuthClaims{
		User: User{Base: Base{ID: 1}},
	}
	c := context.Background()
	validPubIDStr := uuidPkg.NewUUIDV7().String()
	parsedUUID, _ := uuidPkg.ParseUUIDV7(validPubIDStr)

	t.Run("GetOrganizationID Success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		resolver := mockIdentity.NewMockIdentityResolver(ctrl)
		l := getTestLang(t)

		resolver.EXPECT().
			HasAccessByPublicID(gomock.Any(), l, uint64(1), permisEntity.PERMISSION_CODE_ORG_READ, orgEntity.TABLE_NAME, parsedUUID).
			Return(true, nil)
		resolver.EXPECT().
			Resolve(gomock.Any(), orgEntity.TABLE_NAME, parsedUUID).
			Return(uint64(50), nil)

		orgID, pubID, errStatus := claims.GetOrganizationID(c, resolver, l, validPubIDStr, permisEntity.PERMISSION_CODE_ORG_READ)
		if errStatus != nil {
			t.Errorf("expected no error, got %v", errStatus)
		}
		if orgID != 50 {
			t.Errorf("expected organization id = %d, got %d", 50, orgID)
		}
		if pubID != parsedUUID {
			t.Errorf("expected organization public id = %v, got %v", parsedUUID, pubID)
		}
	})

	t.Run("GetBranchID Success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		resolver := mockIdentity.NewMockIdentityResolver(ctrl)
		l := getTestLang(t)

		resolver.EXPECT().
			HasAccessByPublicID(gomock.Any(), l, uint64(1), permisEntity.PERMISSION_CODE_BRANCH_READ, branchEntity.TABLE_NAME, parsedUUID).
			Return(true, nil)
		resolver.EXPECT().
			Resolve(gomock.Any(), branchEntity.TABLE_NAME, parsedUUID).
			Return(uint64(80), nil)

		branchID, _, errStatus := claims.GetBranchID(c, resolver, l, validPubIDStr, permisEntity.PERMISSION_CODE_BRANCH_READ)
		if errStatus != nil {
			t.Errorf("expected no error, got %v", errStatus)
		}
		if branchID != 80 {
			t.Errorf("expected branch id = %d, got %d", 80, branchID)
		}
	})

	t.Run("Invalid UUID format error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		resolver := mockIdentity.NewMockIdentityResolver(ctrl)
		l := getTestLang(t)

		_, _, errStatus := claims.GetOrganizationID(c, resolver, l, "invalid-uuid-format", permisEntity.PERMISSION_CODE_ORG_READ)
		if errStatus == nil || errStatus.Code != http.StatusBadRequest {
			t.Errorf("expected 400 Bad Request for invalid UUID format, got %v", errStatus)
		}
	})
}

func TestAuthClaimsGetID_Errors(t *testing.T) {
	claims := &AuthClaims{
		User: User{Base: Base{ID: 1}},
	}
	c := context.Background()
	validPubIDStr := uuidPkg.NewUUIDV7().String()
	parsedUUID, _ := uuidPkg.ParseUUIDV7(validPubIDStr)
	l := getTestLang(t)

	t.Run("HasAccessByPublicID Error (errResolver)", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		resolver := mockIdentity.NewMockIdentityResolver(ctrl)
		customErr := responsePkg.NewStatusCode(http.StatusInternalServerError)

		resolver.EXPECT().
			HasAccessByPublicID(gomock.Any(), l, uint64(1), permisEntity.PERMISSION_CODE_ORG_READ, orgEntity.TABLE_NAME, parsedUUID).
			Return(false, customErr)

		_, _, errStatus := claims.GetOrganizationID(c, resolver, l, validPubIDStr, permisEntity.PERMISSION_CODE_ORG_READ)
		if errStatus != customErr {
			t.Errorf("expected error status = %v, got %v", customErr, errStatus)
		}
	})

	t.Run("Resolve Error (errResolve)", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		resolver := mockIdentity.NewMockIdentityResolver(ctrl)
		customErr := responsePkg.NewStatusCode(http.StatusNotFound)

		resolver.EXPECT().
			HasAccessByPublicID(gomock.Any(), l, uint64(1), permisEntity.PERMISSION_CODE_ORG_READ, orgEntity.TABLE_NAME, parsedUUID).
			Return(true, nil)
		resolver.EXPECT().
			Resolve(gomock.Any(), orgEntity.TABLE_NAME, parsedUUID).
			Return(uint64(0), customErr)

		_, _, errStatus := claims.GetOrganizationID(c, resolver, l, validPubIDStr, permisEntity.PERMISSION_CODE_ORG_READ)
		if errStatus != customErr {
			t.Errorf("expected error status = %v, got %v", customErr, errStatus)
		}
	})

	t.Run("No Access Unauthorized", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		resolver := mockIdentity.NewMockIdentityResolver(ctrl)

		resolver.EXPECT().
			HasAccessByPublicID(gomock.Any(), l, uint64(1), permisEntity.PERMISSION_CODE_ORG_READ, orgEntity.TABLE_NAME, parsedUUID).
			Return(false, nil)

		_, _, errStatus := claims.GetOrganizationID(c, resolver, l, validPubIDStr, permisEntity.PERMISSION_CODE_ORG_READ)
		if errStatus == nil || errStatus.Code != http.StatusUnauthorized {
			t.Errorf("expected 401 Unauthorized when access is false, got %v", errStatus)
		}
	})
}

func TestAuthClaimsFromContext(t *testing.T) {
	claims := &AuthClaims{User: User{Username: "ctxuser"}}
	c := context.WithValue(context.Background(), AuthClaimsKey, claims)

	retrieved := AuthClaimsFromContext(c)
	if retrieved == nil || retrieved.Username() != "ctxuser" {
		t.Error("failed to retrieve claims from context")
	}

	emptyC := context.Background()
	if AuthClaimsFromContext(emptyC) != nil {
		t.Error("expected nil for empty context")
	}

	invalidC := context.WithValue(context.Background(), AuthClaimsKey, "not-claims-type")
	if AuthClaimsFromContext(invalidC) != nil {
		t.Error("expected nil for invalid type in context")
	}
}

func TestAuthClaimsFromStringError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	c := context.Background()
	resolver := mockIdentity.NewMockIdentityResolver(ctrl)
	l := getTestLang(t)

	t.Run("Empty token string", func(t *testing.T) {
		_, errStatus := AuthClaimsFromString("", "secret", c, resolver, l)
		if errStatus == nil || errStatus.Code != http.StatusUnauthorized {
			t.Errorf("expected 401 Unauthorized for empty token, got %v", errStatus)
		}
	})

	t.Run("Invalid token signature or format", func(t *testing.T) {
		_, errStatus := AuthClaimsFromString("invalid-jwt-token", "secret", c, resolver, l)
		if errStatus == nil || errStatus.Code != http.StatusUnauthorized {
			t.Errorf("expected 401 Unauthorized for invalid token, got %v", errStatus)
		}
	})
}

func TestAuthClaimsFromString_InternalErrors(t *testing.T) {
	modConfig := configEntity.Module{Name: "TestIssuer"}
	authConfig := configEntity.Auth{
		Expires:      time.Hour,
		SignatureKey: "secret",
	}
	pubID := uuidPkg.NewUUIDV7()
	user := User{ID: 1, PublicID: pubID, Username: "Budi"}

	token, _ := AuthTokenGenerate(modConfig, authConfig, user)
	c := context.Background()
	l := getTestLang(t)

	customErr := responsePkg.NewStatusCode(http.StatusInternalServerError)

	t.Run("Resolve User Error (errResolver)", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		resolver := mockIdentity.NewMockIdentityResolver(ctrl)

		resolver.EXPECT().
			Resolve(c, userEntity.TABLE_NAME, pubID).
			Return(uint64(0), customErr)

		_, errStatus := AuthClaimsFromString(token, authConfig.SignatureKey, c, resolver, l)
		if errStatus != customErr {
			t.Errorf("expected error status = %v, got %v", customErr, errStatus)
		}
	})

	t.Run("GetTenantOrganizationIDs Error (errOrganizations)", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		resolver := mockIdentity.NewMockIdentityResolver(ctrl)

		resolver.EXPECT().
			Resolve(c, userEntity.TABLE_NAME, pubID).
			Return(uint64(1), nil)
		resolver.EXPECT().
			GetTenantOrganizationIDs(gomock.Any(), uint64(1)).
			Return(nil, customErr)

		_, errStatus := AuthClaimsFromString(token, authConfig.SignatureKey, c, resolver, l)
		if errStatus != customErr {
			t.Errorf("expected error status = %v, got %v", customErr, errStatus)
		}
	})

	t.Run("GetTenantBranchIDs Error (errBranches)", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		resolver := mockIdentity.NewMockIdentityResolver(ctrl)

		resolver.EXPECT().
			Resolve(c, userEntity.TABLE_NAME, pubID).
			Return(uint64(1), nil)
		resolver.EXPECT().
			GetTenantOrganizationIDs(gomock.Any(), uint64(1)).
			Return([]uint64{10}, nil)
		resolver.EXPECT().
			GetTenantBranchIDs(gomock.Any(), uint64(1)).
			Return(nil, customErr)

		_, errStatus := AuthClaimsFromString(token, authConfig.SignatureKey, c, resolver, l)
		if errStatus != customErr {
			t.Errorf("expected error status = %v, got %v", customErr, errStatus)
		}
	})
}

func TestHTTPStatusError(t *testing.T) {
	l := getTestLang(t)

	t.Run("Error containing invalid keyword", func(t *testing.T) {
		err := errors.New("invalid request parameter")
		status := HTTPStatusError(err, l)
		if status.Code != http.StatusUnauthorized {
			t.Errorf("expected 401 Unauthorized for error with 'invalid', got %d", status.Code)
		}
	})

	t.Run("General internal error", func(t *testing.T) {
		err := errors.New("database connection timeout")
		status := HTTPStatusError(err, l)
		if status.Code != http.StatusInternalServerError {
			t.Errorf("expected 500 Internal Server Error for general error, got %d", status.Code)
		}
	})
}
