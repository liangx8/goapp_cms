package dist

import (
	"fmt"
	"io/ioutil"
//	"cloud.google.com/go/storage"
	"golang.org/x/net/context"
	"google.golang.org/appengine/log"

	"github.com/liangx8/spark"

	"encoding/base64"
	"time"
	"unicode/utf8"
)
type (
	Entity struct{
		Ip string `yaml:"ip"`
		Timestamp time.Time `yaml:"timestamp"`
		Raw string `yaml:"raw,omitempty"`
	}
)
func Touch(ctx context.Context){
	_,r,_ := spark.ReadHttpContext(ctx)
	view:=spark.NewYAMLView()
	log.Infof(ctx,"incoming")
	if r.Method != "PUT" && r.Method != "POST" {
		view.Render(ctx,"有错误,PUT or POST")
		return
	}
	raw,err:=ioutil.ReadAll(r.Body)
	if err != nil{
		log.Errorf(ctx,"%v",err)
	}
	var rawText string
	if utf8.Valid(raw){
		rawText=string(raw)
	} else {
		rawText = base64.StdEncoding.EncodeToString(raw)
	}
	entity := Entity{
		Ip:r.RemoteAddr,
		Raw:rawText,
		Timestamp:time.Now(),
	}
	if err:=add(ctx,entity); err !=nil {
		log.Errorf(ctx,"%v",err)
	}
	view.Render(ctx,entity)
	log.Infof(ctx,r.RemoteAddr)
	log.Infof(ctx,fmt.Sprint(r.Header))
}
