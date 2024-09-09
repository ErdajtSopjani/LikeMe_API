package tests

import (
	"fmt"
	"log"
	"os"

	"github.com/ErdajtSopjani/LikeMe_API/internal/handlers"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// SetupTestDB sets up a connection to the test database
func SetupTestDB() *gorm.DB {
	// get test-db variables from env
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_TEST_NAME")

	// create db connection string
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s dbname=%s password=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbName, dbPassword,
	)

	// connect to the test database
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to the test database: %v", err)
	}

	// create tables in the test database
	err = runMigrations(db)
	if err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	// return test DB connection
	return db
}

func runMigrations(db *gorm.DB) error {
	// use the already defined table models on the test database
	return db.AutoMigrate(
		&handlers.User{},
		&handlers.UserProfile{},
		&handlers.UserToken{},
		&handlers.Tag{},
		&handlers.UserInterest{},
		&handlers.Follow{},
		&handlers.BlockedUser{},
		&handlers.Post{},
		&handlers.Like{},
		&handlers.Comment{},
		&handlers.Message{},
		&handlers.TwoFactor{},
		&handlers.VerificationToken{},
	)
}

// CleanupTestDB clears the database after each test if needed
func CleanupTestDB(db *gorm.DB) {
	tables := []string{
		"users",
		"user_profiles",
		"user_tokens",
		"tags",
		"user_interests",
		"follows",
		"blocked_users",
		"posts",
		"likes",
		"comments",
		"messages",
		"two_factors",
		"verification_tokens",
	}

	// Truncate all tables and reset the primary key sequences
	for _, table := range tables {
		err := db.Exec("TRUNCATE TABLE " + table + " RESTART IDENTITY CASCADE").Error
		if err != nil {
			log.Fatalf("Failed to truncate table %s: %v", table, err)
		}
	}
}
