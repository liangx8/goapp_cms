package expense

import (
	"time"
	"fmt"
	"encoding/base64"
    "crypto/md5"
)

type (
	Expense struct{
		Seq string `json:"seq" yaml:"seq"`
		Amount int `json:"amount" yaml:"amount"`
		CountIn bool `json:"count-in" yaml:"count-in"`
		Remark string `json:"remark" yaml:"remark"`
		Tags []string `json:"tags" yaml:"tags"`
		When int64 `json:"when" yaml: "when"`
		CreatedTime int64 `json:"created-time" yaml:"created-time"`
	}
)

var location *time.Location
func init(){
	location,_ = time.LoadLocation("Asia/Hong_Kong")
}

func JavaTimestampStr(ts int64) string{
	t:=time.Unix(ts/1000,ts%1000 * 1000000).In(location)
	return t.Format(TIMESTAMP)
	
}
func JavaDateStr(ts int64) string{
	t:=time.Unix(ts/1000,ts%1000 * 1000000).In(location)
	return t.Format(DATE)
	
}
func Money(i int) string {
	
	if i >= 0 {
		return fmt.Sprintf("%d.%02d",i /100, i % 100)
	}
	i = -i
	return fmt.Sprintf("-%d.%02d",i /100, i % 100)
}

func JavaTimestampIntNow() int64{
	return time.Now().UnixNano() / 1000000
}
func JavaDateInt(date string) int64{
	d,err:=time.Parse(DATE,date)
	if err != nil {
		panic(err)
	}
	return d.Unix()*1000
}
func NowDateStr() string{
	return time.Now().Format(DATE)
}
func guessExpense(exp *Expense) int{
	retval:=EDIT_EXP
	if exp.CreatedTime==0 {
		exp.CreatedTime=JavaTimestampIntNow()
		retval=NEW_EXP
	}
	if exp.Seq=="" {
		exp.Seq=UniqueId()
	}
	return retval
}
func UniqueId() string{
    enc := base64.NewEncoding("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_/")
    h:=md5.New()
    h.Write([]byte(time.Now().Format(TIMESTAMP)))
	return enc.EncodeToString(h.Sum(nil))[:20]
}
const (
	TIMESTAMP = "2006-01-02 15:04:05.000 -0700"
	DATE      = "2006-01-02"

	NEW_EXP int = 1
	EDIT_EXP int = 2
)
