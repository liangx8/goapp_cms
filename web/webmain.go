package web

import (
	"net/http"
	"time"
	"io/ioutil"
	"io"
	"sync"

	"golang.org/x/net/context"
	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
	"cloud.google.com/go/storage"

	yaml "gopkg.in/yaml.v2"

	
//	"github.com/liangx8/gcloud-helper/gcs"
	"github.com/liangx8/spark/session"
	"github.com/liangx8/spark"

)
const (
	bucketName="pfa.rc-greed.com"
)

func index(ctx context.Context){
	w,r,err := spark.ReadHttpContext(ctx)
	if err != nil { panic (err) }
	
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
	url := r.URL.Path
	if url == "/" {
		url="/index.html"
	}
	obj:=bucket.Object(defaultConfig.SourcesPrefix+url)
	rd,err := obj.NewReader(ctx)
	if err != nil {
		http.NotFound(w,r)
		return
	}
	w.Header().Set("Content-Type",rd.ContentType())
	io.Copy(w,rd)
}


func init(){
	spk:=spark.New(appengine.NewContext)
	spk.AddChain(session.CreateSessionChain(GaeSessionMaker))
	spk.AddChain(initReadConfig)
	
	
	spk.HandleFunc("/",index)
	spk.HandleFunc("/login",login)

	spk.HandleFunc("/admin",admin)
	spk.HandleFunc("/trigger",Notification)
	spk.HandleFunc("/reset",pageReset)
}
func LoadConfig(ctx context.Context){
	client,err:=storage.NewClient(ctx)
	if err != nil {
		log.Errorf(ctx,"%v",err)
		return
	}
	defer client.Close()
	bucket := client.Bucket(bucketName)
	cfgObj:=bucket.Object("config.yaml")
	r,err:=cfgObj.NewReader(ctx)
	if err != nil {
		log.Errorf(ctx,"%v",err)
		defaultConfig.ResetId=session.UniqueId()
		defaultConfig.ResetTimeout=time.Now().AddDate(0,0,2)
		saveConfig(ctx,&defaultConfig)
		return
	}
	cfgBuf ,err := ioutil.ReadAll(r)
	r.Close()
	if err != nil {
		
		log.Errorf(ctx,"%v",err)
		return
	}
	err = yaml.Unmarshal(cfgBuf,&defaultConfig)
	if err != nil {
		log.Errorf(ctx,"%v",err)
	}
}

func initReadConfig(ctx context.Context,chain spark.HandleFunc){
	once.Do(func(){
		log.Infof(ctx,"Service is funcational, Instance(%s)",appengine.InstanceID())
		LoadConfig(ctx)
		parseTemplate()
	})
	chain(ctx)
}

var once sync.Once
