package mgr_test

import (
	"crypto/sha1"
	"testing"

	"rcgreed.bid/ics/entity"
	"rcgreed.bid/ics/lite"
	"rcgreed.bid/ics/mgr"
	"rcgreed.bid/ics/utils"
)

var pwdKid = utils.NewPasswordKit(sha1.New())

func TestLogin(t *testing.T) {
	salt := make([]byte, 6)
	utils.RandomSalt(salt)
	dbi, err := lite.NewDBI("home.db")
	if err != nil {
		t.Fatal(err)
	}
	m := mgr.NewManager(dbi)
	u := entity.User{
		Name:     "admin",
		Pwd:      pwdKid.Create(salt, "1"),
		Password: "1",
		Active:   true,
	}
	m.Init(u)
	ok, err := m.Login("admin", "1")
	if err != nil {
		t.Fatal(err)
	}
	if !ok {
		t.Fatal("test login postive failed!")
	}
	ok, err = m.Login("admin", "aa")
	if err != nil {
		t.Fatal(err)
	}
	if ok {
		t.Fatal("test login nagetive failed")
	}
}
