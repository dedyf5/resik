// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package mysql

import (
	"errors"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/dedyf5/resik/config"
	checkEntity "github.com/dedyf5/resik/entities/check"
)

func TestCheckRepo_Check(t *testing.T) {
	t.Run("Return StatusUp when database ping succeeds", func(t *testing.T) {
		db, mock, err := sqlmock.New(sqlmock.MonitorPingsOption(true))
		if err != nil {
			t.Fatalf("failed to create sqlmock: %v", err)
		}
		mock.ExpectPing()
		mock.ExpectClose()
		t.Cleanup(func() {
			if err := db.Close(); err != nil {
				t.Errorf("failed to close mock db: %v", err)
			}
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("unmet sqlmock expectations: %v", err)
			}
		})

		cfg := config.Config{}
		cfg.Database.HealthCheckTimeout = 2 * time.Second

		repo := New(db, cfg)
		got := repo.Check()

		if got.Name != "database" {
			t.Errorf("Name = %q, want %q", got.Name, "database")
		}
		if got.Status != checkEntity.StatusUp {
			t.Errorf("Status = %v, want %v", got.Status, checkEntity.StatusUp)
		}
		if got.Error != nil {
			t.Errorf("Error = %v, want nil", got.Error)
		}
	})

	t.Run("Return StatusDown when database ping fails", func(t *testing.T) {
		db, mock, err := sqlmock.New(sqlmock.MonitorPingsOption(true))
		if err != nil {
			t.Fatalf("failed to create sqlmock: %v", err)
		}

		wantErr := errors.New("connection refused")
		mock.ExpectPing().WillReturnError(wantErr)
		mock.ExpectClose()
		t.Cleanup(func() {
			if err := db.Close(); err != nil {
				t.Errorf("failed to close mock db: %v", err)
			}
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("unmet sqlmock expectations: %v", err)
			}
		})

		cfg := config.Config{}
		cfg.Database.HealthCheckTimeout = 2 * time.Second

		repo := New(db, cfg)
		got := repo.Check()

		if got.Name != "database" {
			t.Errorf("Name = %q, want %q", got.Name, "database")
		}
		if got.Status != checkEntity.StatusDown {
			t.Errorf("Status = %v, want %v", got.Status, checkEntity.StatusDown)
		}
		if got.Error == nil {
			t.Fatal("expected non-nil Error, got nil")
		}
	})
}
