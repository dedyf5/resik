// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package config

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	langCtx "github.com/dedyf5/resik/ctx/lang"
	"github.com/dedyf5/resik/drivers"
	configEntity "github.com/dedyf5/resik/entities/config"
	"github.com/spf13/viper"
)

type (
	Config struct {
		App       configEntity.App
		Module    configEntity.Module
		HTTP      configEntity.HTTP
		Database  drivers.SQLConfig
		Redis     *drivers.RedisConfig
		RateLimit configEntity.RateLimit
		Auth      configEntity.Auth
		Log       configEntity.Log
	}
)

func Load(module configEntity.ModuleType) *Config {
	if _, err := os.Stat(".env"); err == nil {
		viper.SetConfigType("env")
		viper.SetConfigFile(".env")
		if err := viper.ReadInConfig(); err != nil {
			log.Printf("WARNING: Failed to read .env file: %v", err)
		}
	} else {
		log.Print("INFO: No .env file found, relying on environment variables")
	}

	viper.AutomaticEnv()

	conf := Config{}
	conf.loadApp()
	conf.loadModule(module)
	conf.loadHTTP(module)
	conf.loadDatabase(module)
	conf.loadRedis()
	conf.loadRateLimit()
	conf.loadAuth(module)
	conf.loadLog(module)

	return &conf
}

func getSecretFromFileOrEnv(secretFilePathEnvVarName, fallbackEnvVarName string) string {
	secretFilePath := viper.GetString(secretFilePathEnvVarName)

	if secretFilePath != "" {
		content, err := readSecretFile(secretFilePath)
		if err != nil {
			log.Fatalf("FATAL: Failed to read secret from file pointed by %s (%s): %v", secretFilePathEnvVarName, secretFilePath, err)
		}
		return content
	} else {
		fallbackValue := viper.GetString(fallbackEnvVarName)
		return fallbackValue
	}
}

func readSecretFile(path string) (string, error) {
	if path == "" {
		return "", errors.New("secret file path is empty")
	}

	cleanPath := filepath.Clean(path)

	content, err := os.ReadFile(cleanPath)
	if err != nil {
		return "", fmt.Errorf("failed to read secret file %s: %w", cleanPath, err)
	}

	return strings.TrimSpace(string(content)), nil
}

func (conf *Config) loadApp() {
	conf.App = *GetApp()
}

func GetApp() *configEntity.App {
	return configEntity.NewApp(
		viper.GetString("APP_NAME"),
		viper.GetString("APP_NAME_KEY"),
		viper.GetString("APP_VERSION"),
	)
}

func (conf *Config) loadModule(module configEntity.ModuleType) {
	envStr := viper.GetString(module.Key("MODULE_ENV"))
	env := configEntity.EnvDevelopment
	switch envStr {
	case "staging":
		env = configEntity.EnvStaging
	case "production":
		env = configEntity.EnvProduction
	}
	conf.Module = configEntity.Module{
		Name:        viper.GetString(module.Key("MODULE_NAME")),
		NameKey:     viper.GetString(module.Key("MODULE_NAME_KEY")),
		Type:        module,
		Env:         env,
		LangDefault: langCtx.GetLanguageOrDefault(viper.GetString(module.Key("MODULE_LANG_DEFAULT"))),
		Host:        viper.GetString(module.Key("MODULE_HOST")),
		Port:        viper.GetUint(module.Key("MODULE_PORT")),
		Public: configEntity.Public{
			Host:     viper.GetString(module.Key("MODULE_PUBLIC_HOST")),
			Port:     viper.GetUint(module.Key("MODULE_PUBLIC_PORT")),
			Schema:   viper.GetString(module.Key("MODULE_PUBLIC_SCHEMA")),
			BasePath: viper.GetString(module.Key("MODULE_PUBLIC_BASE_PATH")),
		},
	}
}

func (conf *Config) loadHTTP(module configEntity.ModuleType) {
	if conf.Module.Type != configEntity.ModuleTypeREST {
		return
	}

	conf.HTTP = configEntity.HTTP{
		ReadHeaderTimeout: viper.GetDuration(module.Key("HTTP_READ_HEADER_TIMEOUT")),
		ReadTimeout:       viper.GetDuration(module.Key("HTTP_READ_TIMEOUT")),
		WriteTimeout:      viper.GetDuration(module.Key("HTTP_WRITE_TIMEOUT")),
		IdleTimeout:       viper.GetDuration(module.Key("HTTP_IDLE_TIMEOUT")),
	}
}

func (conf *Config) loadDatabase(module configEntity.ModuleType) {
	db := drivers.SQLConfig{}
	switch viper.GetString("DATABASE_ENGINE") {
	case "mysql":
		db.Engine = drivers.MySQL
	case "postgres":
		db.Engine = drivers.PostgreSQL
	}

	db.Host = viper.GetString("DATABASE_HOST")
	db.Port = viper.GetInt("DATABASE_PORT")
	db.Username = getSecretFromFileOrEnv("DATABASE_USERNAME_PATH_FILE", "DATABASE_USERNAME")
	db.Password = getSecretFromFileOrEnv("DATABASE_PASSWORD_PATH_FILE", "DATABASE_PASSWORD")
	db.Schema = viper.GetString("DATABASE_SCHEMA")
	db.MaxOpenConns = viper.GetInt(module.Key("DATABASE_MAX_OPEN_CONS"))
	db.MaxIdleConns = viper.GetInt(module.Key("DATABASE_MAX_IDLE_CONS"))
	db.ConnMaxLifetime = getDuration(module.Key("DATABASE_CONN_MAX_LIFETIME"))
	db.ConnMaxIdleTime = getDuration(module.Key("DATABASE_CONN_MAX_IDLETIME"))
	db.HealthCheckTimeout = getDuration("DATABASE_HEALTHCHECK_TIMEOUT")
	db.IsDebug = viper.GetBool(module.Key("DATABASE_IS_DEBUG"))
	conf.Database = db
}

func (conf *Config) loadRedis() {
	host := viper.GetString("REDIS_HOST")
	if host == "" {
		return
	}

	redis := drivers.RedisConfig{}

	redis.Host = host
	redis.Port = viper.GetInt("REDIS_PORT")
	redis.Username = getSecretFromFileOrEnv("REDIS_USERNAME_PATH_FILE", "REDIS_USERNAME")
	redis.Password = getSecretFromFileOrEnv("REDIS_PASSWORD_PATH_FILE", "REDIS_PASSWORD")
	redis.Database = viper.GetInt("REDIS_DATABASE")
	redis.PoolSize = viper.GetInt("REDIS_POOL_SIZE")
	redis.HealthCheckTimeout = getDuration("REDIS_HEALTHCHECK_TIMEOUT")

	conf.Redis = &redis
}

func (conf *Config) loadRateLimit() {
	rl := configEntity.RateLimit{}

	rl.Driver = configEntity.RateLimitDriver(viper.GetString("RATE_LIMIT_DRIVER"))
	rl.Period = getDuration("RATE_LIMIT_PERIOD")
	rl.Limit = viper.GetInt64("RATE_LIMIT_LIMIT")
	rl.Prefix = viper.GetString("RATE_LIMIT_PREFIX")

	conf.RateLimit = rl
}

func (conf *Config) loadAuth(module configEntity.ModuleType) {
	conf.Auth = configEntity.Auth{
		Expires:        getDuration(module.Key("AUTH_EXPIRES")),
		SignatureKey:   getSecretFromFileOrEnv("AUTH_SIGNATURE_KEY_PATH_FILE", "AUTH_SIGNATURE_KEY"),
		HashMemory:     viper.GetUint32("AUTH_HASH_MEMORY"),
		HashIterations: viper.GetUint32("AUTH_HASH_ITERATIONS"),
	}
}

func (conf *Config) loadLog(module configEntity.ModuleType) {
	conf.Log = configEntity.Log{
		File: viper.GetString(module.Key("LOG_FILE")),
	}
}

func (conf *Config) AppModuleName() string {
	return conf.App.Name() + " (" + conf.Module.Name + ")"
}

// getDuration retrieves a string value from the configuration using the provided key
// and parses it into a time.Duration.
//
// The value in the configuration should follow the format supported by time.ParseDuration,
// such as "300s", "5m", or "24h". If the key is missing or the value is not a valid
// duration string, the application will log a fatal error and exit.
func getDuration(key string) (duration time.Duration) {
	duration, err := time.ParseDuration(viper.GetString(key))
	if err != nil {
		log.Fatalf("[config] field '%s' has an invalid format: %v (valid examples: '1h', '30m', '3600s')", key, err)
	}
	return
}
