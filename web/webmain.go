package web

import (
	"fmt"
	"sync"

	"golang.org/x/net/context"
	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
	
	"github.com/liangx8/gcloud-helper/gcs"
//	"github.com/liangx8/spark/session"
	"github.com/liangx8/spark"

	"pfa"
	"cloud.google.com/go/storage"
)

func index(ctx context.Context){
	w,r,err:=spark.ReadHttpContext(ctx)
	if err != nil {
		log.Errorf(ctx,"%v",err)
		return
	}
	
	bucket,err := gcs.NewBucket(ctx,pfa.ProjectId,"pfa.rc-greed.com")
	if err != nil {
		fmt.Fprintln(w,err)
		return
	}
	defer bucket.Close()
	err = bucket.Objects(gcs.AttrCallback(func(attrs *storage.ObjectAttrs)error{
		fmt.Fprintln(w,attrs.Name)
		return nil
	}),nil)
	if err != nil {
		fmt.Fprint(w,err)
		return
	}
	fmt.Fprintf(w,"链接:%s\n",r.URL.RequestURI())
}

func init(){
	spk:=spark.New(appengine.NewContext)
	spk.HandleFunc("/pfa/account",pfa.Account)
	spk.HandleFunc("/pfa/",pfa.PFA)
	spk.HandleFunc("/",index)
}
var once sync.Once
