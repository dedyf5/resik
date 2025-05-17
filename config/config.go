// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package config

import (
	"fmt"
	"log"

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
	viper.SetConfigType("env")
	viper.SetConfigFile(fmt.Sprintf("./app/%s/config/.env", module.DirectoryName()))
	if err := viper.ReadInConfig(); err != nil {
		log.Print("ERROR read env", err.Error())
	}

	viper.AutomaticEnv()

	conf := Config{}
	conf.loadApp(module)
	conf.loadHTTP()
	conf.loadDatabase()
	conf.loadAuth()
	conf.loadLog()

	return &conf
}

func (conf *Config) loadApp(module configEntity.Module) {
	envStr := viper.GetString("ENV")
	env := configEntity.EnvDevelopment
	switch envStr {
	case "staging":
		env = configEntity.EnvStaging
	case "production":
		env = configEntity.EnvProduction
	}
	conf.App = configEntity.App{
		Name:        viper.GetString("APP_NAME"),
		Version:     viper.GetString("APP_VERSION"),
		Module:      module,
		Env:         env,
		LangDefault: langCtx.GetLanguageOrDefault(viper.GetString("APP_LANG_DEFAULT")),
		Host:        viper.GetString("APP_HOST"),
		Port:        viper.GetUint("APP_PORT"),
	}
}

func (conf *Config) loadHTTP() {
	conf.HTTP = configEntity.HTTP{
		ClientTimeout: viper.GetUint("HTTP_CLIENT_TIMEOUT"),
	}
}

func (conf *Config) loadDatabase() {
	db := drivers.SQLConfig{}
	switch viper.GetString("DATABASE_ENGINE") {
	case "mysql":
		db.Engine = drivers.MySQL
	case "postgres":
		db.Engine = drivers.PostgreSQL
	}

	db.Host = viper.GetString("DATABASE_HOST")
	db.Port = viper.GetInt("DATABASE_PORT")
	db.Username = viper.GetString("DATABASE_USERNAME")
	db.Password = viper.GetString("DATABASE_PASSWORD")
	db.Schema = viper.GetString("DATABASE_SCHEMA")
	db.MaxOpenConns = viper.GetInt("DATABASE_MAX_OPEN_CONS")
	db.MaxIdleConns = viper.GetInt("DATABASE_MAX_IDLE_CONS")
	db.ConnMaxLifetime = viper.GetInt("DATABASE_CONN_MAX_LIFETIME")
	db.ConnMaxIdleTime = viper.GetInt("DATABASE_CONN_MAX_IDLETIME")
	db.IsDebug = viper.GetBool("DATABASE_IS_DEBUG")
	conf.Database = db
}

func (conf *Config) loadAuth() {
	conf.Auth = configEntity.Auth{
		Expires:      viper.GetUint64("AUTH_EXPIRES"),
		SignatureKey: viper.GetString("AUTH_SIGNATURE_KEY"),
	}
}

func (conf *Config) loadLog() {
	conf.Log = configEntity.Log{
		File: viper.GetString("LOG_FILE"),
	}
}
