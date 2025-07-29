package testutil

import (
	"crypto-indicator-dashboard/pkg/logger"
	"database/sql"
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
)

// TestDB provides a test database instance for tests
type TestDB struct {
	DB     *gorm.DB
	sqlDB  *sql.DB
	Logger logger.Logger
}

// NewTestDB creates a new test database instance
func NewTestDB(t *testing.T) *TestDB {
	t.Helper()

	// Create in-memory SQLite database for testing
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
		Logger: gormLogger.Default.LogMode(gormLogger.Silent),
	})
	require.NoError(t, err, "Failed to connect to test database")

	sqlDB, err := db.DB()
	require.NoError(t, err, "Failed to get underlying sql.DB")

	testLogger := logger.New("test")

	return &TestDB{
		DB:     db,
		sqlDB:  sqlDB,
		Logger: testLogger,
	}
}

// NewTestDBWithPostgres creates a test database using PostgreSQL (for integration tests)
func NewTestDBWithPostgres(t *testing.T) *TestDB {
	t.Helper()

	// Get test database URL from environment
	testDSN := os.Getenv("TEST_DATABASE_URL")
	if testDSN == "" {
		t.Skip("TEST_DATABASE_URL not set, skipping PostgreSQL integration test")
	}

	db, err := gorm.Open(sqlite.Open(testDSN), &gorm.Config{
		Logger: gormLogger.Default.LogMode(gormLogger.Silent),
	})
	require.NoError(t, err, "Failed to connect to test PostgreSQL database")

	sqlDB, err := db.DB()
	require.NoError(t, err, "Failed to get underlying sql.DB")

	testLogger := logger.New("test")

	return &TestDB{
		DB:     db,
		sqlDB:  sqlDB,
		Logger: testLogger,
	}
}

// Migrate runs database migrations for testing
func (tdb *TestDB) Migrate(models ...interface{}) error {
	return tdb.DB.AutoMigrate(models...)
}

// Cleanup cleans up the test database
func (tdb *TestDB) Cleanup() error {
	return tdb.sqlDB.Close()
}

// Truncate truncates all tables in the test database
func (tdb *TestDB) Truncate(tables ...string) error {
	for _, table := range tables {
		if err := tdb.DB.Exec(fmt.Sprintf("DELETE FROM %s", table)).Error; err != nil {
			return err
		}
	}
	return nil
}

// Transaction executes a function within a database transaction
func (tdb *TestDB) Transaction(fn func(tx *gorm.DB) error) error {
	return tdb.DB.Transaction(fn)
}

// WithTx returns a new TestDB instance with a transaction
func (tdb *TestDB) WithTx(t *testing.T) *TestDB {
	t.Helper()

	tx := tdb.DB.Begin()
	require.NoError(t, tx.Error, "Failed to begin transaction")

	t.Cleanup(func() {
		tx.Rollback()
	})

	sqlTx, err := tx.DB()
	require.NoError(t, err, "Failed to get underlying sql.DB from transaction")

	return &TestDB{
		DB:     tx,
		sqlDB:  sqlTx,
		Logger: tdb.Logger,
	}
}