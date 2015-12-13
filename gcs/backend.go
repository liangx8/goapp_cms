package gcs
import (
	"net/http"
	"fmt"
	"io"
	"os"

	"google.golang.org/cloud/storage"
	"google.golang.org/appengine/log"

	"tmpls"
	"wraperror"
	"zpack"
	"utils"
	"gcs/cacher"
)
type data struct{
	Prefixes []string
	Attrs []*storage.ObjectAttrs
}
type loopImp struct{
	da []*data

}
func (loopImp)Start(){}
func (l *loopImp)Each(pfx []string,attrs []*storage.ObjectAttrs)error{
	l.da = append(l.da,&data{Prefixes:pfx,Attrs:attrs})
	return nil
}
func (l *loopImp)End(){
}
func saveFile(bkt *GcsBucket,r *http.Request)error{
	f,fh,err := r.FormFile("filename")
	if err != nil {
		return err
	}
/*
保存在FileHeader中的信息
fh.Filename:banner.svg
fh.Header["Content-Disposition"]=[`form-data; name="filename"; filename="banner.svg"`]
fh.Header["Content-Type"]=["image/svg+xml"]
*/
	/*
	log.Infof(bkt.ctx,"filename:%s\n",fh.Filename)
	for k,v :=range fh.Header {
		log.Infof(bkt.ctx, "%s=>%v\n",k,v)
	}*/
	err=bkt.CreateObject("sources/"+fh.Filename,fh.Header["Content-Type"][0],f)
	if err != nil {
		return err
	}
	return nil
}
func delFiles(bkt *GcsBucket)error {
	if err:=bkt.AllObject("web/",
		func(prefix string)Action{
			if prefix == "web/static/" {return NoAction}
			return Enter
		},
		func(attrs *storage.ObjectAttrs)Action{
			return Delete
		}); err != nil {return err}
	return nil
}
func webFiles(bkt *GcsBucket,name string)error{
	//type Zcallback func(io.Reader,os.FileInfo,string) error
	r,attrs,err:=bkt.ReadObject(name)
	if err != nil { return err }
	defer r.Close()
	zra,err := zpack.NewZipReaderAdaptor(r,attrs.Size)
	if err != nil { return err }
	err = delFiles(bkt)
	if err != nil { return err }
	err=zpack.ZipForEach(zra,func(src io.Reader,fi os.FileInfo,fullname string) error{
		if fi.IsDir() {return nil}
		log.Infof(bkt.ctx,"%s\n",fullname)
		err := bkt.CreateObject("web/"+fullname,utils.GuessMimeType(fullname),src)
		if err != nil { return nil }
		return nil
	})
	if err != nil { return err }
	return nil
}
func Backend(w http.ResponseWriter,r *http.Request){
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	var merr wraperror.MultiError
	bucket,err := GetBucket(r,BUCKET_NAME)
	if err != nil {
		fmt.Fprintln(w,err)
		return
	}
	prefix := r.FormValue("prefix")
	defer bucket.Close()
	loop := &loopImp{da:make([]*data,0,10)}

	err=bucket.ListDir("sources/"+prefix,loop)
	if err == NoFound {
		err = wraperror.Printf(err,"dir '%s' is not exists, use default",prefix)
		merr = append(merr,err)
		prefix=""
		err=bucket.ListDir("sources/",loop)
	}
	if err != nil {
		merr = append(merr,err)
		fmt.Fprintln(w,err)
		goto err_end
	}
	if r.Method == "POST" {
		err=saveFile(bucket,r)
		if err != nil {
			merr = append(merr,err)
		}
	} else {
		// assums apply parameter is avaliable
		apply := r.FormValue("apply")
		if apply != "" {
			err = webFiles(bucket,apply)
			if err != nil {
				merr = append(merr,err)
			} else {
				err = cacher.NewCacheContext(r).Flush()
				merr = append(merr,err)
			}
		}
	}

err_end:
	t :=tmpls.Tmpl.Lookup("backend")
	t.Execute(w,map[string]interface{}{
		"errors":merr,
		"errcnt":len(merr),
		"title":"backend manager",
		"list":loop.da,
		"url":"gcsb?prefix="+prefix,
		"prompt":"choose file",
	})
}

