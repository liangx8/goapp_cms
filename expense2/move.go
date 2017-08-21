package expense2

import (
	"golang.org/x/net/context"

	"github.com/liangx8/spark"
	"github.com/liangx8/gcloud-helper/gcs"

	"tmpl"
)
func Move(ctx context.Context){
	

	// move 指文件 gs://pfa.rc-greed.com/tmpl/move.tmpl 
	page := tmpl.MustBuildTemplate(ctx,"account")
	bkt,closer,err:=gcs.MakeBucket(ctx,BUCKET_NAME)
	if err != nil {
		panic(err)
	}
	defer closer()
	list := make([]string,0,10)
	gcs.Objects(ctx,bkt,gcs.StringTranslate(func(name string)error{
		list=append(list,name[len(PREFIX):])
		return nil
	}),PREFIX)
	err=spark.NewPage(page).Render(ctx,list)
	if err != nil {
		panic(err)
	}
}
const(
	BUCKET_NAME="pfa.rc-greed.com"
	PREFIX="expense/"
)
