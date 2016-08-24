package web

import (
	"net/http"
	"html/template"
	"time"

	"golang.org/x/net/context"
	"github.com/liangx8/spark"
	"github.com/liangx8/spark/session"
)
func pageReset(ctx context.Context){

	data:=make(map[string]string)
	data["title"]="Reset"
	w,r,_ := spark.ReadHttpContext(ctx)
	resetId := r.FormValue("resetid")
	if resetId == defaultConfig.ResetId && time.Now().Before(defaultConfig.ResetTimeout) {
		data["message"]="reset passord"
		pwd := r.FormValue("pwd")
		defaultConfig.Passphase=pwdEncrypt(resetId[:4],pwd)
		defaultConfig.ResetId=""
		saveConfig(ctx,&defaultConfig)
	} else {
		data["message"]="Invail linked!, try ?resetid=....&pwd=your password"
	}
	tmpl.Lookup("loginresult").Execute(w,data)
}
func login(ctx context.Context){
	data:=make(map[string]string)
	s,_ := session.GetSession(ctx)
	w,r,err := spark.ReadHttpContext(ctx)
	if err != nil {panic(err)}
	pc := r.FormValue("passcode")
	var loginOk bool
	if pc == "" {
		data["status"]="failure"
		loginOk=false
	} else {
		if authorized(ctx,pc) {
			data["title"]="successful"
			data["message"]="Successful"
			loginOk=true
		} else {
			data["title"]="failure"
			data["message"]="Failure"
			loginOk=false
		}
	}
	s.Put("",&loginOk)
	tmpl.Lookup("loginresult").Execute(w,data)
}
func pageLogin(w http.ResponseWriter,src string){
	data:=make(map[string]string)
	data["src"]=src
	tmpl.Lookup("login").Execute(w,data)
}
func parseTemplate(){
	tmpl =template.Must(template.New("login").Parse(body))
	template.Must(tmpl.New("loginresult").Parse(login_result))
	template.Must(tmpl.New("admin").Parse(adminPage))
}
var tmpl *template.Template
const (
	body = `
<html>
<head><title>LOGIN</title>
<style type="text/css">
html, body {
font-family: "Roboto", sans-serif;
color: #333333;
background-color: #35ea78;
margin: 0px;
}
h1 {
color: #d04526
background-color: #eaeaea;
padding: 20px;
border-bottom: 1px dashed #2b3848;
}
pre {
margin: 20px;
padding: 20px;
border: 2px solid #2b3848;
background-color: #eaeaea;
}
</style>
</head><body>
<h1>Authorization</h1>
<pre style="font-weight: bold;">Please input the pass pharse for authorized to access page {{.src}}</pre>
<form action="login" method="POST">
pass code:<input name="passcode" /><input type="submit" /><input type="hidden" name="src" value="{{.src}}" />
</form>

</body>
</html>

`
	login_result =`
<html>
<head><title>{{.title}}</title>
<style type="text/css">
html, body {
font-family: "Roboto", sans-serif;
color: #333333;
background-color: #a8ea78;
margin: 0px;
}
h1 {
color: #d04526
background-color: #eaeaea;
padding: 20px;
border-bottom: 1px dashed #2b3848;
}
pre {
margin: 20px;
padding: 20px;
border: 2px solid #2b3848;
background-color: #eaeaea;
}
</style>
</head><body>
<h1>{{.title}}</h1>
<pre style="font-weight: bold;">{{.message}}</pre>
</body>
</html>
`
)
