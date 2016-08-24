package web

import (

	"golang.org/x/net/context"
	
	"github.com/liangx8/spark"
	"google.golang.org/appengine/log"
	yaml "gopkg.in/yaml.v2"
)

func Notification(ctx context.Context){
	w,r,_ := spark.ReadHttpContext(ctx)
	log.Infof(ctx,"bucket is change! %v",r)
	LoadConfig(ctx)
	b,err := yaml.Marshal(defaultConfig)
	if err != nil {
		log.Errorf(ctx,"%v",err)
	}
	w.Write(b)
	
	
}
