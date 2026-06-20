// Package database opens and configures the GORM/SQLite connection.
package database

import (
	"fmt"
	"log"
	"os"
	"time"

	"restapi/internal/book"
	"restapi/internal/config"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

// Open connects to SQLite, tunes the pool, pings, and runs migrations.
func Open(cfg config.Config) (*gorm.DB, error) {
	gl := gormlogger.New(
		log.New(os.Stdout, "", log.LstdFlags),
		gormlogger.Config{
			LogLevel:                  gormlogger.Warn,
			IgnoreRecordNotFoundError: true, // not-found is a normal API outcome, not an error
		},
	)

	db, err := gorm.Open(sqlite.Open(cfg.DBPath), &gorm.Config{Logger: gl})
	if err != nil {
		return nil, fmt.Errorf("open database: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("get sql.DB: %w", err)
	}
	// SQLite is a single-writer engine; cap open connections to avoid
	// "database is locked" errors under concurrency.
	sqlDB.SetMaxOpenConns(1)
	sqlDB.SetMaxIdleConns(1)
	sqlDB.SetConnMaxLifetime(time.Hour)

	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("ping database: %w", err)
	}

	if err := db.AutoMigrate(&book.Book{}); err != nil {
		return nil, fmt.Errorf("migrate: %w", err)
	}

	return db, nil
}

// Close releases the underlying connection pool.
func Close(db *gorm.DB) error {
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}
