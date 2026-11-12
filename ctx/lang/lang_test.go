// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package lang

import (
	"context"
	"net/http"
	"os"
	"path/filepath"
	"testing"

	"github.com/dedyf5/resik/entities/common"
	"golang.org/x/text/language"
)

func TestGetLanguageReqOrDefault(t *testing.T) {
	defaultLang := language.English
	reqLang := language.Indonesian

	t.Run("Returns reqLang when not nil", func(t *testing.T) {
		got := GetLanguageReqOrDefault(defaultLang, &reqLang)
		if got != reqLang {
			t.Errorf("expected %v, got %v", reqLang, got)
		}
	})

	t.Run("Returns defaultLang when reqLang is nil", func(t *testing.T) {
		got := GetLanguageReqOrDefault(defaultLang, nil)
		if got != defaultLang {
			t.Errorf("expected %v, got %v", defaultLang, got)
		}
	})
}

func TestGetLanguageOrDefault(t *testing.T) {
	t.Run("Valid avaiable language", func(t *testing.T) {
		got := GetLanguageOrDefault("id")
		if id := language.Indonesian; got != id {
			t.Errorf("expected %v, got %v", id, got)
		}
	})

	t.Run("Invalid/Unavailable language falls back to Default", func(t *testing.T) {
		got := GetLanguageOrDefault("invalid-lang-code")
		if got != Default {
			t.Errorf("expected Default (%v), got %v", Default, got)
		}
	})
}

func TestLanguageIsAvailable(t *testing.T) {
	t.Run("Empty language returns REQUIRED error", func(t *testing.T) {
		ok, status := LanguageIsAvailable("")
		if ok || status == nil || status.Code != http.StatusBadRequest {
			t.Errorf("expected false and 400 Bad Request, got ok=%v, status=%v", ok, status)
		}
	})

	t.Run("Unavailable language returns OUT_OF_RANGE error", func(t *testing.T) {
		ok, status := LanguageIsAvailable("fr")
		if ok || status == nil || status.Code != http.StatusBadRequest {
			t.Errorf("expected false and 400 Bad Request for unsupported language, got ok=%v, status=%v", ok, status)
		}
	})

	t.Run("Available language returns true", func(t *testing.T) {
		ok, status := LanguageIsAvailable("id")
		if !ok || status != nil {
			t.Errorf("expected true and nil status, got ok=%v, status=%v", ok, status)
		}
	})
}

func TestGetLanguageAvailable(t *testing.T) {
	t.Run("Valid language tag", func(t *testing.T) {
		tag, status := GetLanguageAvailable("ja")
		if status != nil || tag == nil || *tag != language.Japanese {
			t.Errorf("expected Japanese tag, got tag=%v, status=%v", tag, status)
		}
	})

	t.Run("Invalid language returns error status", func(t *testing.T) {
		tag, status := GetLanguageAvailable("unknown")
		if status == nil || tag != nil {
			t.Errorf("expected nil tag and status error, got tag=%v, status=%v", tag, status)
		}
	})
}

func TestLangKey_String(t *testing.T) {
	if ContextKey.String() != "lang" {
		t.Errorf("expected 'lang', got %q", ContextKey.String())
	}
}

func TestFromContext(t *testing.T) {
	t.Run("Success retrieving Lang from context", func(t *testing.T) {
		l := &Lang{}
		c := context.WithValue(context.Background(), ContextKey, l)
		got, status := FromContext(c)
		if status != nil || got != l {
			t.Errorf("expected Lang=%v and status=nil, got=%v, status=%v", l, got, status)
		}
	})

	t.Run("Failure casting Lang from context", func(t *testing.T) {
		c := context.WithValue(context.Background(), ContextKey, "not-a-lang")
		got, status := FromContext(c)
		if status == nil || got != nil || status.Code != http.StatusInternalServerError {
			t.Errorf("expected 500 Internal Server Error status, got=%v, status=%v", got, status)
		}
	})
}

func TestLangMethodsWithTempDir(t *testing.T) {
	tmpDir := t.TempDir()

	tomlContent := `
	[validation.field.name]
	other = "Name"

	[validation.quote_val]
	other = "\"{{.Val}}\""
	
	[test_message]
	other = "Hello {{.Val}}"
	`
	for _, v := range Available {
		filePath := filepath.Join(tmpDir, "active."+v.String()+".toml")
		if err := os.WriteFile(filePath, []byte(tomlContent), 0600); err != nil {
			t.Fatalf("failed to write dummy toml file: %v", err)
		}
	}

	reqLang := language.Indonesian
	langInst := NewLangLocaleDir(Default, &reqLang, "en", tmpDir)

	t.Run("NewLang default constructor", func(t *testing.T) {
		l := NewLang(Default, &reqLang, "id")
		if l == nil {
			t.Error("expected Lang instance, got nil")
		}
	})

	t.Run("LanguageReqOrDefault", func(t *testing.T) {
		if got := langInst.LanguageReqOrDefault(); got != reqLang {
			t.Errorf("expected %v, got %v", reqLang, got)
		}
	})

	t.Run("Validation field", func(t *testing.T) {
		fieldKey := langInst.ValidationField("name")
		if expected := "validation.field.name"; fieldKey != expected {
			t.Errorf("expected %s, got %s", expected, fieldKey)
		}

		name := langInst.GetValidationFieldName("name")
		if name == "" {
			t.Error("expected localized field name, got empty")
		}

		nameWithQuote := langInst.GetValidationFieldNameWithQuote("name")
		if expected := "\"Name\""; nameWithQuote != expected {
			t.Errorf("expected %s with quote, got %s", expected, nameWithQuote)
		}
	})

	t.Run("GetByMessageID and GetByTemplateData", func(t *testing.T) {
		msg := langInst.GetByMessageID("test_message")
		if msg == "" {
			t.Error("expected localized test_message, got empty")
		}

		templateMsg := langInst.GetByTemplateData("test_message", common.Map{"Val": "World"})
		if expected := "Hello World"; templateMsg != expected {
			t.Errorf("expected %q, got %q", expected, templateMsg)
		}
	})
}
