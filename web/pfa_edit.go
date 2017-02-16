package web
import (
	"google.golang.org/appengine/log"
	
	"net/http"
	"golang.org/x/net/context"
	"html/template"
	"github.com/liangx8/spark"
	"expense"
	"strconv"
	"strings"
)
func pfaEdit(ctx context.Context){
	w,r,_ := spark.ReadHttpContext(ctx)
	if r.Method != "POST" {
		panic("必须是PUT或者POST");
	}
	page = template.Must(template.New("page").Parse(expedit))
	data:= make(map[string]interface{})
	var exp expense.Expense
	log.Infof(ctx,"%s",r.FormValue("tags"))
	parseExpense(r,&exp)
	log.Infof(ctx,"%s",exp.Tags)
	if exp.Seq != "" {
		data["seq"]=exp.Seq
	}
	if exp.Remark != "" {
		data["remark"]=exp.Remark
	}
	if exp.Amount != 0 {
		data["amount"]=expense.Money(exp.Amount)
	}
	data["countin"]=exp.CountIn
	if exp.When != 0 {
		data["when"]=expense.JavaDateStr(exp.When)
	} else {
		data["when"]=expense.NowDateStr()
	}
	if exp.CreatedTime != 0 {
		data["createdtime"]=exp.CreatedTime
	}
	if exp.Tags != nil {
		data["tags"]=strings.Join(exp.Tags,",")
	}
	
	
	data["name"]=r.FormValue("account")
	ary := make([]string,0,30)
	ary = append(ary,expense.Money(exp.Amount))
	data["tagslist"]=ary
	page.Execute(w,data)
}

func parseExpense(r *http.Request,exp *expense.Expense) {
	_,ok:=r.Form["seq"]
	if ok {
		exp.Seq= r.FormValue("seq")
	}
	_,ok = r.Form["countin"]
	if ok {
		exp.CountIn=true
	} else {
		exp.CountIn=false
	}
	_,ok = r.Form["amount"]
	if ok {
		f,_:=strconv.ParseFloat(r.FormValue("amount"),64)
		exp.Amount = int(f * 100)
	}
	_,ok = r.Form["remark"]
	if ok {
		exp.Remark=r.FormValue("remark")
	}
	str := r.FormValue("when")
	if str!="" {
		exp.When=expense.JavaDateInt(str)
	}
	_,ok = r.Form["createdtime"]
	if ok {
		exp.CreatedTime,_=strconv.ParseInt(r.FormValue("createdtime"),10,64)
	}
	_,ok = r.Form["tags"]
	if ok {
		str = r.FormValue("tags")
		exp.Tags=strings.Split(str,",")
	}
	
	
}
func pfaList(ctx context.Context){
	w,r,_ := spark.ReadHttpContext(ctx)

	pg := template.Must(template.New("page").Parse(explist))
	data:= make(map[string]interface{})
	account :=r.FormValue("account")
	if account== ""{
		data["title"]="必须提供一个帐套名称"
	} else {
		data["title"]=account
		data["ok"]=true
	}
	
	pg.Execute(w,data)
}
const (
	expedit=`<!DOCTYPE HTML>
<html>
<head>
	<title>Expense Account {{.name}}</title>
</head>
<body>
<datalist id="tagslist">
{{range .tagslist}}
<option value="{{.}}" />
{{end}}
</datalist>
<form action="/edit" method="POST">
<table>
<tr><td>Amount</td><td><input type="number" step="0.01" required name="amount" placeholder="金额" value="{{if .amount}}{{.amount}}{{end}}" /></td></tr>
<tr><td>Count In</td><td><input type="checkbox" {{if .countin}}checked{{end}} value="true" name="countin" /></td></tr>
<tr><td>Tags</td><td><input name="tags" list="tagslist" placeholder="输入适合的标签" value="{{if .tags}}{{.tags}}{{end}}"/></td></tr>
<tr><td>Date</td><td><input name="when" type="date" value="{{if .when}}{{.when}}{{end}}" /></td></tr>
<tr>
  <td>Remark</td><td>
    <textarea name="remark" placeholder="这笔费用的描述">{{if .remark}}{{.remark}}{{end}}</textarea>
  </td>
</tr>

<tr><td colspan=2><input type="submit" /></td></tr>
</table>
{{if .seq}}
<input type="hidden" name="seq" value="{{.seq}}"/>
{{end}}
{{if .createdtime}}
<input type="hidden" name="createdtime" value="{{.createdtime}}" />
{{end}}
</form>
{{.data}}
</body>
</html>`
	explist=`<!DOCTYPE HTML>
<html>
<head><title>{{.title}}</title></head>
<body>
<h3>{{.title}}</h3>
{{if .ok}}<form action="/edit" method="POST"><input type="submit" value="添加" /></form>{{end}}
</body>
</html>
`
)
