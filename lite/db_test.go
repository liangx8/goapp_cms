package lite_test

import (
	"testing"

	"rcgreed.bid/ics/entity"
	"rcgreed.bid/ics/lite"
)

func TestDBInit(t *testing.T) {
	dbi, err := lite.NewDBI("home.db")
	if err != nil {
		t.Fatal(err)
	}
	if err := dbi.Init(); err != nil {
		t.Fatal(err)
	}
	user := entity.User{Seq: 1, Name: "admin", Password: "xx", Updated: 1, Active: true, Remark: "remark"}
	t.Log(lite.InsertSQL(user))
	t.Fail()
}
