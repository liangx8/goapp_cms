package one

import(
	"net/http"
//	"html/template"



	//	"golang.org/x/net/context"

	"google.golang.org/appengine"
	"google.golang.org/appengine/log"

//	"zpack"
	//	"wraperror"
	"tmpls"
)

func backEnd(w http.ResponseWriter,r *http.Request){
	c:=appengine.NewContext(r)
	var err error
	var filelist []string
	if r.Method == "POST" {

		ty := r.FormValue("type");
		switch ty{
		case "zipupload":

			f,_,er := r.FormFile("filename")
			if er != nil {
				log.Errorf(c,"%v",er)
				err=er
			} else {
				filelist,err=saveZip(c,f,true)
				if err != nil {
					log.Errorf(c,"%v",er)

				}
			}
		case "gzipupload":
			f,_,er := r.FormFile("filename")
			if er != nil {
				log.Errorf(c,"%v",er)

				err=er
			} else {
				filelist,err=saveZip(c,f,false)

				if err != nil {
					log.Errorf(c,"%v",er)


				}
			}
		default:

		}
	}
	tmpls.Header.Execute(w,map[string]string{"title":"Upload"})
	t := tmpls.Tmpl.Lookup("upload.form")
	t.Execute(w,map[string]string{"url":"/showdoor","prompt":"zip file","type":"zipupload"});
	t.Execute(w,map[string]string{"url":"/showdoor","prompt":"gzip file","type":"gzipupload"});
	if err != nil {

		t=tmpls.Tmpl.Lookup("error")
		t.Execute(w,map[string]interface{}{"error":err})

	} else {
		if r.Method == "POST" {
			t=tmpls.Tmpl.Lookup("list.table")
			t.Execute(w,map[string]interface{}{"list":filelist})
		}
	}
	tmpls.Tailer.Execute(w,nil);
}
