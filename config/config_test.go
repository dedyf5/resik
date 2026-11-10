// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package config

import (
	"context"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"time"

	"github.com/dedyf5/resik/drivers"
	configEntity "github.com/dedyf5/resik/entities/config"
	"github.com/spf13/viper"
)

func TestLoad_SuccessREST(t *testing.T) {
	viper.Reset()
	t.Cleanup(viper.Reset)
	viper.AutomaticEnv()

	setMinEnv(t)

	conf := Load(configEntity.ModuleTypeREST)

	if conf == nil {
		t.Fatal("expected loaded config to be non-nil")
	}

	if expected := "Resik App"; conf.App.Name() != expected {
		t.Errorf("expected App Name = %q, got %q", expected, conf.App.Name())
	}

	if expected := "resik"; conf.App.NameKey() != expected {
		t.Errorf("expected App Name Key = %q, got %q", expected, conf.App.NameKey())
	}

	if expected := "v1.0.0"; conf.App.Version() != expected {
		t.Errorf("expected App Version = %q, got %q", expected, conf.App.Version())
	}

	if conf.Database.Engine != drivers.MySQL {
		t.Errorf("expected Database Engine = %v, got %v", drivers.MySQL, conf.Database.Engine)
	}

	if expected := 10 * time.Second; conf.HTTP.ReadTimeout != expected {
		t.Errorf("expected HTTP ReadTimeout = %v, got %v", expected, conf.HTTP.ReadTimeout)
	}

	expectedName := "Resik App (Rest API)"
	if got := conf.AppModuleName(); got != expectedName {
		t.Errorf("expected AppModuleName = %q, got %q", expectedName, got)
	}
}

func TestLoad_NonRestModuleHTTP(t *testing.T) {
	viper.Reset()
	t.Cleanup(viper.Reset)
	viper.AutomaticEnv()

	setMinEnv(t)
	t.Setenv("GRPC_AUTH_EXPIRES", "24h")

	conf := Load(configEntity.ModuleTypeGRPC)
	if conf == nil {
		t.Fatal("expected loaded config to be non-nil")
	}

	if conf.HTTP.ReadTimeout != 0 {
		t.Errorf("expected HTTP ReadTimeout = %v, got %v", 0, conf.HTTP.ReadTimeout)
	}
}

func TestLoadModuleEnvProduction(t *testing.T) {
	viper.Reset()
	t.Cleanup(viper.Reset)
	viper.AutomaticEnv()

	setMinEnv(t)
	t.Setenv("REST_MODULE_ENV", "production")

	conf := Load(configEntity.ModuleTypeREST)
	if conf == nil {
		t.Fatal("expected loaded config to be non-nil")
	}

	if conf.Module.Env != configEntity.EnvProduction {
		t.Errorf("expected App Environment = %v, got %v", configEntity.EnvProduction, conf.Module.Env)
	}
}

func TestLoad_WithEnvFile(t *testing.T) {
	viper.Reset()
	t.Cleanup(viper.Reset)
	viper.AutomaticEnv()

	tmpDir := t.TempDir()
	t.Chdir(tmpDir)

	err := os.WriteFile(".env", []byte("APP_NAME=ResikEnvFile\n"), 0600)
	if err != nil {
		t.Fatalf("failed to create temporary .env file: %v", err)
	}

	setMinEnv(t)

	conf := Load(configEntity.ModuleTypeREST)
	if conf == nil {
		t.Fatalf("expected loaded config to be non-nil")
	}
}

func TestLoad_EnvConfigReadError(t *testing.T) {
	viper.Reset()
	t.Cleanup(viper.Reset)
	viper.AutomaticEnv()

	tmpDir := t.TempDir()
	t.Chdir(tmpDir)

	err := os.Mkdir(".env", 0700)
	if err != nil {
		t.Fatalf("failed to create temporary .env directory: %v", err)
	}

	setMinEnv(t)

	conf := Load(configEntity.ModuleTypeREST)
	if conf == nil {
		t.Fatal("expected loaded config to be non-nil")
	}
}

func TestLoad_WithSecretFile(t *testing.T) {
	viper.Reset()
	t.Cleanup(viper.Reset)
	viper.AutomaticEnv()

	tmpDir := t.TempDir()
	t.Chdir(tmpDir)

	secretFileName := "db_password.txt"
	err := os.WriteFile(secretFileName, []byte("my-secret-password-from-file\n"), 0600)
	if err != nil {
		t.Fatalf("failed to create temporary secret file: %v", err)
	}

	setMinEnv(t)
	t.Setenv("DATABASE_PASSWORD_PATH_FILE", secretFileName)

	conf := Load(configEntity.ModuleTypeREST)
	if conf == nil {
		t.Fatal("expected loaded config to be non-nil")
	}

	if expected := "my-secret-password-from-file"; conf.Database.Password != expected {
		t.Errorf("expected Database Password = %q, got %q", expected, conf.Database.Password)
	}
}

func TestLoad_SecretFileReadError(t *testing.T) {
	viper.Reset()
	t.Cleanup(viper.Reset)
	viper.AutomaticEnv()

	if os.Getenv("BE_CRASHER") == "1" {
		setMinEnv(t)
		t.Setenv("DATABASE_PASSWORD_PATH_FILE", "non_existent_secret_file.txt")
		Load(configEntity.ModuleTypeREST)
		return
	}

	ctx := context.Background()

	// #nosec G204 G702
	cmd := exec.CommandContext(ctx, os.Args[0], "-test.run=TestLoad_SecretFileReadError")
	cmd.Env = append(cmd.Environ(), "BE_CRASHER=1")
	err := cmd.Run()

	var exitErr *exec.ExitError
	if errors.As(err, &exitErr) && !exitErr.Success() {
		return
	}

	t.Fatalf("process ran with err %v, want exit status 1", err)
}

func TestLoad_DatabaseEnginePostgres(t *testing.T) {
	viper.Reset()
	t.Cleanup(viper.Reset)
	viper.AutomaticEnv()

	setMinEnv(t)
	t.Setenv("DATABASE_ENGINE", "postgres")

	conf := Load(configEntity.ModuleTypeREST)
	if conf == nil {
		t.Fatal("expected loaded config to be non-nil")
	}

	if conf.Database.Engine != drivers.PostgreSQL {
		t.Errorf("expected Database Engine = %v, got %v", drivers.PostgreSQL, conf.Database.Engine)
	}
}

func TestLoad_DatabaseEngineEmpty(t *testing.T) {
	viper.Reset()
	t.Cleanup(viper.Reset)
	viper.AutomaticEnv()

	if os.Getenv("BE_CRASHER") == "1" {
		setMinEnv(t)
		t.Setenv("DATABASE_ENGINE", "")
		Load(configEntity.ModuleTypeREST)
		return
	}

	ctx := context.Background()

	// #nosec G204 G702
	cmd := exec.CommandContext(ctx, os.Args[0], "-test.run=TestLoad_DatabaseEngineEmpty")
	cmd.Env = append(cmd.Environ(), "BE_CRASHER=1")
	err := cmd.Run()

	var exitErr *exec.ExitError
	if errors.As(err, &exitErr) && !exitErr.Success() {
		return
	}

	t.Fatalf("process ran with err %v, want exit status 1", err)
}

func TestLoad_WithoutRedis(t *testing.T) {
	viper.Reset()
	t.Cleanup(viper.Reset)
	viper.AutomaticEnv()

	setMinEnv(t)
	t.Setenv("REDIS_HOST", "")

	conf := Load(configEntity.ModuleTypeREST)
	if conf == nil {
		t.Fatal("expected loaded config to be non-nil")
	}

	if conf.Redis != nil {
		t.Errorf("expected empty Redis, got %v", conf.Redis)
	}
}

func TestLoad_RateLimitDriverEmpty(t *testing.T) {
	viper.Reset()
	t.Cleanup(viper.Reset)
	viper.AutomaticEnv()

	setMinEnv(t)
	t.Setenv("RATE_LIMIT_DRIVER", "")

	conf := Load(configEntity.ModuleTypeREST)
	if conf == nil {
		t.Fatal("expected loaded config to be non-nil")
	}

	if conf.RateLimit.Driver != configEntity.RateLimitDriverMemory {
		t.Errorf("expected RateLimit Driver = %v, got %v", configEntity.RateLimitDriverMemory, conf.RateLimit.Driver)
	}
}

func TestLoad_RateLimitDriverRedis(t *testing.T) {
	viper.Reset()
	t.Cleanup(viper.Reset)
	viper.AutomaticEnv()

	setMinEnv(t)
	t.Setenv("RATE_LIMIT_DRIVER", "redis")

	conf := Load(configEntity.ModuleTypeREST)
	if conf == nil {
		t.Fatal("expected loaded config to be non-nil")
	}

	if conf.RateLimit.Driver != configEntity.RateLimitDriverRedis {
		t.Errorf("expected RateLimit Driver = %v, got %v", configEntity.RateLimitDriverRedis, conf.RateLimit.Driver)
	}
}

func setMinEnv(t *testing.T) {
	// App env keys
	t.Setenv("APP_NAME", "Resik App")
	t.Setenv("APP_NAME_KEY", "resik")
	t.Setenv("APP_VERSION", "v1.0.0")

	// Module REST env keys
	t.Setenv("REST_MODULE_NAME", "Rest API")
	t.Setenv("REST_MODULE_NAME_KEY", "rest-api")
	t.Setenv("REST_MODULE_HOST", "127.0.0.1")
	t.Setenv("REST_MODULE_PORT", "8080")

	// HTTP timeouts
	t.Setenv("REST_HTTP_READ_HEADER_TIMEOUT", "5s")
	t.Setenv("REST_HTTP_READ_TIMEOUT", "10s")
	t.Setenv("REST_HTTP_WRITE_TIMEOUT", "10s")
	t.Setenv("REST_HTTP_IDLE_TIMEOUT", "60s")

	// Database config
	t.Setenv("DATABASE_ENGINE", "mysql")
	t.Setenv("DATABASE_HOST", "localhost")
	t.Setenv("DATABASE_PORT", "3306")
	t.Setenv("DATABASE_USERNAME", "root")
	t.Setenv("DATABASE_PASSWORD", "secret")
	t.Setenv("DATABASE_SCHEMA", "resik_db")
	t.Setenv("DATABASE_HEALTHCHECK_TIMEOUT", "3s")
	t.Setenv("REST_DATABASE_MAX_OPEN_CONS", "20")
	t.Setenv("REST_DATABASE_MAX_IDLE_CONS", "5")
	t.Setenv("REST_DATABASE_CONN_MAX_LIFETIME", "15m")
	t.Setenv("REST_DATABASE_CONN_MAX_IDLETIME", "10m")

	// Redis config
	t.Setenv("REDIS_HOST", "127.0.0.1")
	t.Setenv("REDIS_PORT", "6379")
	t.Setenv("REDIS_USERNAME", "default")
	t.Setenv("REDIS_DATABASE", "0")
	t.Setenv("REDIS_POOL_SIZE", "10")
	t.Setenv("REDIS_HEALTHCHECK_TIMEOUT", "3s")

	// Rate Limit config
	t.Setenv("RATE_LIMIT_DRIVER", "memory")
	t.Setenv("RATE_LIMIT_PERIOD", "1m")
	t.Setenv("RATE_LIMIT_LIMIT", "100")
	t.Setenv("RATE_LIMIT_PREFIX", "resik:rl:")

	// Auth config
	t.Setenv("REST_AUTH_EXPIRES", "24h")
	t.Setenv("REST_AUTH_SIGNATURE_KEY", "super-secret-signature-key")

	// Log config
	t.Setenv("REST_LOG_FILE", "app.log")
}

func TestReadSecretFile(t *testing.T) {
	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "secret.txt")
	// #nosec G101
	secretContent := "  my-super-secret-token\n"

	err := os.WriteFile(tmpFile, []byte(secretContent), 0600)
	if err != nil {
		t.Fatalf("failed to create temporary secret file: %v", err)
	}

	t.Run("Success reading secret file with trim space", func(t *testing.T) {
		val, err := readSecretFile(tmpFile)
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		expected := "my-super-secret-token"
		if val != expected {
			t.Errorf("expected %q, got %q", expected, val)
		}
	})

	t.Run("Error empty path", func(t *testing.T) {
		_, err := readSecretFile("")
		if err == nil {
			t.Error("expected error for empty path, got nil")
		}
	})

	t.Run("Error non-existent file", func(t *testing.T) {
		_, err := readSecretFile(filepath.Join(tmpDir, "non-existent.txt"))
		if err == nil {
			t.Error("expected error for non-existent file, got nil")
		}
	})
}

func TestGetIntOrDefault(t *testing.T) {
	viper.Reset()
	t.Cleanup(viper.Reset)

	viper.Set("VALID_INT", 100)

	t.Run("Retrieve existing valid integer", func(t *testing.T) {
		val := getIntOrDefault("VALID_INT", 10)
		if val != 100 {
			t.Errorf("expected %d, got %d", 100, val)
		}
	})

	t.Run("Retrieve fallback default when key is missing or zero", func(t *testing.T) {
		val := getIntOrDefault("MISSING_INT", 25)
		if val != 25 {
			t.Errorf("expected %d, got %d", 25, val)
		}
	})
}

func TestGetUint32OrDefault(t *testing.T) {
	viper.Reset()
	t.Cleanup(viper.Reset)
	viper.AutomaticEnv()

	viper.Set("VALID_UINT32", uint32(16384))

	t.Run("Retrieve existing valid uint32", func(t *testing.T) {
		val := getUint32OrDefault("VALID_UINT32", 1024)
		if val != 16384 {
			t.Errorf("expected %d, got %d", 16384, val)
		}
	})

	t.Run("Retrieve fallback default when key is missing or zero", func(t *testing.T) {
		val := getUint32OrDefault("MISSING_UINT32", 2048)
		if val != 2048 {
			t.Errorf("expected %d, got %d", 2048, val)
		}
	})
}

func TestGetDuration_EmptyError(t *testing.T) {
	viper.Reset()
	t.Cleanup(viper.Reset)
	viper.AutomaticEnv()

	if os.Getenv("BE_CRASHER") == "1" {
		setMinEnv(t)
		t.Setenv("DATABASE_HEALTHCHECK_TIMEOUT", "")
		Load(configEntity.ModuleTypeREST)
		return
	}

	ctx := context.Background()

	// #nosec G204 G702
	cmd := exec.CommandContext(ctx, os.Args[0], "-test.run=TestGetDuration_EmptyError")
	cmd.Env = append(cmd.Environ(), "BE_CRASHER=1")
	err := cmd.Run()

	var exitErr *exec.ExitError
	if errors.As(err, &exitErr) && !exitErr.Success() {
		return
	}

	t.Fatalf("process ran with err %v, want exit status 1", err)
}

func TestGetDuration_InvalidFormatError(t *testing.T) {
	viper.Reset()
	t.Cleanup(viper.Reset)
	viper.AutomaticEnv()

	if os.Getenv("BE_CRASHER") == "1" {
		setMinEnv(t)
		t.Setenv("DATABASE_HEALTHCHECK_TIMEOUT", "invalid-duration-format")
		Load(configEntity.ModuleTypeREST)
		return
	}

	ctx := context.Background()

	// #nosec G204 G702
	cmd := exec.CommandContext(ctx, os.Args[0], "-test.run=TestGetDuration_InvalidFormatError")
	cmd.Env = append(cmd.Environ(), "BE_CRASHER=1")
	err := cmd.Run()

	var exitErr *exec.ExitError
	if errors.As(err, &exitErr) && !exitErr.Success() {
		return
	}

	t.Fatalf("process ran with err %v, want exit status 1", err)
}

func TestGetDurationOrDefault(t *testing.T) {
	viper.Reset()
	t.Cleanup(viper.Reset)

	viper.Set("VALID_DURATION", "30s")

	t.Run("Retrieve existing valid duration", func(t *testing.T) {
		val := getDurationOrDefault("VALID_DURATION", 5*time.Second)
		if expected := 30 * time.Second; val != expected {
			t.Errorf("expected %v, got %v", expected, val)
		}
	})

	t.Run("Retrieve fallback default when duration key is missing", func(t *testing.T) {
		val := getDurationOrDefault("MISSING_DURATION", 10*time.Second)
		if expected := 10 * time.Second; val != expected {
			t.Errorf("expected %v, got %v", expected, val)
		}
	})
}

func TestGetDurationOrDefault_InvalidFormat(t *testing.T) {
	viper.Reset()
	t.Cleanup(viper.Reset)

	viper.Set("INVALID_DURATION", "not-a-duration")

	if os.Getenv("BE_CRASHER") == "1" {
		getDurationOrDefault("INVALID_DURATION", 5*time.Second)
		return
	}

	ctx := context.Background()

	// #nosec G204 G702
	cmd := exec.CommandContext(ctx, os.Args[0], "-test.run=TestGetDurationOrDefault_InvalidFormat")
	cmd.Env = append(cmd.Environ(), "BE_CRASHER=1")
	err := cmd.Run()

	var exitErr *exec.ExitError
	if errors.As(err, &exitErr) && !exitErr.Success() {
		return
	}

	t.Fatalf("process ran with err %v, want exit status 1", err)
}

func TestGetApp(t *testing.T) {
	viper.Reset()
	t.Cleanup(viper.Reset)

	viper.AutomaticEnv()

	t.Setenv("APP_NAME", "Resik Test Engine")
	t.Setenv("APP_NAME_KEY", "resik-test")
	t.Setenv("APP_VERSION", "v1.2.3")

	app := GetApp()

	if expected := "Resik Test Engine"; app.Name() != expected {
		t.Errorf("expected %q, got %q", expected, app.Name())
	}

	if expected := "resik-test"; app.NameKey() != expected {
		t.Errorf("expected %q, got %q", expected, app.NameKey())
	}

	if expected := "v1.2.3"; app.Version() != expected {
		t.Errorf("expected %q, got %q", expected, app.Version())
	}
}
