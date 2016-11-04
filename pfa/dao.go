package pfa

import (
	"golang.org/x/net/context"
	"gopkg.in/yaml.v2"
	"cloud.google.com/go/storage"
	"google.golang.org/appengine/log"

)

type Dao struct{
	ctx context.Context
	client *storage.Client
	b *storage.BucketHandle
	account string
}
func NewDao(ctx context.Context,a string) (*Dao,error){
	cli,err := storage.NewClient(ctx)
	if err != nil { return nil,err }
	bucket := cli.Bucket(bucketName)
	return &Dao{ctx:ctx,client:cli,account:a,b:bucket},nil
}

func (dao *Dao)Save(es []Expense)error{
	b,err := yaml.Marshal(es)
	if err != nil { return err }
	fn := prefix + dao.account + ".yaml"
	log.Errorf(dao.ctx,fn)
	oh := dao.b.Object(fn)
	w:=oh.NewWriter(dao.ctx)
	_,err = w.Write(b)
	if err != nil { return err }
	err=w.Close()
	if err != nil { return err }
	return nil
}
func (dao *Dao)Close()error{
	return dao.client.Close()
}
const (
	prefix      ="expense/"
	bucketName = "pfa.rc-greed.com"
)
