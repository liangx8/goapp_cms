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

type ErrorResult struct{
	Title,Head,Message string
	Err error
}

func DefaultErrorResut(title string,err error) ErrorResult{
	return ErrorResult{
		Title:title,
		Head:title,
		Message:func() string{if err==nil{return ""} else {return err.Error()}}(),
		Err:err,
	}
}

var tempError = template.Must(template.New("").Parse(error_tmpl))

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
			aName := attrs.Name[8:]
			length := len(aName)
			account = append(account,aName[:length-5])
			return nil
		}),&q)
		if e!= nil {
			log.Errorf(ctx,"%v",err)
		}
		tmpl,er := buildTemplate(ctx,bucket,"tmpl/account_list.tmpl")
		if er != nil {
			tmpl = template.Must(template.New("").Parse(account_list_tmpl))
		}
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
		account:=r.FormValue("name")
		if account == ""{
			tempError.Execute(w,DefaultErrorResut("页面没有找到",nil))
			return
		}
		es := make([]Expense,0,50)
		dao,err := NewDao(ctx,account)
		if err == nil {
			err=dao.Load(&es)
		}
		if err != nil {
			tempError.Execute(w,DefaultErrorResut("出错"+account,err))
			return
		}
		defer dao.Close()
		tmpl,er:=buildTemplate(ctx,bucket,"tmpl/account.tmpl")
		if er != nil {
			tempError.Execute(w,er)
			return
		}
		model :=make(map[string]interface{})
		initModel(model)
		model["accountName"]=account;
		model["data"]=es
		tmpl.Execute(w,model)
	})
	if err != nil {
		log.Errorf(ctx,"%v",err)
		return
	}
}
func initModel(data map[string]interface{}){
	data["toDate"]=JavaDateStr
	data["toTimestamp"]=JavaTimestampStr
	data["money"]=Money
}
func buildTemplate(ctx context.Context,bkt *gcs.Bucket,name string) (*template.Template,*ErrorResult){
//		tmplObj:=bkt.Object("tmpl/account_list.tmpl")
	tmplObj:=bkt.Object(name)
	objr,err := tmplObj.NewReader(ctx)
	if err == storage.ErrObjectNotExist {
		return nil, &ErrorResult{
			Title   :"找不到模板",
			Head    :"你必须定义一个模板",
			Message :"必须定义一个模板"+name+"以便于显示内容",
			Err     :err,
		}
	} else {
		if buf,err := ioutil.ReadAll(objr); err == nil {
			tmpl, err := template.New("").Parse(string(buf))
			if err == nil { return tmpl,nil }
		}
		return nil, &ErrorResult{
			Title   :"模板错误",
			Head    :"模板错误",
			Message :"读取模板时错误",
			Err     :err,
		}
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
	error_tmpl=`<html>
<head><title>Error: {{.Title}}</title>
<style type="text/css">
html, body {
	font-family: "Roboto", sans-serif;
	color: #333333;
	background-color: #ea5343;
	margin: 0px;
}
h1 {
	color: #d04526;
	background-color: #ffffff;
	padding: 20px;
	border-bottom: 1px dashed #2b3848;
}
pre {
	margin: 20px;
	padding: 20px;
	border: 2px solid #2b3848;
	background-color: #ffffff;
}
</style>
</head><body>
<h1>Error</h1>
<pre style="font-weight: bold;">{{.Head}}</pre>
<pre>{{.Message}}</pre>
{{if .Err}}
<pre>{{.Err}}</pre>
{{end}}
</body>
</html>`

)
