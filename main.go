package main

import (
	"log"
	"net/http"

	"rcgreed.bid/ics/ctrl"
	"rcgreed.bid/ics/install"
)

func main() {
	err := ctrl.LocaleInit("messages")
	if err != nil {
		log.Print(err)
	}
	_, err = install.LoadCfg()
	if err != nil {
		log.Print(err)

	}

	//log.Fatal(http.ListenAndServe(":8000", http.FileServer(http.Dir("web/"))))
	log.Println("Server start")
	log.Fatal(http.ListenAndServe(":8000", ctrl.Route()))
}
