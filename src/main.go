package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/ErdajtSopjani/LikeMe_API/src/pkg/config"
)

var app config.AppConfig

func main() {
	app.IsProd = false

	srv := &http.Server{
		Addr:    ":8080",
		Handler: routes(&app),
	}

	fmt.Println("Starting server on :8080")
	err := srv.ListenAndServe()
	log.Fatalln(err)
}
