package expense2
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
		Rseq string `json:"rseq" yaml:"rseq"`
	}
	Relation struct{
		Seq string `yaml:"seq"`
		Remark string `yaml:"seq"`
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
func UniqueId() string{
    enc := base64.NewEncoding("/_0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
    h:=md5.New()
    h.Write([]byte(time.Now().Format(TIMESTAMP)))
	return enc.EncodeToString(h.Sum(nil))[:20]
}
const (
	TIMESTAMP = "2006-01-02 15:04:05.000 -0700"
	DATE      = "2006-01-02"

)
