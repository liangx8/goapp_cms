package utils_test

import (
	"crypto/sha1"
	"testing"

	"rcgreed.bid/ics/utils"
)

var pwdKit = utils.NewPasswordKit(sha1.New())

func TestPwdkit(t *testing.T) {
	salt := make([]byte, 5)
	utils.RandomSalt(salt)

	enc := pwdKit.Create(salt, "1")
	t.Logf("1:%p", enc)
	if !pwdKit.Verify(enc, "1") {
		t.Fatal("postive failed")
	}
	t.Logf("2:%p", enc)
	if pwdKit.Verify(enc, "9") {
		t.Fatal("nagetive failed")
	}
}
func TestSlice(t *testing.T) {
	sl1 := []int{1, 2, 3, 4}
	sl2 := sl1[:]
	t.Log(1, sl1)
	t.Log(2, sl2)
	sl2[3] = 99
	t.Log(1, sl1)
	t.Log(2, sl2)
	t.FailNow()
}
