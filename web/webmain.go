package web

import (
	"fmt"
	"html/template"
//	"sync"

	"golang.org/x/net/context"
	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
//	"google.golang.org/appengine/file"
	
	"github.com/liangx8/gcloud-helper/gcs"
//	"github.com/liangx8/spark/session"
	"github.com/liangx8/spark"

//	"pfa"
	"expense"
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
	accounts := make([]string,0,10)
	err = expense.AllAccount(bucket,func(name string){
		accounts=append(accounts,name)
	})
	if err != nil {
		fmt.Fprintf(w,"read objects error:%v",err)
		log.Errorf(ctx,"%v",err)
		return
	}

	data := make(map[string]interface{})
	data["account"]=accounts
	data["url"]=r.URL.RequestURI();
	page.Execute(w,data)
}

func init(){
	spk:=spark.New(appengine.NewContext)
	spk.HandleFunc("/pfa/",pfa)
	spk.HandleFunc("/",index)
	spk.HandleFunc("/edit",pfaEdit)
	spk.HandleFunc("/list",pfaList)
	page = template.Must(template.New("page").Parse(homeTemplate))
}
var page *template.Template
const (
	projectId= "personal-financial-140007"
	bucketName="pfa.rc-greed.com"
	homeTemplate=`<!DOCTYPE HTML>
<html>
<head>
	<title>Expense Account List{{.url}}</title>
</head>
<body>
<form action="list" method="POST">
新建帐套:<input name="account" /><input type="submit" />
</form>
{{range .account}}
{{.}}
{{end}}
</body>
</html>`

)
