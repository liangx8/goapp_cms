package expense

import (
	"github.com/liangx8/spark/helper"

	"time"
	"fmt"
	"encoding/base64"
    "crypto/md5"
	"strings"
	"sort"
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
func CompleteExpense(exp *Expense) {
	if exp.CreatedTime==0 {
		exp.CreatedTime=JavaTimestampIntNow()
	}
	if exp.Seq=="" {
		exp.Seq=UniqueId()
	}
}
func UniqueId() string{
    enc := base64.NewEncoding("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_/")
    h:=md5.New()
    h.Write([]byte(time.Now().Format(TIMESTAMP)))
	return enc.EncodeToString(h.Sum(nil))[:20]
}

func cmpSeq(l,r interface{})bool{
	return strings.Compare(l.(Expense).Seq,r.(Expense).Seq) < 0
}

func Merge(src,dst []Expense)([]Expense,int,int){
	srt := helper.NewSorter(src,cmpSeq)
	sort.Sort(srt)
	srt = helper.NewSorter(dst,cmpSeq)
	sort.Sort(srt)
	idxdst := 0
	updateCount :=0;
	addition :=make([]Expense,0,20)
outer:
	for idxsrc,v := range src{
		if idxdst >= len(dst) {
			break
		}
		x := strings.Compare(v.Seq,dst[idxdst].Seq)
		for x > 0 {
			addition = append(addition,dst[idxdst])
			idxdst ++
			if idxdst >= len(dst) {
				break outer
			}
			x = strings.Compare(v.Seq,dst[idxdst].Seq)
		}
		if x==0 {
			src[idxsrc]=dst[idxdst]
			idxdst ++
			updateCount ++
			continue
		}
		if x <0 {
			continue
		}
	}
	if idxdst < len(dst){
		addition = append(addition,dst[idxdst:]...)
	}
	return append(src,addition...),len(addition),updateCount

}
const (
	TIMESTAMP = "2006-01-02 15:04:05.000 -0700"
	DATE      = "2006-01-02"

)
