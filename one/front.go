package one

import(
	"net/http"
	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
//	"golang.org/x/net/context"
	"wraperror"
)

func front(w http.ResponseWriter,r *http.Request){
	c := appengine.NewContext(r)
	url := r.URL.Path
	var f *File
	var err error
	if url == "/" {
		f,err=getByName(c,"index.html")
	} else {
		f,err = getByName(c,url[1:])
	}
	if err != nil{
		log.Errorf(c,"%v",wraperror.Printf(err,"path %s is not exist anymore",url))
		http.NotFound(w,r)
		return
	}
	w.Header().Set("Content-Type",f.MimeType)
	w.Write(f.Content)
}
