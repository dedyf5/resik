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
	"strings"
	"time"

	"github.com/pkg/errors"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
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
	Engine          SQLEngine
	Host            string
	Port            int
	Username        string
	Password        string
	Schema          string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime int
	ConnMaxIdleTime int
	IsDebug         bool
}

func NewMySQLConnection(config SQLConfig) (*sql.DB, func(), error) {
	dsnCfgs := map[string]string{
		"charset":    "utf8",
		"parseTime":  "True",
		"loc":        "Asia%2FJakarta",
		"autocommit": "True",
	}
	dsnCfgStr := make([]string, 0)
	for key, val := range dsnCfgs {
		dsnCfgStr = append(dsnCfgStr, fmt.Sprintf("%s=%s", key, val))
	}

	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?%s",
		config.Username,
		config.Password,
		config.Host,
		config.Port,
		config.Schema,
		strings.Join(dsnCfgStr, "&"),
	)

	conn, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := conn.PingContext(ctx); err != nil {
		return nil, nil, err
	}

	conn.SetMaxOpenConns(30)
	if config.MaxOpenConns != 0 {
		conn.SetMaxOpenConns(config.MaxOpenConns)
	}

	conn.SetMaxIdleConns(5)
	if config.MaxIdleConns != 0 {
		conn.SetMaxIdleConns(config.MaxIdleConns)
	}

	conn.SetConnMaxLifetime(time.Second * 300)
	if config.ConnMaxLifetime != 0 {
		conn.SetConnMaxLifetime(time.Second * time.Duration(config.ConnMaxLifetime))
	}

	conn.SetConnMaxIdleTime(time.Second * 300)
	if config.ConnMaxIdleTime != 0 {
		conn.SetConnMaxIdleTime(time.Second * time.Duration(config.ConnMaxIdleTime))
	}

	go func() {
		ticker := time.NewTicker(time.Second * 10)
		defer ticker.Stop()
		for range ticker.C {
			if err := conn.Ping(); err != nil {
				os.Exit(0)
			}
		}
	}()

	cleanup := func() {
		if err := conn.Close(); err != nil {
			log.Printf("failed to close mysql connection %e", err)
		}
	}

	return conn, cleanup, nil
}

func NewGorm(
	dialect SQLEngine,
	conn *sql.DB,
) (*gorm.DB, func(), error) {
	var gormConfig = &gorm.Config{
		SkipDefaultTransaction: true,
	}
	var gormDialect gorm.Dialector
	if dialect.String() == "mysql" {
		gormDialect = mysql.New(mysql.Config{Conn: conn})
	} else {
		return nil, func() {}, errors.WithStack(errors.New("sql dialect is not available"))
	}

	gormDB, err := gorm.Open(gormDialect, gormConfig)
	if err != nil {
		return nil, nil, err
	}

	cleanup := func() {
		if err := conn.Close(); err != nil {
			log.Printf("failed to close mysql connection %e", err)
		}
	}

	return gormDB, cleanup, nil
}
