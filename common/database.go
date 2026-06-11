package common

import (
	"log"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var (
	UsingSQLite    bool
	UsingMySQL     bool
	UsingPostgreSQL bool
)

// InitDatabase initializes the database connection
func InitDatabase() (*gorm.DB, error) {
	dsn := GetEnvOrDefault("SQL_DSN", "")
	if dsn == "" {
		// Default to SQLite
		UsingSQLite = true
		return gorm.Open(sqlite.Open("data.db"), &gorm.Config{})
	}

	// Determine database type from DSN
	switch {
	case len(dsn) > 6 && dsn[:6] == "mysql:":
		UsingMySQL = true
		return gorm.Open(mysql.Open(dsn), &gorm.Config{})
	case len(dsn) > 11 && dsn[:11] == "postgresql:":
		UsingPostgreSQL = true
		return gorm.Open(postgres.Open(dsn), &gorm.Config{})
	default:
		// Try SQLite as fallback
		UsingSQLite = true
		return gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	}
}

// GetEnvOrDefault returns the value of an environment variable or a default value
func GetEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// LogDatabaseType logs the database type being used
func LogDatabaseType() {
	switch {
	case UsingSQLite:
		log.Println("Using SQLite database")
	case UsingMySQL:
		log.Println("Using MySQL database")
	case UsingPostgreSQL:
		log.Println("Using PostgreSQL database")
	}
}
