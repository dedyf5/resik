// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package config

import (
	"fmt"
	"log"
	"os"
	"strings"

	langCtx "github.com/dedyf5/resik/ctx/lang"
	"github.com/dedyf5/resik/drivers"
	configEntity "github.com/dedyf5/resik/entities/config"
	"github.com/spf13/viper"
)

const AppLogoASCII string = `
 ____           _ _    
|  _ \ ___  ___(_) | __
| |_) / _ \/ __| | |/ /
|  _ <  __/\__ \ |   < 
|_| \_\___||___/_|_|\_\

`

type (
	Config struct {
		App      configEntity.App
		HTTP     configEntity.HTTP
		Database drivers.SQLConfig
		Auth     configEntity.Auth
		Log      configEntity.Log
	}
)

func Load(module configEntity.Module) *Config {
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
	conf.loadApp(module)
	conf.loadHTTP(module)
	conf.loadDatabase(module)
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

func readSecretFile(filepath string) (string, error) {
	if filepath == "" {
		return "", fmt.Errorf("secret file path is empty")
	}
	content, err := os.ReadFile(filepath)
	if err != nil {
		return "", fmt.Errorf("failed to read secret file %s: %w", filepath, err)
	}
	return strings.TrimSpace(string(content)), nil
}

func (conf *Config) loadApp(module configEntity.Module) {
	envStr := viper.GetString(module.Key("APP_ENV"))
	env := configEntity.EnvDevelopment
	switch envStr {
	case "staging":
		env = configEntity.EnvStaging
	case "production":
		env = configEntity.EnvProduction
	}
	conf.App = configEntity.App{
		Name:        viper.GetString(module.Key("APP_NAME")),
		Version:     viper.GetString(module.Key("APP_VERSION")),
		Module:      module,
		Env:         env,
		LangDefault: langCtx.GetLanguageOrDefault(viper.GetString(module.Key("APP_LANG_DEFAULT"))),
		Host:        viper.GetString(module.Key("APP_HOST")),
		Port:        viper.GetUint(module.Key("APP_PORT")),
		Public: configEntity.AppPublic{
			Host:     viper.GetString(module.Key("APP_PUBLIC_HOST")),
			Port:     viper.GetUint(module.Key("APP_PUBLIC_PORT")),
			Schema:   viper.GetString(module.Key("APP_PUBLIC_SCHEMA")),
			BasePath: viper.GetString(module.Key("APP_PUBLIC_BASE_PATH")),
		},
	}
}

func (conf *Config) loadHTTP(module configEntity.Module) {
	conf.HTTP = configEntity.HTTP{
		ClientTimeout: viper.GetUint(module.Key("HTTP_CLIENT_TIMEOUT")),
	}
}

func (conf *Config) loadDatabase(module configEntity.Module) {
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
	db.ConnMaxLifetime = viper.GetInt(module.Key("DATABASE_CONN_MAX_LIFETIME"))
	db.ConnMaxIdleTime = viper.GetInt(module.Key("DATABASE_CONN_MAX_IDLETIME"))
	db.IsDebug = viper.GetBool(module.Key("DATABASE_IS_DEBUG"))
	conf.Database = db
}

func (conf *Config) loadAuth(module configEntity.Module) {
	conf.Auth = configEntity.Auth{
		Expires:      viper.GetUint64(module.Key("AUTH_EXPIRES")),
		SignatureKey: getSecretFromFileOrEnv("AUTH_SIGNATURE_KEY_PATH_FILE", "AUTH_SIGNATURE_KEY"),
	}
}

func (conf *Config) loadLog(module configEntity.Module) {
	conf.Log = configEntity.Log{
		File: viper.GetString(module.Key("LOG_FILE")),
	}
}
