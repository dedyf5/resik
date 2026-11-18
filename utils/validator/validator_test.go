// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package validator

import (
	"errors"
	"net/http"
	"testing"
	"time"

	langCtx "github.com/dedyf5/resik/ctx/lang"
	"golang.org/x/text/language"
)

type FullPayload struct {
	ID          string `json:"id" protobuf:"name=id_proto" validate:"required,uuidv7"`
	Name        string `json:"name" validate:"required,notblank"`
	Email       string `json:"email" validate:"required,email"`
	SortOrder   string `json:"sort_order" validate:"oneof_order=name created_at"`
	BirthDate   string `json:"birth_date" validate:"past_or_present"`
	EmptyDate   string `json:"empty_date" validate:"past_or_present"`
	Bio         string `json:"bio" validate:"min=5,max=10"`
	Age         int    `json:"age" validate:"min=18,max=60,numeric"`
	Code        string `json:"code" validate:"len=4"`
	Other       string `json:"other" validate:"eq=test"`
	CustomField string `validate:"required"`
}

func TestValidator_NewAndTranslators(t *testing.T) {
	v := New(language.Indonesian)
	if v == nil {
		t.Fatal("expected non-nil validator")
	}

	langs := []language.Tag{
		language.Indonesian,
		language.Japanese,
		language.English,
		language.Make("de"),
	}

	for _, l := range langs {
		tr := v.Translator(l)
		if tr == nil {
			t.Errorf("expected translator for language %v", l)
		}
		lt := LanguageToTranslator(l)
		if lt == nil {
			t.Errorf("expected locales translator for language %v", l)
		}
	}
}

func TestValidate_Struct_SuccessAndNil(t *testing.T) {
	v := New(language.Indonesian)

	t.Run("Nil payload returns internal server error status", func(t *testing.T) {
		err := v.Struct(nil, nil)
		if err == nil || err.Code != http.StatusInternalServerError {
			t.Error("expected internal server error for nil struct")
		}
	})

	t.Run("Valid payload passes validation successfully", func(t *testing.T) {
		validUUID := "01a1ecb8-cd44-7d06-982b-12cb796801b6"
		pastDate := time.Now().Add(-48 * time.Hour).Format(time.RFC3339)

		payload := FullPayload{
			ID:          validUUID,
			Name:        "Budi Susanti",
			Email:       "email@budi.com",
			SortOrder:   "name",
			BirthDate:   pastDate,
			EmptyDate:   "",
			Bio:         "Golang",
			Age:         25,
			Code:        "1234",
			Other:       "test",
			CustomField: "value",
		}

		err := v.Struct(payload, nil)
		if err != nil {
			t.Errorf("expected nil error for valid payload, got %v", err)
		}
	})
}

func TestValidate_Struct_ValidationErrorsAndReasons(t *testing.T) {
	v := New(language.English)

	t.Run("Invalid payload triggers comprehensive error reasons", func(t *testing.T) {
		payload := FullPayload{
			ID:          "invalid-uuid",
			Name:        "",
			Email:       "invalid-email",
			SortOrder:   "bad-order",
			BirthDate:   "2050-01-01T00:00:00Z",
			Bio:         "a",     // Terlalu pendek (TOO_SHORT)
			Age:         100,     // Terlalu besar (TOO_LARGE)
			Code:        "12345", // Panjang salah (INVALID_LENGTH)
			Other:       "wrong", // Tidak sama dengan eq=test
			CustomField: "",
		}

		errEn := v.Struct(payload, nil)
		if errEn == nil || errEn.Code != http.StatusBadRequest {
			t.Errorf("expected bad request status for english, got %v", errEn)
		}

		langID := langCtx.NewLang(language.English, &language.Indonesian, "id")
		errID := v.Struct(payload, langID)
		if errID == nil || errID.Code != http.StatusBadRequest {
			t.Errorf("expected bad request status for indonesian, got %v", errID)
		}

		langJA := langCtx.NewLang(language.English, &language.Japanese, "ja")
		errJA := v.Struct(payload, langJA)
		if errJA == nil || errJA.Code != http.StatusBadRequest {
			t.Errorf("expected bad request status for japanese, got %v", errJA)
		}
	})
}

func TestErrorFormatter_NonValidationErrors(t *testing.T) {
	v := New(language.English)
	err := errors.New("generic error")
	status := v.ErrorFormatter(err, nil)
	if status.Code != http.StatusBadRequest {
		t.Errorf("expected status code 400 for non-validation error")
	}
}
