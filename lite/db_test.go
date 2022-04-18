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
	if err := dbi.Init(entity.User{}); err != nil {
		t.Fatal(err)
	}
}
func TestDBLoad(t *testing.T) {
	user := entity.User{Seq: 1}
	dbi, err := lite.NewDBI("home.db")
	if err != nil {
		t.Fatal(err)
	}
	if err := dbi.Load(&user); err != nil {
		t.Fatal(err)
	}

}
