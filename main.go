package main

import (
	"log"
	"net/http"
	"reflect"

	"rcgreed.bid/ics/ctrl"
	"rcgreed.bid/ics/entity"
	"rcgreed.bid/ics/lite"
	"rcgreed.bid/ics/utils"
)

func main() {
	dbi, err := lite.NewDBI("")
	if err != nil {
		log.Fatal(err)
		return
	}
	defer dbi.Close()
	test()

	//log.Fatal(http.ListenAndServe(":8000", http.FileServer(http.Dir("web/"))))
	log.Println("Server start")
	log.Fatal(http.ListenAndServe(":8000", ctrl.Route()))
}
func test() {
	for ii := 0; ii < 10; ii++ {
		log.Print(utils.MakeID())
	}
	obj := reflect.TypeOf((*entity.User)(nil))
	ent := obj.Elem()
	log.Println(ent)
	for x := 0; x < ent.NumField(); x++ {
		log.Println(ent.Field(x))
	}
}
