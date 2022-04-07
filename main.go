package main

import (
	"log"
	"net/http"

	"rcgreed.bid/ics/lite"
	"rcgreed.bid/ics/mgr"
	"rcgreed.bid/ics/view"
)

func ok(r *http.Request) any {
	return "hello world"
}
func main() {
	dbi, err := lite.NewDBI("")
	if err != nil {
		log.Fatal(err)
		return
	}
	defer dbi.Close()

	//log.Fatal(http.ListenAndServe(":8000", http.FileServer(http.Dir("web/"))))
	log.Println("Server start")
	log.Fatal(http.ListenAndServe(":8000", mgr.Filter(ok, view.JsonView())))
}
