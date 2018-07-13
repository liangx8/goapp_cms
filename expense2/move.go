package expense2

import (
	"golang.org/x/net/context"

	"github.com/liangx8/spark"
	"github.com/liangx8/gcloud-helper/gcs"

	"tmpl"
)
func Move(ctx context.Context){
	

	// move 指文件 gs://pfa.rc-greed.com/tmpl/move.tmpl
	pp:=spark.NewParamParser(ctx)
	account := NO_STRING
	if err := pp.Populate("account",&account);err != nil {
		panic(err)
	}
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
	data := make(map[string]interface{})
	data["list"]=list
	data["account"]=account
	err=spark.NewTemplateView(page).Render(ctx,data)
	if err != nil {
		panic(err)
	}
}
const(
	BUCKET_NAME="pfa.rc-greed.com"
	PREFIX="expense/"
	NO_STRING="__no_string__"
)
