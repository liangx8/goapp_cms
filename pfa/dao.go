package pfa

import (
	"golang.org/x/net/context"
	"gopkg.in/yaml.v2"
	"cloud.google.com/go/storage"
	"google.golang.org/appengine/log"

	"io/ioutil"
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
const (
	prefix      ="expense/"
	bucketName = "pfa.rc-greed.com"
)
