package gcs



import(
	"net/http"
	"fmt"
	"io/ioutil"
//	"os"

//	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
	"google.golang.org/cloud/storage"
//	"google.golang.org/appengine/file"
	"golang.org/x/net/context"

	//	"wraperror"
	//	"zpack"
	"gcs/cacher"
)


func Front(w http.ResponseWriter,r *http.Request){
	var err error
	var file *cacher.CacheFile
	bucket,err := GetBucket(r,BUCKET_NAME)
	if err != nil{
		bucket.Context(func(ctx context.Context){
			log.Errorf(ctx,"%v",err)
		})

		w.Header().Set("Content-Type","text/plain")
		fmt.Fprintf(w,"%v",err)
	}
	defer bucket.Close()
	url := r.URL.Path
	if url == "/" {
		url = "/index.html"
	}
	cache:=cacher.NewCacheContext(r)
	file,err = cache.Load(url,func(f *cacher.CacheFile)error{
		bucket.Context(func(ctx context.Context){
			log.Infof(ctx,"touch cache %s",url)
		})

		rd,attrs,err:=bucket.ReadObject("web"+url)
		if err != nil { return err }
		f.ContentType=attrs.ContentType
		f.Content,err=ioutil.ReadAll(rd)
		if err != nil { return err }
		return nil
	})

	if err == storage.ErrObjectNotExist {
		http.NotFound(w,r)
		return
	}
	if err != nil{
		bucket.Context(func(ctx context.Context){
			log.Errorf(ctx,"%v",err)
		})

		w.Header().Set("Content-Type","text/plain")
		fmt.Fprintf(w,"%v",err)
	}
	w.Header().Set("Content-Type",file.ContentType)
	_,err = w.Write(file.Content)
	if err != nil{
		bucket.Context(func(ctx context.Context){
			log.Infof(ctx,"touch cache %s",url)
		})
		fmt.Fprintf(w,"%v",err)
	}

}


const (
	BUCKET_NAME="www.rc-greed.com"
)
