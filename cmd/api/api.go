package api

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/ErdajtSopjani/LikeMe_API/pkg/config"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var app config.AppConfig
var db *gorm.DB

func init() {
	app.IsProd = false

	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file", err)
	}
}

func initDB() (*gorm.DB, error) {
	/* retrieve db variables from env */
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	dbSSLMode := os.Getenv("DB_SSLMODE")

	/* create db connection string */
	db, err := gorm.Open(
		postgres.Open(
			fmt.Sprintf(
				"host=%s port=%s user=%s dbname=%s password=%s sslmode=%s",
				dbHost, dbPort, dbUser, dbName, dbPassword, dbSSLMode)),
		&gorm.Config{})

	return db, err
}

func Run() {
	port := ":3333"
	host := "localhost"
	addr := host + port

	/*
		errs := make(chan error)
		go func() {
		   	c := make(chan os.Signal, 1)
		   	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		   	errs <- fmt.Errorf("%s", <-c)
		}()
	*/

	log.Printf("Server running on %s\n", addr)

	var err error
	db, err := initDB()
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}
	log.Printf("Database connection established")

	err = http.ListenAndServe(addr, router(db))
	log.Fatal(err)

}

func router(db *gorm.DB) http.Handler {
	r := mux.NewRouter()

	r.HandleFunc("/api/v1/users", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, %q", "users")
	})

	return r
}
