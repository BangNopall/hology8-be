package pgsql_conn_mock

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	// "gorm.io/gorm/logger"
)

func NewMockDB(t *testing.T) (*gorm.DB, func() error, sqlmock.Sqlmock) {
	mockDB, mock, err := sqlmock.New()

	if err != nil {
		t.Fatal(err)
	}

	db, err := gorm.Open(postgres.New(postgres.Config{
		Conn: mockDB,
	}), &gorm.Config{
		// Logger: logger.Default.LogMode(logger.Silent),
	})

	if err != nil {
		t.Fatal(err)
	}

	return db, mockDB.Close, mock
}
