// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package response

import (
	"net/http"
	"testing"

	"github.com/pkg/errors"
	"google.golang.org/grpc/codes"
)

func TestStatus_Constructors(t *testing.T) {
	t.Run("NewStatusCode", func(t *testing.T) {
		st := NewStatusCode(http.StatusOK)
		if st.Code != http.StatusOK {
			t.Errorf("Code = %v, want %v", st.Code, http.StatusOK)
		}
		if st.Caller == "" {
			t.Error("Expected Caller to be populated")
		}
	})

	t.Run("NewStatusMessage", func(t *testing.T) {
		err := errors.New("original error")
		customError := "Custom Error"
		st := NewStatusMessage(http.StatusInternalServerError, customError, err)
		if st.Code != http.StatusInternalServerError {
			t.Errorf("Code = %v, want %v", st.Code, http.StatusInternalServerError)
		}
		if st.Message != customError {
			t.Errorf("Message = %v, want '%v'", st.Message, customError)
		}
		if st.CauseError == nil {
			t.Error("Expected CauseError to be set")
		}
		if len(st.StackTrace) == 0 {
			t.Error("Expecated CauseError to be recorded")
		}
	})

	t.Run("NewStatusData", func(t *testing.T) {
		data := map[string]string{"key": "val"}
		st := NewStatusData(http.StatusOK, data)
		if st.Code != http.StatusOK {
			t.Errorf("Code = %v, want %v", st.Code, http.StatusOK)
		}
		if st.Data == nil {
			t.Error("Expected Data to be set")
		}
	})

	t.Run("NewStatusMessageData with error", func(t *testing.T) {
		err := errors.New("some error")
		someData := "some-data"
		st := NewStatusMessageData(http.StatusBadRequest, "Bad", someData, err)
		if st.Data != someData {
			t.Errorf("Data = %v, want '%v'", st.Data, someData)
		}
		if st.CauseError == nil {
			t.Error("Expected CauseError to be set with error is provided")
		}
	})

	t.Run("NewStatusDataMeta", func(t *testing.T) {
		meta := &Meta{Total: 10, PageCurrent: 1, Limit: 5}
		st := NewStatusDataMeta(http.StatusOK, "data", meta)
		if st.Meta != meta {
			t.Errorf("Meta = %+v, want %+v", st.Meta, meta)
		}
	})

	t.Run("NewStatusSuccess", func(t *testing.T) {
		st := NewStatusSuccess(http.StatusOK, "Success", "data")
		if st.Message != "Success" {
			t.Errorf("Message = %v, want '%v'", st.Message, "Success")
		}
	})

	t.Run("NewStatusError", func(t *testing.T) {
		err := errors.New("internal error")
		st := NewStatusError(http.StatusInternalServerError, err)
		if st.Code != http.StatusInternalServerError {
			t.Errorf("Code = %v, want %v", st.Code, http.StatusInternalServerError)
		}
		if st.CauseError == nil {
			t.Error("Expected CauseError to be set")
		}
	})

	t.Run("NewStatusBadRequest and BadRequests extraction", func(t *testing.T) {
		err := errors.New("invalid format")
		st := NewStatusBadRequest("en", "email", "Invalid email", "INVALID", err)
		if !st.IsError() {
			t.Error("Expected IsError() to be true")
		}
		badReqs := st.BadRequests()
		if n := len(badReqs); n != 1 {
			t.Errorf("Bad request detail = %d, want %d", n, 1)
		}
		if field := badReqs[0].GetFieldViolations()[0].GetField(); field != "email" {
			t.Errorf("Field = %v, want %v", field, "email")
		}
	})
}

func TestStatus_Methods(t *testing.T) {
	t.Run("IsError and Error", func(t *testing.T) {
		successSt := &Status{Code: http.StatusOK, Message: "OK"}
		if successSt.IsError() {
			t.Errorf("Expected IsError() to be false for status %v", http.StatusOK)
		}
		if successSt.Error() != "" {
			t.Errorf("Expected Error() to be empty, got %v", successSt.Error())
		}

		badRequest := "Bad Request"
		errSt := &Status{Code: http.StatusBadRequest, Message: badRequest}
		if !errSt.IsError() {
			t.Errorf("Expected IsError() to be true for status %v", http.StatusBadRequest)
		}
		if got := errSt.Error(); got != badRequest {
			t.Errorf("Error() = %v, want %v", got, badRequest)
		}
	})

	t.Run("MessageOrDefault", func(t *testing.T) {
		mesCusExp := "Custom Not Found"
		stCustom := &Status{Code: http.StatusNotFound, Message: mesCusExp}
		if got := stCustom.MessageOrDefault(); got != mesCusExp {
			t.Errorf("Got '%v', want '%v'", got, mesCusExp)
		}

		mesDefExp := statusText(http.StatusNotFound)
		stDefault := &Status{Code: http.StatusNotFound}
		if got := stDefault.MessageOrDefault(); got != mesDefExp {
			t.Errorf("Got '%v', want '%v'", got, mesDefExp)
		}
	})

	t.Run("CauseErrorMesssageOrDefault", func(t *testing.T) {
		rootCause := "root cause"
		stWithCause := &Status{
			Code:       http.StatusInternalServerError,
			CauseError: errors.New(rootCause),
		}
		if got := stWithCause.CauseErrorMessageOrDefault(); got != rootCause {
			t.Errorf("Got '%v', want '%v'", got, rootCause)
		}

		serverError := "server error"
		stWithoutCause := &Status{Code: http.StatusInternalServerError, Message: serverError}
		if got := stWithoutCause.CauseErrorMessageOrDefault(); got != serverError {
			t.Errorf("Got '%v', want '%v'", got, serverError)
		}
	})

	t.Run("HTTPStatusCodeToGRPCCode and GRPCStatus", func(t *testing.T) {
		codesMap := map[int]codes.Code{
			http.StatusOK:                  codes.OK,
			http.StatusCreated:             codes.OK,
			http.StatusBadRequest:          codes.InvalidArgument,
			http.StatusUnauthorized:        codes.Unauthenticated,
			http.StatusForbidden:           codes.PermissionDenied,
			http.StatusNotFound:            codes.NotFound,
			http.StatusTooManyRequests:     codes.ResourceExhausted,
			http.StatusInternalServerError: codes.Internal,
			http.StatusServiceUnavailable:  codes.Unavailable,
			999:                            codes.Unknown,
		}

		for httpCode, wantCode := range codesMap {
			gotCode := HTTPStatusCodeToGRPCCode(httpCode)
			if gotCode != wantCode {
				t.Errorf("HTTPStatusCodeToGRPCCode(%d) = %v, want %v", httpCode, gotCode, wantCode)
			}
		}

		st := NewStatusCode(http.StatusBadRequest)
		grpcSt := st.GRPCStatus()
		if grpcSt.Code() != codes.InvalidArgument {
			t.Errorf("GRPCStatus code = %v, want %v", grpcSt.Code(), codes.InvalidArgument)
		}
	})
}

func statusText(code int) string {
	return http.StatusText(code)
}
