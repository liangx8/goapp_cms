package expense

import (
	"time"
	"fmt"
)

type (
	Expense struct{
		Seq string
		Amount int
		CountIn bool `json:"count-in" yaml:"count-in"`
		Remark string
		Tags []string
		When int64
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
	return fmt.Sprintf("%d.%02d",i /100, i % 100)
}


const (
	TIMESTAMP = "2006-01-02 15:04:05.000 -0700"
	DATE      = "2006-01-02"
)
