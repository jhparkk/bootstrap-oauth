package main

import (
	"log"
	"net/http"

	"jhpark.sinsiway.com/bootstrap-oauth/app"
)

func main() {
	m := app.MakeHandler("./test.db")
	defer m.Close()

	err := http.ListenAndServe(":3000", m)
	if err != nil {
		log.Println(err)
	}
}
