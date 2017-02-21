package expense

import (
	"golang.org/x/net/context"
	"cloud.google.com/go/storage"
	"google.golang.org/appengine/log"

	"gopkg.in/yaml.v2"

	
	"github.com/liangx8/gcloud-helper/gcs"


	"io/ioutil"
	"fmt"

)

type (
	ExpenseCloud struct {
		client *storage.Client
		Save func([]Expense)error
		Load func(*[]Expense)error
	}
)

func NewCloud(ctx context.Context,ac string) (*ExpenseCloud,error){
	cli,err := storage.NewClient(ctx)
	if err != nil {return nil,err }
	bucket := cli.Bucket(bucketName)
	filename := prefix + ac + ".yaml"
	return &ExpenseCloud{
		client : cli,
		Save:  func(es []Expense)error {
			oh := bucket.Object(filename)
			objw := oh.NewWriter(ctx)
			defer objw.Close()
			buf,err := yaml.Marshal(es)
			if err != nil { return err }
			_,err = objw.Write(buf)
			if err != nil { return err }
			return nil
		},
		Load: func(es *[]Expense)error{
			oh := bucket.Object(filename)
			objr,err := oh.NewReader(ctx)
			if err == storage.ErrObjectNotExist {
				return nil
			}
			if err != nil { return err }
			defer objr.Close()
			buf,err := ioutil.ReadAll(objr)
			if err != nil { return err }
			return yaml.Unmarshal(buf,es)
		},
	},nil
}
func (ec *ExpenseCloud)Close()error{
	return ec.client.Close()
}
func (ec *ExpenseCloud)InsertOrUpdate(ex Expense) error {
	exps:= make([]Expense,0,200)
	err :=ec.Load(&exps)
	if err != nil {
		return err
	}
	isInsert := true
	for i,e := range exps{
		if e.Seq == ex.Seq {
			exps[i]=ex
			isInsert = false
			break
		}
	}
	if isInsert {
		exps = append(exps,ex)
	}
	err = ec.Save(exps)
	if err != nil {
		return err
	}
	return nil
}
func (ec *ExpenseCloud)Delete(seq string) error{
	exps:= make([]Expense,0,200)
	dst:= make([]Expense,0,200)
	err :=ec.Load(&exps)
	if err != nil {
		return err
	}
	change :=false
	for _,e := range exps{
		if e.Seq == seq {
			change = true
			continue
		}
		dst = append(dst,e)
	}
	if change {
		return ec.Save(dst)
	}else {
		
		return fmt.Errorf("%s is not exists",seq)
	}
}
// TODO: work here
func (ec *ExpenseCloud)Merge(exps []Expense,result func(addCount,updateCount int, err error)){
	var org []Expense
	err := ec.Load(&org)
	if err != nil {
		result(0,0,err)
		return
	}
	x,a,u :=Merge(org,exps)
	ec.Save(x)
	result(a,u,nil)

}
func AllAccount(bucket *gcs.Bucket,one func (act string)) error{
	q := storage.Query{
		Delimiter:"/",
		Prefix:prefix,
	}
	e:=bucket.Objects(gcs.AttrCallback(func(attrs *storage.ObjectAttrs) error{
		aName := attrs.Name[8:]
		length := len(aName)
		one(aName[:length-5])
		return nil
	}),&q)
	if e != nil {
		return e
	}else {
		one("")
		return nil
	}
}


func OldData(ctx context.Context) []Expense{
	cli,err := storage.NewClient(ctx)
	if err != nil {
		log.Errorf(ctx,"%v",err)
		return nil
	}

	bucket := cli.Bucket(bucketName)
	filename := "expense/expense.yaml"
	oh := bucket.Object(filename)
	objr,err := oh.NewReader(ctx)
	if err != nil {
		log.Errorf(ctx,"%v",err)
		return nil
	}
	buf,err := ioutil.ReadAll(objr)
	if err != nil {
		log.Errorf(ctx,"%v",err)
		return nil
	}
	objs := make ([]map[string]interface{},0)
	err = yaml.Unmarshal(buf,&objs)
	if err != nil {
		log.Errorf(ctx,"%v",err)
		return nil
	}
	exps :=make([]Expense,len(objs))
	for i,obj:=range objs{
		tags := make([]string,0,2)
		exps[i].Amount=obj["amount"].(int)
		exps[i].CountIn=obj["count-in"].(bool)
		exps[i].When=int64(obj["when"].(int))
		exps[i].CreatedTime=int64(obj["update"].(int))
		mils:=obj["miles"].(int)
		if(mils != 0){
			exps[i].Remark=fmt.Sprintf("里程:%d",mils)
		}
		s := obj["type"].(string)
		if s != "" {
			tags = append(tags,s)
		}
		s = obj["sub-type"].(string)
		if s != "" {
			tags = append(tags,s)
		}
		if len(tags)>0{
			exps[i].Tags=tags
		}
		exps[i].Remark=exps[i].Remark + obj["remark"].(string)
		
	}
	return exps
}
const (
//	projectId= "personal-financial-140007"
	bucketName="pfa.rc-greed.com"
	prefix = "expense2/"
)
