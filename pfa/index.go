package pfa

import (
	"net/http"

	"html/template"
	"io/ioutil"

	"golang.org/x/net/context"

	"google.golang.org/appengine/log"

	"cloud.google.com/go/storage"

	"github.com/liangx8/spark"
	"github.com/liangx8/gcloud-helper/gcs"
)

func index(ctx context.Context){
	bucket,err :=gcs.NewBucket(ctx,ProjectId,"pfa.rc-greed.com")
	if err != nil {
		log.Errorf(ctx,"%v",err)
		return
	}
	defer bucket.Close()
	err=spark.HttpHandleFunc(ctx,func(w http.ResponseWriter, r *http.Request){
		q := storage.Query{
			Delimiter:"/",
			Prefix:"expense/",
		}
		account := make([]string,0,10)
		e:=bucket.Objects(gcs.AttrCallback(func(attrs *storage.ObjectAttrs) error{
			account = append(account,attrs.Name[8:])
			return nil
		}),&q)
		if e!= nil {
			log.Errorf(ctx,"%v",err)
		}
		tmpl := buildTemplate(ctx,bucket,"tmpl/account_list.tmpl",account_list_tmpl)
		w.Header().Add("Content-Type","text/html;charset=UTF-8")
		tmpl.Execute(w,account)
	})
	if err != nil {
		log.Errorf(ctx,"%v",err)
	}
}
// 列出帐套中的内容
func Account(ctx context.Context){
	bucket,err :=gcs.NewBucket(ctx,ProjectId,"pfa.rc-greed.com")
	if err != nil {
		log.Errorf(ctx,"%v",err)
		return
	}
	defer bucket.Close()
	err = spark.HttpHandleFunc(ctx,func(w http.ResponseWriter, r *http.Request){
		buildTemplate(ctx,bucket,"tmpl/account.tmpl",no_tmpl).Execute(w,nil)
	})
	if err != nil {
		log.Errorf(ctx,"%v",err)
		return
	}

}
func buildTemplate(ctx context.Context,bkt *gcs.Bucket,name,defaultStr string) *template.Template{
//		tmplObj:=bkt.Object("tmpl/account_list.tmpl")
	tmplObj:=bkt.Object(name)
	objr,err := tmplObj.NewReader(ctx)
	if err == storage.ErrObjectNotExist {
		return template.Must(template.New("").Parse(defaultStr))
	} else {
		if buf,err := ioutil.ReadAll(objr); err == nil {
			tmpl, err := template.New("").Parse(string(buf))
			if err == nil { return tmpl }
		}
		return template.Must(template.New("").Parse(defaultStr))
	}
	panic("never reach here")
}
const (
	ProjectId="personal-financial-140007"
	account_list_tmpl = `<html>
<head>
<title>全部帐套</title>
</head>
<body>
<h3>全部帐套</h3>
<ul>
{{range .}}
<li><a href="account?name={{.}}">{{.}}</a></li>
{{end}}
</ul>
</body>
</html>`
	no_tmpl = "<html><head><title>你必须定义一个模板</title></head><body>必须定义一个模板以便于显示内容</body></html>"
)
