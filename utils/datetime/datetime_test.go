// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package datetime

import (
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/dedyf5/resik/config"
	"github.com/dedyf5/resik/ctx"
	langCtx "github.com/dedyf5/resik/ctx/lang"
	"github.com/dedyf5/resik/ctx/log"
	configEntity "github.com/dedyf5/resik/entities/config"
	"golang.org/x/text/language"
)

func TestFromString(t *testing.T) {
	ctx := env()

	t.Run("Empty string returns nil and no error", func(t *testing.T) {
		res, err := FromString("", time.RFC3339, ctx)
		if res != nil {
			t.Errorf("expected nil time pointer, got %v", res)
		}
		if err != nil {
			t.Errorf("expected no status error, got %v", err)
		}
	})

	t.Run("Valid datetime string returns parsed time", func(t *testing.T) {
		timeStr := "2026-11-15T10:00:00Z"
		res, err := FromString(timeStr, time.RFC3339, ctx)
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		if res == nil {
			t.Fatal("expected non-nil time pointer, got nil")
		}
		if res.Year() != 2026 || res.Day() != 15 {
			t.Errorf("parsed time mismatch: got %v", res)
		}
	})

	t.Run("Invalid datetime string returns status error", func(t *testing.T) {
		invalidStr := "not-a-datetime"
		res, err := FromString(invalidStr, time.RFC3339, ctx)
		if err == nil {
			t.Error("expected non-nil error status, got nil")
		} else if expected := http.StatusInternalServerError; err.Code != expected {
			t.Errorf("expected error status code %d, got %d", expected, err.Code)
		}
		if res != nil {
			t.Errorf("expected nil time pointer on failure, got %v", res)
		}
	})
}

func env() *ctx.Ctx {
	conf := config.Config{
		Module: configEntity.Module{
			Name:        "REST",
			NameKey:     "rest",
			LangDefault: language.English,
			Host:        "0.0.0.0",
			Port:        8081,
		},
	}
	c := context.WithValue(
		context.Background(),
		langCtx.ContextKey,
		langCtx.NewLangLocaleDir(
			language.English,
			&language.English,
			"",
			fmt.Sprintf("%s%s", "../../", langCtx.LocaleDir),
		),
	)
	ctx, _ := ctx.NewCtx(c, log.Get(configEntity.Log{}, conf.Module))

	return ctx
}
