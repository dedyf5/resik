// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package config

import (
	"errors"
	"fmt"
	"log/slog"
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

var logger *slog.Logger

func init() {
	opts := &slog.HandlerOptions{
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == slog.TimeKey {
				t := a.Value.Time().UTC()
				return slog.Time(slog.TimeKey, t)
			}

			return a
		},
	}

	logger = slog.New(slog.NewTextHandler(os.Stdout, opts))

	slog.SetDefault(logger)

	logger = logger.With("component", "config")
}

func Load(module configEntity.ModuleType) *Config {
	if _, err := os.Stat(".env"); err == nil {
		viper.SetConfigType("env")
		viper.SetConfigFile(".env")
		if err := viper.ReadInConfig(); err != nil {
			logger.Warn("failed to read .env file: " + err.Error())
		}
	} else {
		logger.Info("no .env file found, relying on environment variables")
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
			logExitError("failed to read secret from file pointed by " + secretFilePathEnvVarName + " (" + secretFilePath + "): " + err.Error())
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
	envKey := module.Key("MODULE_ENV")
	envValue := configEntity.Env(viper.GetString(envKey))
	if envValue == "" {
		envValue = configEntity.EnvDevelopment
		logWarnEmptyValue(envKey, envValue)
	}

	conf.Module = configEntity.Module{
		Name:        viper.GetString(module.Key("MODULE_NAME")),
		NameKey:     viper.GetString(module.Key("MODULE_NAME_KEY")),
		Type:        module,
		Env:         envValue,
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
	db.MaxOpenConns = getIntOrDefault(module.Key("DATABASE_MAX_OPEN_CONS"), 30)
	db.MaxIdleConns = getIntOrDefault(module.Key("DATABASE_MAX_IDLE_CONS"), 5)
	db.ConnMaxLifetime = getDurationOrDefault(module.Key("DATABASE_CONN_MAX_LIFETIME"), 5*time.Minute)
	db.ConnMaxIdleTime = getDurationOrDefault(module.Key("DATABASE_CONN_MAX_IDLETIME"), 5*time.Minute)
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

	driverKey := "RATE_LIMIT_DRIVER"
	rl.Driver = configEntity.RateLimitDriver(viper.GetString(driverKey))
	if rl.Driver == "" {
		rl.Driver = configEntity.RateLimitDriverMemory
		logWarnEmptyValue(driverKey, configEntity.RateLimitDriverMemory)
	}

	rl.Period = getDuration("RATE_LIMIT_PERIOD")
	rl.Limit = viper.GetInt64("RATE_LIMIT_LIMIT")
	rl.Prefix = viper.GetString("RATE_LIMIT_PREFIX")

	conf.RateLimit = rl
}

func (conf *Config) loadAuth(module configEntity.ModuleType) {
	conf.Auth = configEntity.Auth{
		Expires:        getDuration(module.Key("AUTH_EXPIRES")),
		SignatureKey:   getSecretFromFileOrEnv("AUTH_SIGNATURE_KEY_PATH_FILE", "AUTH_SIGNATURE_KEY"),
		HashMemory:     getUint32OrDefault("AUTH_HASH_MEMORY", 8*1024),
		HashIterations: getUint32OrDefault("AUTH_HASH_ITERATIONS", 3),
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

// getIntOrDefault retrieves an integer value from the configuration using the provided key
// and returns it. If the key is missing or the value is zero, it will return `defaultValue`.
func getIntOrDefault(key string, defaultValue int) (value int) {
	if val := viper.GetInt(key); val != 0 {
		value = val
	} else {
		value = defaultValue
		logWarnEmptyValue(key, defaultValue)
	}
	return
}

// getUint32OrDefault retrieves a uint32 value from the configuration using the provided key
// and returns it. If the key is missing or the value is zero, it will return `defaultValue`.
func getUint32OrDefault(key string, defaultValue uint32) (value uint32) {
	if val := viper.GetUint32(key); val != 0 {
		value = val
	} else {
		value = defaultValue
		logWarnEmptyValue(key, defaultValue)
	}
	return
}

// getDuration retrieves a string value from the configuration using the provided key
// and parses it into a time.Duration.
//
// The value in the configuration should follow the format supported by time.ParseDuration,
// such as "300s", "5m", or "24h". If the key is missing or the value is not a valid
// duration string, the application will log a fatal error and exit.
func getDuration(key string) time.Duration {
	str := viper.GetString(key)
	if str == "" {
		logExitError("field '" + key + "' is required")
	}

	duration, err := time.ParseDuration(str)
	if err != nil {
		logErrorInvalidFormatDuration(key, err)
	}
	return duration
}

// getDurationOrDefault retrieves a string value from the configuration using the provided key
// and parses it into a time.Duration.
//
// The value in the configuration should follow the format supported by time.ParseDuration,
// such as "300s", "5m", or "24h". If the key is missing or empty, it will return `defaultValue`.
// If the value is not a valid duration string, the application will log a fatal error and exit.
func getDurationOrDefault(key string, defaultValue time.Duration) time.Duration {
	str := viper.GetString(key)
	if str == "" {
		logWarnEmptyValue(key, defaultValue)
		return defaultValue
	}

	duration, err := time.ParseDuration(str)
	if err != nil {
		logErrorInvalidFormatDuration(key, err)
	}
	return duration
}

func logWarnEmptyValue(key string, defaultValue any) {
	logger.Warn("field '" + key + "' is empty, using default value '" + fmt.Sprintf("%v", defaultValue) + "'")
}

func logErrorInvalidFormatDuration(key string, err error) {
	logExitError("field '" + key + "' has an invalid format: " + err.Error() + " (valid examples: '1h', '30m', '3600s')")
}

func logExitError(msg string) {
	logger.Error(msg)
	os.Exit(1)
}
