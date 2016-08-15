package web

import (
	"net/http"
	"fmt"

	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
	"google.golang.org/cloud/storage"

	
	"gcs"
)
const (
	bucketName="pfa.rc-greed.com"
)

func hello(w http.ResponseWriter,r *http.Request){
	fmt.Fprintln(w,"connect to google clould storage...")
	ctx := appengine.NewContext(r)
	
	client, err:=storage.NewClient(ctx)
	if err != nil {
		log.Errorf(ctx,"%v",err)
		return
	}
	defer func(){
		if err := client.Close(); err != nil {
			log.Errorf(ctx,"%v",err)
		}
	}()
	bucket := client.Bucket(bucketName)
	each := gcs.All(ctx,bucket)
	err=each(func(obj *storage.ObjectHandle)error{
		attrs,err := obj.Attrs(ctx)
		if err != nil {
			return err
		}
		fmt.Fprintln(w,attrs.Name)
		return nil
	})
	if err != nil {
		log.Errorf(ctx,"%v",err)
	}
}
func init(){
	http.HandleFunc("/",hello)
}
