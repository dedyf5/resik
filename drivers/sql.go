// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package drivers

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/pkg/errors"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type SQLEngine string

const (
	MySQL      SQLEngine = "mysql"
	PostgreSQL SQLEngine = "postgres"
)

func (engine SQLEngine) String() string {
	return string(engine)
}

type SQLConfig struct {
	Engine             SQLEngine
	Host               string
	Port               int
	Username           string
	Password           string
	Schema             string
	MaxOpenConns       int
	MaxIdleConns       int
	ConnMaxLifetime    time.Duration
	ConnMaxIdleTime    time.Duration
	HealthCheckTimeout time.Duration
	IsDebug            bool
}

func dataSourceName(config SQLConfig, mysqlMultiStatements bool) string {
	var dsnCfg map[string]string
	format := ""

	switch config.Engine {
	case MySQL:
		dsnCfg = map[string]string{
			"charset":         "utf8",
			"parseTime":       "True",
			"loc":             "UTC",
			"multiStatements": strconv.FormatBool(mysqlMultiStatements),
			"autocommit":      "True",
		}
		format = "%s:%s@tcp(%s:%d)/%s?%s"
	case PostgreSQL:
		dsnCfg = map[string]string{
			"sslmode":  "disable",
			"TimeZone": "UTC",
		}
		format = "postgres://%s:%s@%s:%d/%s?%s"
	}

	dsnCfgList := make([]string, 0, len(dsnCfg))
	for key, val := range dsnCfg {
		dsnCfgList = append(dsnCfgList, fmt.Sprintf("%s=%s", key, val))
	}

	return fmt.Sprintf(
		format,
		config.Username,
		config.Password,
		config.Host,
		config.Port,
		config.Schema,
		strings.Join(dsnCfgList, "&"),
	)
}

func NewConnection(config SQLConfig, mysqlMultiStatements bool) (*sql.DB, func(), error) {
	dsn := dataSourceName(config, mysqlMultiStatements)

	driverName := config.Engine.String()

	conn, err := sql.Open(driverName, dsn)
	if err != nil {
		return nil, nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := conn.PingContext(ctx); err != nil {
		return nil, nil, err
	}

	if config.MaxOpenConns != 0 {
		conn.SetMaxOpenConns(config.MaxOpenConns)
	}
	if config.MaxIdleConns != 0 {
		conn.SetMaxIdleConns(config.MaxIdleConns)
	}
	if config.ConnMaxLifetime != 0 {
		conn.SetConnMaxLifetime(config.ConnMaxLifetime)
	}
	if config.ConnMaxIdleTime != 0 {
		conn.SetConnMaxIdleTime(config.ConnMaxIdleTime)
	}

	cleanup := func() {
		if err := conn.Close(); err != nil {
			log.Printf("failed to close %s connection %e", driverName, err)
		}
	}

	return conn, cleanup, nil
}

func NewGorm(
	dialect SQLEngine,
	conn *sql.DB,
	config SQLConfig,
) (*gorm.DB, func(), error) {
	var gormConfig = &gorm.Config{
		SkipDefaultTransaction: true,
	}
	if config.IsDebug {
		gormConfig.Logger = logger.New(
			log.New(os.Stdout, "\r\n", log.LstdFlags),
			logger.Config{
				SlowThreshold: time.Second,
				LogLevel:      logger.Silent,
				Colorful:      true,
			},
		)
	}
	var gormDialect gorm.Dialector
	switch dialect {
	case MySQL:
		gormDialect = mysql.New(mysql.Config{Conn: conn})
	case PostgreSQL:
		gormDialect = postgres.New(postgres.Config{Conn: conn})
	default:
		return nil, func() {}, errors.WithStack(errors.New("sql dialect is not available"))
	}

	gormDB, err := gorm.Open(gormDialect, gormConfig)
	if err != nil {
		return nil, nil, err
	}

	cleanup := func() {
		if err := conn.Close(); err != nil {
			log.Printf("failed to close %s connection %e", dialect, err)
		}
	}

	if config.IsDebug {
		return gormDB.Debug(), cleanup, nil
	}

	return gormDB, cleanup, nil
}
