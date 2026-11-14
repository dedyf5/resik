// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package ctx

import (
	"context"
	"fmt"
	"testing"

	langCtx "github.com/dedyf5/resik/ctx/lang"
	logCtx "github.com/dedyf5/resik/ctx/log"
	configEntity "github.com/dedyf5/resik/entities/config"
	permisEntity "github.com/dedyf5/resik/entities/permission"
	"golang.org/x/text/language"
)

func TestAllCtx(t *testing.T) {
	loggerInst := logCtx.Get(configEntity.Log{File: ""}, configEntity.Module{NameKey: "test_module"})

	t.Run("Fail when language is missing from context", func(t *testing.T) {
		c := context.Background()
		c = context.WithValue(c, logCtx.KeyCallerHolderContext, &logCtx.CallerHolder{})

		ctx, errStatus := NewCtx(c, loggerInst)
		if errStatus == nil {
			t.Error("expecated error status due to missing language, got nil")
		}
		if ctx != nil {
			t.Errorf("expecated nil ctx instance on error, got %v", ctx)
		}
	})

	t.Run("Success NewCtx with valid language context", func(t *testing.T) {
		langInst := langCtx.NewLangLocaleDir(language.English, &language.Indonesian, "", fmt.Sprintf("%s%s", "../../", langCtx.LocaleDir))

		c := context.Background()
		c = context.WithValue(c, logCtx.KeyCallerHolderContext, &logCtx.CallerHolder{})
		c = context.WithValue(c, langCtx.ContextKey, langInst)

		ctx, errStatus := NewCtx(c, loggerInst)
		if errStatus != nil {
			t.Fatalf("expected no error, got status: %v", errStatus)
		}
		if ctx == nil {
			t.Fatal("expecated ctx instance, got nil")
		}

		if ctx.Lang() == nil {
			t.Error("expected lang instance, got nil")
		}
		if ctx.Lang() != langInst {
			t.Error("Lang() missmatch")
		}

		if ctx.Log() == nil {
			t.Error("expected log instance, got nil")
		}
		if ctx.Log() != loggerInst {
			t.Error("Log() missmatch")
		}

		_ = ctx.UserClaims()
		_, _, _ = ctx.GetOrganizationID(nil, "org-public-id", permisEntity.PERMISSION_CODE_ORG_READ)
		_, _, _ = ctx.GetBranchID(nil, "branch-public-id", permisEntity.PERMISSION_CODE_BRANCH_READ)
	})
}
