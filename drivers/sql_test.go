// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package drivers

import "testing"

const (
	host = "127.0.0.1"
)

func TestSQLEngine_String(t *testing.T) {
	if MySQL.String() != "mysql" {
		t.Errorf("expected 'mysql', got %v", MySQL.String())
	}
	if PostgreSQL.String() != "postgres" {
		t.Errorf("expected 'postgres', got %v", PostgreSQL)
	}
}

func TestDataSourceName(t *testing.T) {
	t.Run("MySQL DataSource Name", func(t *testing.T) {
		cfg := SQLConfig{
			Engine:   MySQL,
			Host:     "localhost",
			Port:     3306,
			Username: "root",
			Password: "secret",
			Schema:   "test_db",
		}
		dsn := dataSourceName(cfg, true)
		if dsn == "" {
			t.Errorf("expected valid MySQL DSN, got empty string")
		}
	})

	t.Run("PostgreSQL DataSource Name", func(t *testing.T) {
		cfg := SQLConfig{
			Engine:   PostgreSQL,
			Host:     "localhost",
			Port:     5432,
			Username: "postgres",
			Password: "secret",
			Schema:   "test_db",
		}
		dsn := dataSourceName(cfg, false)
		if dsn == "" {
			t.Errorf("expected valid PostgreSQL DSN, got empty string")
		}
	})
}

func TestNewConnection_InvalidConfig(t *testing.T) {
	t.Run("Connection failure with invalid port", func(t *testing.T) {
		cfg := SQLConfig{
			Engine:   MySQL,
			Host:     host,
			Port:     1,
			Username: "root",
			Password: "wrong",
			Schema:   "invalid",
		}

		conn, cleanup, err := NewConnection(cfg, false)
		if err == nil {
			t.Errorf("expected connection error, got nil")
			if cleanup != nil {
				cleanup()
			}
			if conn != nil {
				err = conn.Close()
				if err != nil {
					t.Errorf("failed closing connection, got %v", err)
				}
			}
		} else if conn != nil {
			t.Errorf("expected nil connection on error, got %v", conn)
		}
	})
}

func TestNewGorm_InvalidDialect(t *testing.T) {
	t.Run("Gorm initialization failure with unsupported dialect", func(t *testing.T) {
		cfg := SQLConfig{IsDebug: false}
		db, cleanup, err := NewGorm(SQLEngine("oracle"), nil, cfg)
		if err == nil {
			t.Error("expected error for unsupported dialect, got nil")
			if cleanup != nil {
				cleanup()
			}
		}
		if db != nil {
			t.Errorf("expected nil gorm DB on error, got %v", db)
		}
	})
}
