package web

import (
	"io"
	"errors"
	"mime"
	"path/filepath"
	"google.golang.org/appengine/log"
	"cloud.google.com/go/storage"
	
	"golang.org/x/net/context"
	"github.com/liangx8/spark"
	"github.com/liangx8/spark/zip"
	"github.com/liangx8/spark/session"
	"github.com/liangx8/gcloud-helper/gcs"
)

func admin(ctx context.Context){
	bucketName := defaultConfig.BucketName
	w,r,err := spark.ReadHttpContext(ctx)
	if err != nil {
		log.Errorf(ctx,"%v",err)
		return
	}
	s,err :=session.GetSession(ctx)
	if err != nil {
		log.Errorf(ctx,"%v",err)
		return
	}
	if !s.Get("",nil) {
		pageLogin(w,r.URL.Path)
		return
	}
	data := make(map[string]interface{})
	action:=r.FormValue("action")
	switch action{
	case "upload":
		client, err:=storage.NewClient(ctx)
		if err != nil {
			log.Errorf(ctx,"%v",err)
			return
		}
		defer func(){
			if err := client.Close(); err != nil {
				log.Errorf(ctx,"%v",err)
			}
		}()
		bkt := &gcs.Bucket{B:client.Bucket(bucketName),C:ctx}

		f,_,err:=r.FormFile("filename")
		size,err := f.Seek(0,2)
		if err != nil {
			log.Errorf(ctx,"%v",err)
			return
		}
		err=bkt.Delete(defaultConfig.SourcesPrefix,func(name string){
			log.Infof(ctx,"delete object %s",name)
		})
		if err != nil {
			log.Errorf(ctx,"%v",err)
			return
		}
		list := make([]string,0,100)

		err =zip.Each(f,size,func(rc io.ReadCloser,name string)error{
			list = append(list,name)

			return bkt.NewObjectWriter(filepath.Join(defaultConfig.SourcesPrefix,name),guess,func(cw *storage.Writer)error{
				_,err:=io.Copy(cw,rc)
				return err

			})
			
		})
		if err != nil && err != Break{
			log.Errorf(ctx,"%v",err)
			return
		}
		data["list"]=list
	}
	t :=tmpl.Lookup("admin")
	
	err=t.Execute(w,data)
	if err != nil {
		log.Errorf(ctx,"%v",err)
	}
	
}
func guess(name string)string{
	return mime.TypeByExtension(filepath.Ext(name))
}
var Break = errors.New("break")
const adminPage =`
<html>
<head><title>Admin</title></head>
<body>

<form action="" method="POST" enctype="multipart/form-data">
Select a zip file<input type="file" name="filename" /><input type="submit" />
<input type="hidden" name="action" value="upload" />
</form>
<h3>Welcome</h3>
{{range $idx, $elem := .list}}{{$idx}} {{$elem}}<br />{{end}}
</body>
</html>
`


/*


Welcome

*/
