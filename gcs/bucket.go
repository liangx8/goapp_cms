package gcs

import (
	"net/http"
	"io"
	"errors"

	"golang.org/x/net/context"
	"google.golang.org/appengine"
//	"google.golang.org/appengine/file"
	"google.golang.org/appengine/log"
	"google.golang.org/cloud/storage"
//	"google.golang.org/appengine/user"

	"wraperror"
)
type Loop interface{
	Start()
	Each([]string,[]*storage.ObjectAttrs) error
	End()
}
type GcsBucket struct{
	bucket *storage.BucketHandle
	client *storage.Client
	ctx context.Context
}
func (bucket *GcsBucket)Close(){
	bucket.client.Close()
}
func (bucket *GcsBucket)Context(cb func(context.Context)){cb(bucket.ctx)}
func (bucket *GcsBucket)Bucket(cb func(*storage.BucketHandle)){cb(bucket.bucket)}
func GetBucket(r *http.Request,name string) (*GcsBucket,error){
	ctx := appengine.NewContext(r)
	client,err := storage.NewClient(ctx)
	if err != nil {
		return nil,err
	}
	return &GcsBucket{bucket:client.Bucket(name),ctx:ctx,client:client},nil
}
func (bkt *GcsBucket)ListDir(name string,loop Loop)error{
	log.Infof(bkt.ctx,"list prefix:%s",name)
	query := &storage.Query{Prefix:name,Delimiter:"/"}
	emptyResult := true

	loop.Start()
	defer loop.End()
	for query != nil {
		objs, err:=bkt.bucket.List(bkt.ctx,query)
		if err != nil {
			return wraperror.Printf(err,"ListDir error for %s",query.Prefix)
		}
		if (len(objs.Prefixes)>0 || len(objs.Results) > 0 ){
			err=loop.Each(objs.Prefixes,objs.Results)
			if err != nil {
				return err
			}
			emptyResult = false
		}

		query = objs.Next
	}
	if emptyResult {
		return NoFound
	}
	return nil
}
func (bkt *GcsBucket)CreateObject(name string,contentType string,r io.Reader)error{
	w :=bkt.bucket.Object(name).NewWriter(bkt.ctx)
	w.ContentType = contentType
	if _,err:=io.Copy(w,r); err != nil {
		w.CloseWithError(err)
		return err
	}
	w.Close()
	return nil
}
func (bkt *GcsBucket)ReadObject(name string)(r io.ReadCloser,attrs *storage.ObjectAttrs,err error){
	obj := bkt.bucket.Object(name)
	if r,err = obj.NewReader(bkt.ctx); err != nil {
		return
	}
	if attrs,err = obj.Attrs(bkt.ctx); err != nil {
		return
	}
	return
}
func (bkt *GcsBucket)DeleteObject(name string)error {
	obj := bkt.bucket.Object(name)
	if err := obj.Delete(bkt.ctx) ; err != nil {
		return err
	}
	return nil
}
func (bkt *GcsBucket)AllObject(prefix string,
	prefixAction func(string)Action,
	objAction func(*storage.ObjectAttrs)Action) error{
	query := &storage.Query{Prefix:prefix,Delimiter:"/"}

	for query != nil {
		objs, err:=bkt.bucket.List(bkt.ctx,query)
		if err != nil {
			return wraperror.Printf(err,"prefix error for %s",query.Prefix)
		}
		for _,o := range objs.Results{
			act := objAction(o)
			if act == Delete {
				// delete it
				err:=bkt.DeleteObject(o.Name)
				if err != nil {
					return err
				}
			}
			if act == Break { return nil}
		}
		for _,p := range objs.Prefixes {
			act :=prefixAction(p)

			if act == NoAction {continue}
			if act == Break { return nil}
			err = bkt.AllObject(p,prefixAction,objAction)
			if err != nil {
				return err
			}

		}

		query = objs.Next
	}

	return nil
}

var NoFound=errors.New("No Found")
type Action int
const (
	NoAction Action = iota
	Enter
	Delete
	Break
)
