package web

import (
	"fmt"
//	"sync"

	"golang.org/x/net/context"
	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
	"google.golang.org/appengine/file"
	
	"github.com/liangx8/gcloud-helper/gcs"
//	"github.com/liangx8/spark/session"
	"github.com/liangx8/spark"

//	"pfa"
	"cloud.google.com/go/storage"
)


func index(ctx context.Context){
	w,r,err:=spark.ReadHttpContext(ctx)
	if err != nil {
		log.Errorf(ctx,"%v",err)
		return
	}
	
	bucket,err := gcs.NewBucket(ctx,projectId,bucketName)
	if err != nil {
		fmt.Fprintln(w,err)
		log.Errorf(ctx,"%v",err)
		return
	}
	defer bucket.Close()
	x,err := file.DefaultBucketName(ctx)
	if err != nil {
		fmt.Fprintln(w,err)
		return
	}
	fmt.Fprintln(w,x)
	err = bucket.Objects(gcs.AttrCallback(func(attrs *storage.ObjectAttrs)error{
		fmt.Fprintln(w,attrs.Name)
		return nil
	}),nil)
	if err != nil {
		fmt.Fprintf(w,"read objects error:%v",err)
		log.Errorf(ctx,"%v",err)
		return
	}
	fmt.Fprintf(w,"链接:%s\n",r.URL.RequestURI())
}

func init(){
	spk:=spark.New(appengine.NewContext)
	spk.HandleFunc("/",index)
}

const (
	projectId= "personal-financial-140007"
	bucketName="pfa.rc-greed.com"
)
