package api

import (
	"fmt"
	"log"
	"net/http"

	"github.com/ErdajtSopjani/LikeMe_API/pkg/config"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
)

var app config.AppConfig

func init() {
	app.IsProd = false

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

	go func() {
		log.Printf("Server running on %s\n", addr)
		//		db := getUserModelDB()
		//		defer db.Close()
		var db *gorm.DB
		err := http.ListenAndServe(addr, router(db))
		log.Fatalln(err)
	}()

}

func router(db *gorm.DB) http.Handler {
	r := mux.NewRouter()

	r.HandleFunc("/api/v1/users", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, %q", "users")
	})

	return r
}
