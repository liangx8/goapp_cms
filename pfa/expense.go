package pfa

import (
	"time"
	"fmt"
)


type (
	Expense struct {
		Seq,Amount,Miles int
		CountIn bool `json:"count-in" yaml:"count-in"`
		Remark,Type string
		SubType string `json:"sub-type" yaml:"sub-type"`
		When,Update int64
		
	}
	order struct{
		es []Expense
		less func(l,r Expense) bool
	}
)

func whenOrder(l,r Expense) bool{
	if l.When == r.When {
		return l.Update < r.Update
	}
	return l.When < r.When
}
func updateOrder(l,r Expense) bool{
	return l.Update < r.Update
}

func (odr *order)Len() int{
	return len(odr.es)
}
func (odr *order)Less(i, j int)bool {
	oi,oj := odr.es[i],odr.es[j]
	return odr.less(oi,oj)
}
func (odr *order)Swap(i,j int){
	odr.es[i],odr.es[j] = odr.es[j],odr.es[i] 
}

var location *time.Location

func init(){
	location,_ = time.LoadLocation("Asia/Hong_Kong")
}
func JavaTimestampStr(ts int64) string{
	t:=time.Unix(ts/1000,0).In(location)
	return t.Format(TIMESTAMP)
	
}
func JavaDateStr(ts int64) string{
	t:=time.Unix(ts/1000,0).In(location)
	return t.Format(DATE)
	
}
func Money(i int) string {
	return fmt.Sprintf("%d.%02d",i /100, i % 100)
}


const (
	TIMESTAMP = "2006-01-02 15:04:05.000 -0700"
	DATE      = "2006-01-02"
)
