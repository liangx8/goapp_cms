package pfa

import (
	"time"
)


type (
	Expense struct {
		Seq,Amount,Miles int
		CountIn bool `json:"count-in" yaml:"count-in"`
		Remark,Type string
		SubType string `json:"sub-type" yaml:"sub-type"`
		When,Update int64
		
	}
)


func JavaTimestampStr(ts int64) string{
	t:=time.Unix(ts/1000,0)
	return t.Format(TIMESTAMP)
	
}
func JavaDateStr(ts int64) string{
	t:=time.Unix(ts/1000,0)
	return t.Format(DATE)
	
}

const (
	TIMESTAMP = "2006-01-02 15:04:05.000 -0700"
	DATE      = "2006-01-02"
)
