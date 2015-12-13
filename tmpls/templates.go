package tmpls
import (
	"html/template"
)
var Tmpl *template.Template
var Header *template.Template
var Tailer *template.Template

func init(){
	Header = template.Must(template.New("header").Parse(HEAD))
	Tailer = template.Must(template.New("tailer").Parse(TAIL))
	Tmpl = template.Must(template.New("upload.form").Parse(UPLOAD_FORM))
	template.Must(Tmpl.New("list.table").Parse(LIST_TABLE))
	template.Must(Tmpl.New("error").Parse(ERROR_DIV))
	template.Must(Tmpl.New("backend").Parse(BACKEND))
}

const (
	HEAD=`<!DOCTYPE HTML><html><meta http-equiv="content-type" content="text/html" /><head><title>{{.title}}</title></head><body>`
	TAIL="</body></html>"
	UPLOAD_FORM=`<form action="{{.url}}" method="POST" enctype="multipart/form-data"> {{.prompt}}:<input type="file" name="filename" /><input type="submit" /><input type="hidden" name="type" value="{{.type}}" />
</form>`
	LIST_TABLE="<ul>\n{{range $idx,$elem := .list}}  <li>{{$elem}}</li>\n{{end}}</ul>"
	ERROR_DIV=`<div style="background-color:red;color:white"><pre>{{.error}}</pre></div>`

	BACKEND = HEAD +
		`{{$url := .url}}<form action="{{.url}}" method="POST" enctype="multipart/form-data"> {{.prompt}}:<input type="file" name="filename" /><input type="submit" /><input type="hidden" /></form>
{{if .errcnt}}<div style="background-color:red;color:white">{{range $_,$err := .errors}}{{$err}}<br />{{end}}</div>{{end}}
{{range $idx0,$elem0 := .list}}
	<table><tr><th>name</th><th>apply</th></tr>
	{{range $idx1,$elem1 := $elem0.Prefixes}} <tr><td>{{$elem1}}</td><td>open</td></tr> {{end}}
	{{range $idx1,$elem1 := $elem0.Attrs}} <tr><td>{{$elem1.Name}}</td><td><a href="{{$url}}&apply={{$elem1.Name}}">apply</a></td></tr> {{end}}
{{end}}</table>` +
	TAIL
)
