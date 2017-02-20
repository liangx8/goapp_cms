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
func pfaDelete(ctx context.Context){
	w,r,_ := spark.ReadHttpContext(ctx)
	account := r.FormValue("account")
	pg := template.Must(template.New("page").Parse(expdelete))
	data:= make(map[string]interface{})
	if account == "" {
		data["success"]=false
		data["message"]="必须提供帐套名称"
	} else {
		seq := r.FormValue("seq")
		cloud,err := expense.NewCloud(ctx,account)
		if err != nil {
			data["success"]=false
			data["message"]=err
			goto renderPage
		}
		err = cloud.Delete(seq)
		if err != nil {
			data["success"]=false
			data["message"]=err
			goto renderPage
		}
			data["success"]=true
	}
renderPage:
	pg.Execute(w,data)
	
}
func pfaEdit(ctx context.Context){
	w,r,_ := spark.ReadHttpContext(ctx)
	if r.Method != "POST" {
		panic("必须是PUT或者POST");
	}
	pg := template.Must(template.New("page").Parse(expedit))
	data:= make(map[string]interface{})
	r.ParseForm()
	_,ok :=r.Form["account"]
	if ! ok {
		data["error"]="必须提供帐套"
		data["name"]="错误"
		pg.Execute(w,data)
		return
	}
	accountName :=r.FormValue("account")
	data["account"]=accountName
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
	_,ok =r.Form["save"]
	if ok {
		cloud,err := expense.NewCloud(ctx,accountName)
		if err != nil {
			data["error"]=err
			data["title"]="错误"
			pg.Execute(w,data)
			return
		}
		defer cloud.Close()
		expense.CompleteExpense(&exp)
		err = cloud.InsertOrUpdate(exp)
		if err != nil {
		data["error"]=err
			data["title"]="错误"
			pg.Execute(w,data)
			return
		}
	}
	ary := make([]string,0,30)
	ary = append(ary,expense.Money(exp.Amount))
	data["tagslist"]=ary
	pg.Execute(w,data)
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
	r.ParseForm()
	_,ok := r.Form["account"]
	if ok {
		account :=r.FormValue("account")
		data["account"]=account
		data["ok"]=true
		cloud,err := expense.NewCloud(ctx,account)
		if err != nil {
			data["error"]=err
			data["title"]="错误"
			pg.Execute(w,data)
			return
		}
		defer cloud.Close()
		var exps []expense.Expense
		err = cloud.Load(&exps)
		if err != nil {
			data["error"]=err
			data["title"]="错误"
			pg.Execute(w,data)
			return
		}
		data["data"]=exps
		data["showdate"]=expense.JavaDateStr
		data["showts"]=expense.JavaTimestampStr
		data["showmoney"]=expense.Money
	} else {
		data["account"]="必须提供一个帐套名称"
	}
	
	pg.Execute(w,data)
}
const (
	expedit=`<!DOCTYPE HTML>
<html>
<head>
	<title>Expense Account {{.title}}</title>
</head>
<body>
{{if .error}}
<h3>{{.error}}</h3>
{{else}}
{{if .inputerror}}
{{.inputerror}}
{{end}}
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
<input type="hidden" name="account" value="{{.account}}"/>
{{if .createdtime}}
<input type="hidden" name="createdtime" value="{{.createdtime}}" />
{{end}}
<input type="hidden" name="save" value="save" />
</form>
{{end}}
</body>
</html>`
	explist=`{{$dateStr := .showdate}}{{$moneyStr := .showmoney}}{{$ac := .account}}{{$tsStr := .showts}}
<!DOCTYPE HTML>
<html>
<head><title>{{.title}}</title></head>
<body>
{{if .ok}}

<h3>列出{{.account}}的内容</h3>
<form action="/edit" method="POST">
<input type="submit" value="添加" />
<input type="hidden" name="account" value="{{.account}}" />
</form>
{{else}}
<h3>{{.account}}</h3>
{{end}}
<table>
{{range .data}}
<tr>
<td>{{call $dateStr .When}}</td>
<td>{{call $moneyStr .Amount}}</td>
<td>{{.Tags}}</td><td>{{.Remark}}</td>
<td>{{call $tsStr .CreatedTime}}</td>
<td><a href="/delete?seq={{.Seq}}&account={{$ac}}">删除</a></td>
</tr>
{{end}}
</table>
</body>
</html>
`
	expdelete=`{"success":{{.success}}{{if .message}},"message":"{{.message}}"{{end}} }`
)
