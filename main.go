package main

import (
	"log"
	"net/http"

	"rcgreed.bid/ics/ctrl"
	"rcgreed.bid/ics/lite"
)

func main() {
	ctrl.LocaleInit("messages")
	dbi, err := lite.NewDBI("")
	if err != nil {
		log.Fatal(err)
		return
	}
	defer dbi.Close()

	//log.Fatal(http.ListenAndServe(":8000", http.FileServer(http.Dir("web/"))))
	log.Println("Server start")
	log.Fatal(http.ListenAndServe(":8000", ctrl.Route()))
}
