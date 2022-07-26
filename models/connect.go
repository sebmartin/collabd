package models

import (
	"log"
	"os"
	"path"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func Connect() (*gorm.DB, error) {
	return ConnectWithDB(sqlite.Open("db/models.sqlite"))
}

func ConnectWithSqlLite(dbname string) (*gorm.DB, error) {
	return ConnectWithDB(sqlite.Open(dbname))
}

// Helper for testing
func ConnectWithTestDB() (*gorm.DB, func()) {
	dbtmpdir, _ := os.MkdirTemp("", "collabd_tests_")
	dbtmppath := path.Join(dbtmpdir, "_tests.sqlite")
	db, err := ConnectWithSqlLite(dbtmppath)
	if err != nil {
		log.Fatalf("Failed to create test database: %s", err)
	}
	sqlDB, _ := db.DB()

	return db, func() {
		sqlDB.Close()
		os.Remove(dbtmppath)
		os.Remove(dbtmpdir)
	}
}

func ConnectWithDB(dialector gorm.Dialector) (*gorm.DB, error) {
	db, err := gorm.Open(dialector, &gorm.Config{})
	if err != nil {
		return nil, err
	}

	db.AutoMigrate(
		&Player{},
		&Session{},
	)

	return db, err
}
