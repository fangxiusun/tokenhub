package model

import (
	"log"
	"sync"

	"github.com/your-username/your-project/common"
	"gorm.io/gorm"
)

var (
	DB             *gorm.DB
	batchUpdates   chan func()
	batchUpdatesOnce sync.Once
)

// Init initializes the database
func Init() error {
	var err error
	DB, err = common.InitDatabase()
	if err != nil {
		return err
	}

	// Enable logging in debug mode
	if common.GetEnvOrDefault("GIN_MODE", "debug") == "debug" {
		DB = DB.Debug()
	}

	// Auto migrate tables
	err = DB.AutoMigrate(
		&User{},
		&Token{},
		&Channel{},
		&Ability{},
		&Log{},
		&Option{},
		&Pricing{},
	)
	if err != nil {
		return err
	}

	// Initialize batch updates
	initBatchUpdates()

	common.LogDatabaseType()
	log.Println("Database initialized successfully")
	return nil
}

// Close closes the database connection
func Close() {
	if DB != nil {
		sqlDB, err := DB.DB()
		if err == nil {
			sqlDB.Close()
		}
	}
}

// GetDB returns the database instance
func GetDB() *gorm.DB {
	return DB
}

func initBatchUpdates() {
	batchUpdatesOnce.Do(func() {
		batchUpdates = make(chan func(), 1000)
		go func() {
			for update := range batchUpdates {
				update()
			}
		}()
	})
}

// EnqueueUpdate adds an update to the batch queue
func EnqueueUpdate(update func()) {
	select {
	case batchUpdates <- update:
	default:
		log.Println("Warning: Batch update queue is full")
	}
}
