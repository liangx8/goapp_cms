package one

import(
	"io"
	"os"
	"net/http"
	"html/template"



	"appengine"

	"zpack"
	"wraperror"
)

var tmpl *template.Template
var header *template.Template
var tailer *template.Template
func init(){
	header = template.Must(template.New("header").Parse(HEAD))
	tailer = template.Must(template.New("tailer").Parse(TAIL))
	tmpl = template.Must(template.New("upload.form").Parse(UPLOAD_FORM))
	template.Must(tmpl.New("list.table").Parse(LIST_TABLE))
	template.Must(tmpl.New("error").Parse(ERROR_DIV))
}
func backEnd(w http.ResponseWriter,r *http.Request){
	c:=appengine.NewContext(r)
	var filelist []string
	var err error


	//	type zCallback func(io.Reader,os.FileInfo) error
	each :=	func(rd io.Reader,fi os.FileInfo,fullname string,ccc appengine.Context) error{
		if !fi.IsDir(){
			filelist = append(filelist,fullname)
			er:=store(ccc,rd,fullname,guessMimeType(fullname))
			if er!= nil {
				return wraperror.Printf(er,"the file '%s' in your upload is too large",fullname)
			}
		}
		return nil
	}


	if r.Method == "POST" {
		ty := r.FormValue("type");
		switch ty{
		case "zipupload":
			filelist =make([]string,0,10)
			f,_,er := r.FormFile("filename")
			if er != nil {
				c.Errorf("%v",er)
				err=er
			} else {
				err=deleteAll(c,func(cc appengine.Context)error{
					return zpack.ZipForEach(f,
						func(rd io.Reader,fi os.FileInfo,fullname string) error{
							return each(rd,fi,fullname,cc)
						})
				})
				if err != nil {
						c.Errorf("%v",err)
				}
			}
		case "gzipupload":
			filelist =make([]string,0,10)
			filelist =make([]string,0,10)
			f,_,er := r.FormFile("filename")
			if er != nil {
				c.Errorf("%v",er)
				err=er
			} else {
				err=deleteAll(c,func(cc appengine.Context)error{
					return zpack.TarForEach(f,
						func(rd io.Reader,fi os.FileInfo,fullname string) error{
							return each(rd,fi,fullname,cc)
						})
				})
				if err != nil {
					c.Errorf("%v",err)
				}

			}
		default:

		}
	}
	header.Execute(w,map[string]string{"title":"Upload"})
	t := tmpl.Lookup("upload.form")
	t.Execute(w,map[string]string{"url":"/showdoor","prompt":"zip file","type":"zipupload"});
	t.Execute(w,map[string]string{"url":"/showdoor","prompt":"gzip file","type":"gzipupload"});
	if err != nil {

		t=tmpl.Lookup("error")
		t.Execute(w,map[string]interface{}{"error":err})

	} else {
		if r.Method == "POST" {
			t=tmpl.Lookup("list.table")
			t.Execute(w,map[string]interface{}{"list":filelist})
		}
	}
	tailer.Execute(w,nil);
}

const (
	HEAD=`<!DOCTYPE HTML><html><meta http-equiv="content-type" content="text/html" /><head><title>{{.title}}</title></head><body>`
	TAIL="</body></html>"
	UPLOAD_FORM=`<form action="{{.url}}" method="POST" enctype="multipart/form-data"> {{.prompt}}:<input type="file" name="filename" /><input type="submit" /><input type="hidden" name="type" value="{{.type}}" />
</form>`
	LIST_TABLE="<ul>\n{{range $idx,$elem := .list}}  <li>{{$elem}}</li>\n{{end}}</ul>"
	ERROR_DIV=`<div style="background-color:red;color:white"><pre>{{.error}}</pre></div>`
)
