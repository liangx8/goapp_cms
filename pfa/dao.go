package pfa

import (
	"golang.org/x/net/context"
	"gopkg.in/yaml.v2"
	"cloud.google.com/go/storage"
	"google.golang.org/appengine/log"

	"io/ioutil"
	"sort"
)

type Dao struct{
	client *storage.Client
	Account func() string
	Save func([]Expense) error
	Load func(*[]Expense) error
}
func NewDao(ctx context.Context,a string) (*Dao,error){
	cli,err := storage.NewClient(ctx)
	if err != nil { return nil,err }
	bucket := cli.Bucket(bucketName)
	fn := prefix + a + ".yaml"
	return &Dao{
		client  :cli,
		Account :func()string{
			return a
		},
		Save    : func(es []Expense)error{
			buf,err := yaml.Marshal(es)
			if err != nil { return err }
			log.Errorf(ctx,fn)
			oh := bucket.Object(fn)
			w:=oh.NewWriter(ctx)
			_,err = w.Write(buf)
			if err != nil { return err }
			err=w.Close()
			if err != nil { return err }
			return nil
		},
		Load    :func(es *[]Expense) error{
			oh := bucket.Object(fn)
			objr,err := oh.NewReader(ctx)
			if err != nil { return err }
			buf,err :=ioutil.ReadAll(objr)
			if err != nil { return err }
			err = yaml.Unmarshal(buf,es)
			return err
		},
	},nil
}

func (dao *Dao)Close()error{
	return dao.client.Close()
}
// msg 中的 countUpdate 更新的项目,countAdd 新增的项目
// 为 Expense 实现 sort.Interface方法，然后保存数据和已存在数据之间的合并处理
func (dao *Dao)Merge(newData []Expense,msg func(countUpdate,countAdd int)) error{
	old := make([]Expense,0)
	addExpense := make([]Expense, 0, 10)
	err := dao.Load(&old)
	if err != nil {
		return err
	}
	odr := &order{es:old,less:updateOrder}
	sort.Sort(odr)
	odr.es=newData
	sort.Sort(odr)
	oldIdx := 0
	newIdx := 0
	cntUpdate := 0
	for {
		if old[oldIdx].Update == newData[newIdx].Update {
			old[oldIdx]=newData[newIdx]
			cntUpdate ++
			oldIdx ++
			if oldIdx == len(old) { break }
			newIdx ++
			if newIdx == len(newData) { break }
			continue
		}
		if old[oldIdx].Update < newData[newIdx].Update {
			oldIdx ++
			if oldIdx == len(old) { break }
			continue
		}
		// old one > new one
		addExpense = append(addExpense,newData[newIdx])
		newIdx ++
		if newIdx == len(newData) { break }		
	}
	for i := newIdx; i< len(newData) ; i++ {
		addExpense = append(addExpense,newData[i])
	}
	old = append(old, addExpense...)
	err = dao.Save(old)
	if err != nil {
		return err
	}
	msg(cntUpdate,len(addExpense))
	return nil
}
const (
	prefix      ="expense/"
	bucketName = "pfa.rc-greed.com"
)
