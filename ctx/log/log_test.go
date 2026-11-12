// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package log

import (
	"bytes"
	"context"
	"path/filepath"
	"sync"
	"testing"

	configEntity "github.com/dedyf5/resik/entities/config"
)

func resetLog() {
	once = sync.Once{}
	log = nil
}

func TestEnsureCorrelationID(t *testing.T) {
	c := context.Background()

	t.Run("Generates new correlation ID when missing", func(t *testing.T) {
		id, newC := EnsureCorrelationID(c)
		if id == "" {
			t.Errorf("expected non-empty correlation ID")
		}
		expectedID, ok := newC.Value(KeyCorrelationIDContext).(string)
		if !ok || expectedID != id {
			t.Errorf("expected correlation ID stored in context")
		}
	})

	t.Run("Returns existing correlation ID when present", func(t *testing.T) {
		existingID := "test-correlation-123"
		testC := context.WithValue(c, KeyCorrelationIDContext, existingID)
		id, _ := EnsureCorrelationID(testC)
		if id != existingID {
			t.Errorf("expected %q, got %q", existingID, id)
		}
	})
}

func TestKey_String(t *testing.T) {
	if KeyCorrelationIDContext.String() != "correlation_id" {
		t.Errorf("expected 'correlation_id', got %q", KeyCorrelationIDContext.String())
	}
	if KeyCallerHolderContext.String() != "caller" {
		t.Errorf("expected 'caller', got %q", KeyCallerHolderContext.String())
	}
}

func TestLogInitializationAndMethods(t *testing.T) {
	resetLog()

	logConfig := configEntity.Log{File: ""}
	moduleConfig := configEntity.Module{NameKey: "test_module"}

	loggerInst := Get(logConfig, moduleConfig)
	if loggerInst == nil {
		t.Fatal("expected logger instance, got nil")
	}

	queryString := "page=1&limit=10"
	loggerInst.Path = "/api/v1/test"
	loggerInst.QueryString = &queryString
	loggerInst.CorrelationID = "corr-test-id"

	t.Run("ZapFields construction", func(t *testing.T) {
		fields := loggerInst.ZapFields()
		if got := len(fields); got != 4 {
			t.Errorf("expected 4 fields, got %d", got)
		}
	})

	t.Run("Log level methods execution", func(t *testing.T) {
		loggerInst.Debug("debug message")
		loggerInst.Info("info message")
		loggerInst.Warn("warn message")
		loggerInst.Error("error message")
		loggerInst.Errorf("error format message: %s", "test")
	})
}

func TestLogInitializationWithFile(t *testing.T) {
	resetLog()

	tmpDir := t.TempDir()
	logFilePath := filepath.Join(tmpDir, "app.log")

	logConfig := configEntity.Log{File: logFilePath}
	moduleConfig := configEntity.Module{NameKey: "file_module"}

	loggerInst := Get(logConfig, moduleConfig)
	if loggerInst == nil {
		t.Fatal("expected logger instance with file, got nil")
	}

	loggerInst.Info("testing log file output")
}

func TestMaskBinaryFields(t *testing.T) {
	t.Run("Nil data handling", func(t *testing.T) {
		if res := maskBinaryFields(nil); res != nil {
			t.Errorf("expected nil, got %v", res)
		}
	})

	t.Run("Bytes Buffer handling", func(t *testing.T) {
		buf := bytes.NewBufferString(`{"status":"ok"}`)
		res := maskBinaryFields(buf)
		if res == nil {
			t.Error("expected processed buffer data")
		}
	})

	t.Run("Nil bytes.Buffer pointer handling", func(t *testing.T) {
		var nilBuf *bytes.Buffer
		res := maskBinaryFields(nilBuf)
		if res != nil {
			t.Errorf("expected nil for nil buffer, got %v", res)
		}
	})

	t.Run("Bytes Buffer with non-JSON content", func(t *testing.T) {
		buf := bytes.NewBufferString("plain text non-json data")
		res := maskBinaryFields(buf)
		strRes, ok := res.(string)
		if !ok || strRes == "" {
			t.Errorf("expected binary string representation, got %v", res)
		}
	})

	t.Run("Pointer to value handling", func(t *testing.T) {
		val := "pointer string value"
		res := maskBinaryFields(&val)
		if res != val {
			t.Errorf("expected '%s', got %v", val, res)
		}
	})

	t.Run("Nil pointer handling", func(t *testing.T) {
		var nilPtr *string
		res := maskBinaryFields(nilPtr)
		if res != nil {
			t.Errorf("expected nil for nil pointer, got %v", res)
		}
	})

	t.Run("Slice of non-byte elements (handleSliceMasking)", func(t *testing.T) {
		stringSlice := []string{"apple", "banana"}
		res := maskBinaryFields(stringSlice)
		sliceRes, ok := res.([]any)
		if !ok {
			t.Fatalf("expected []any slice result, got %T", res)
		}
		if len(sliceRes) != 2 || sliceRes[0] != "apple" {
			t.Errorf("unexpected slice masking result: %v", res)
		}
	})

	t.Run("Struct masking with cache, unexported fields, and binary fields", func(t *testing.T) {
		type ComplexPayload struct {
			PublicName string `json:"public_name"`
			privateKey string
			Data       []byte `json:"data"`
			NilData    []byte `json:"nil_data"`
		}

		payload := ComplexPayload{
			PublicName: "resik",
			privateKey: "secret-internal",
			Data:       []byte("raw binary content"),
			NilData:    nil,
		}

		masked1 := maskBinaryFields(payload)
		if masked1 == nil {
			t.Fatal("expected masked result")
		}

		payload2 := ComplexPayload{
			PublicName: "resik-v2",
			privateKey: "secret-internal-2",
			Data:       []byte("more binary"),
			NilData:    nil,
		}

		masked2 := maskBinaryFields(payload2)
		m, ok := masked2.(map[string]any)
		if !ok || m["public_name"] != "resik-v2" {
			t.Errorf("expected cached struct masking to work properly, got %v", masked2)
		}
	})

	t.Run("Struct field with protobuf tag name resolution", func(t *testing.T) {
		type ProtoPayload struct {
			FieldWithProto string `protobuf:"bytes,1,opt,name=proto_field_name,json=protoFieldName"`
		}

		payload := ProtoPayload{
			FieldWithProto: "proto_value",
		}

		masked := maskBinaryFields(payload)
		_, ok := masked.(map[string]any)
		if !ok {
			t.Errorf("expected map result from struct masking, got %T", masked)
		}
	})

	t.Run("Map with non-string keys (handleMapMasking fallback)", func(t *testing.T) {
		nonStringMap := map[int]any{
			100: "value_one",
			200: "value_two",
		}
		masked := maskBinaryFields(nonStringMap)
		m, ok := masked.(map[string]any)
		if !ok {
			t.Fatalf("expected map[string]any result, got %T", masked)
		}

		if m["100"] != "value_one" || m["200"] != "value_two" {
			t.Errorf("expected non-string keys to be formatted as string keys, got %v", m)
		}
	})

	maskedText := "[MASKED]"

	t.Run("Sensitive fields masking in Struct", func(t *testing.T) {
		type TestPayload struct {
			Username string
			Password string
			Token    string
		}

		payload := TestPayload{
			Username: "budi",
			Password: "secretpassword",
			Token:    "bearer-token-xyz",
		}

		masked := maskBinaryFields(payload)
		m, ok := masked.(map[string]any)
		if !ok {
			t.Errorf("expected map result from masking, got %T", masked)
		}

		if m["Username"] != "budi" {
			t.Errorf("expected username %q, got %q", "budi", m["Username"])
		}

		if m["Password"] != maskedText {
			t.Errorf("expected password %q, got %q", maskedText, m["Password"])
		}

		if m["Token"] != maskedText {
			t.Errorf("expected token %q, got %q", maskedText, m["Token"])
		}
	})

	t.Run("Map and Slice masking", func(t *testing.T) {
		dataMap := map[string]any{
			"password": "123",
			"normal":   "val",
		}
		maskedMap := maskBinaryFields(dataMap).(map[string]any)
		if maskedMap["password"] != maskedText {
			t.Errorf("expected password %q, got %q", maskedText, maskedMap["password"])
		}

		byteSlice := []byte("binary data")
		maskedSlice := maskBinaryFields(byteSlice)
		if _, ok := maskedSlice.(string); !ok {
			t.Errorf("expected binary slice formatted as string description")
		}
	})
}
