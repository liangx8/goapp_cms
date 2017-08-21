package tmpl

import (
	"html/template"
	"io/ioutil"
	"utils"

//	"cloud.google.com/go/storage"
	"golang.org/x/net/context"

	"github.com/liangx8/gcloud-helper/gcs"
)
func MustBuildTemplate(ctx context.Context,name string) *template.Template{
	return template.Must(BuildTemplate(ctx,name))
}
func BuildTemplate(ctx context.Context,name string) (*template.Template, error){
	bh,closer,err := gcs.MakeBucket(ctx,BUCKET_NAME)
	if err != nil { return nil,err}
	defer closer()
	fname := prefix + "/" + name + ".tmpl"
	oh := bh.Object(fname)
	objr,err := oh.NewReader(ctx)
	if err != nil {
		return nil,utils.WrappingError(err,"New Reader of object " +fname + " error in bucket 'pfa.rc-greed.com'")
	}
	defer objr.Close()
	buf,err := ioutil.ReadAll(objr)
	if err != nil {
		return nil,utils.WrappingError(err,"Read object " +fname + " error in bucket 'pfa.rc-greed.com'")
	}
	return template.New(name).Parse(string(buf))
}

const (
	prefix="tmpl"
	BUCKET_NAME="pfa.rc-greed.com"
)
