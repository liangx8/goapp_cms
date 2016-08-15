package gcs
import (

	"golang.org/x/net/context"
	"google.golang.org/cloud/storage"
)


func All(ctx context.Context,
	bucket *storage.BucketHandle) func( func(*storage.ObjectHandle)error) error {
	var ea func(string,func(*storage.ObjectHandle)error) error
	ea = func(prefix string,cb func(*storage.ObjectHandle)error) error{
		query := &storage.Query{Prefix:prefix,Delimiter:"/"}
		var pf []string
//		for query != nil {
		objs,err := bucket.List(ctx,query)
		if err != nil{ return err}
		pf=objs.Prefixes
		for _,res := range objs.Results {
			err = cb(bucket.Object(res.Name))
			if err != nil { return nil }
		}
		// not sure what is objs.Next, studying it in future
//			query = objs.Next
//		}
		for _,p:=range pf {
			err := ea(p,cb)
			if err != nil { return nil }
		}
		return nil
	}
	return func(cb func(*storage.ObjectHandle)error) error{
		return ea("",cb)
	}
}
