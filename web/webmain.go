package web

import (
	"net/http"
	"fmt"

	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
	"cloud.google.com/go/storage"

	
	"github.com/liangx8/gcloud-helper/gcs"
	"github.com/liangx8/spark/session"
)
const (
	bucketName="pfa.rc-greed.com"
)

func hello(w http.ResponseWriter,r *http.Request){

	s:=session.Get(w,r)
	if !s.Get("",nil) {
		fmt.Fprintln(w,"Are you going to login?")
		return
	}
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
	err = gcs.All(ctx,bucket,gcs.AttrCallback(func(obj *storage.ObjectAttrs)error{
		fmt.Fprintln(w,obj.Name)
		return nil
	}))

	if err != nil {
		log.Errorf(ctx,"%v",err)
	}
	fmt.Fprintln(w,s.Id())

}

func login(w http.ResponseWriter,r *http.Request){
	s:=session.Get(w,r)
	t:=true
	s.Put("",&t)
	fmt.Fprintln(w,"OK")
}
func logout(w http.ResponseWriter,r *http.Request){
	s:=session.Get(w,r)
	t:=false
	s.Put("",&t)
	fmt.Fprintln(w,"logout")
}
func init(){
	SessionInit()
	http.HandleFunc("/",hello)
	http.HandleFunc("/login",login)
	http.HandleFunc("/logout",logout)
}
