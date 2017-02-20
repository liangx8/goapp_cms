package web

import (
	"golang.org/x/net/context"
	"google.golang.org/appengine/log"

	"github.com/liangx8/spark"
	"github.com/liangx8/gcloud-helper/gcs"

	"expense"
	
	"fmt"
	"strings"
	"io"
	"io/ioutil"
	"unicode/utf8"
	"encoding/json"
)
/*
 * 第一个字节
 *   1 数据上传，下一个字节未帐套都名称的长度,后面的指定长度内容必须是帐套的名称(utf-8编码)
 *   2 数据下载，下一个字节
 *      3 请求帐套的列表
 *      4 请求帐套内容，下一个字节为帐套的名称长度,后面的指定长度内容必须是帐套的名称(utf-8编码)(未实现)
 */
func pfa(ctx context.Context){
	w,r,_ := spark.ReadHttpContext(ctx)
	log.Infof(ctx,"开始")
	if r.Method != "PUT" && r.Method != "POST" {
		fmt.Fprintf(w,returnError,"必须是PUT或者POST");
		return
	}
	buf := make([]byte,2)
	num,_ := r.Body.Read(buf)
	if num <2{
		fmt.Fprintf(w,returnError,"最少需要2个字节")
		return
	}
	fileName := r.URL.RequestURI()
	if !strings.HasSuffix(fileName,".json"){
		fmt.Fprintf(w,returnError,"目前只会用JSON")
		return
	}
	if buf[0] == data_outgoing {
		if buf[1]==account_request {
			bucket,err := gcs.NewBucket(ctx,projectId,bucketName)
			if err != nil {
				fmt.Fprintf(w,returnError,err)
				return
			}
			err = expense.AllAccount(bucket,jsonListAccount(w))
			if err != nil {
				fmt.Fprintf(w,returnError,err)
			}
			return
		}
		comp := buf[1]
		if comp==expense_request || comp==expense_request_new {
			num,_=r.Body.Read(buf[0:1])
			if num != 1 {
				fmt.Fprintf(w,returnError,"期望帐套名称的长度")
				return
			}
			length := int(buf[0])
			buf = make([]byte,length)
			num,_ := r.Body.Read(buf)
			if num != length{
				fmt.Fprintf(w,returnError,"帐套名称长度不符")
				return
			}
			account := string(buf)
			if !utf8.ValidString(account){
				fmt.Fprintf(w,returnError,"必须提供一个UTF-8编码的账户名")
				return
			}
			//			fmt.Fprintf(w,"请求帐套: %s",account)
			if comp==expense_request {
				exps:=expense.OldData(ctx)
				if exps == nil {
					fmt.Fprint(w,"{}")
				} else {
					buf, _ = json.Marshal(exps)
					w.Write(buf)
				}
			} else {
				cloud,err := expense.NewCloud(ctx,account)
				if err != nil {
					fmt.Fprintf(w,returnError,err)
					return
				}
				defer cloud.Close()
				var exps []expense.Expense
				err = cloud.Load(&exps)
				if err != nil {
					fmt.Fprintf(w,returnError,"必须提供一个UTF-8编码的账户名")
				}
				if exps == nil {
					fmt.Fprint(w,"{}")
				} else {
					buf, _ = json.Marshal(exps)
					w.Write(buf)
				}
			}
		}
	}
	if buf[0] == data_incoming {
		acLen := int(buf[1])
		buf= make([]byte,acLen)
		num,_ = r.Body.Read(buf)
		if num != acLen {
			fmt.Fprintf(w,returnError,"帐套名称长度不符")
			return
		}
		account := string(buf)
		if !utf8.ValidString(account){
			fmt.Fprintf(w,returnError,"必须提供一个UTF-8编码的账户名")
			return
		}
		// TODO: 读数据，然后保存
		buf,err := ioutil.ReadAll(r.Body)
		if err != nil {
			fmt.Fprintf(w,returnError,err)
			return
		}
		var exps []expense.Expense
		err = json.Unmarshal(buf,&exps)
		if err != nil {
			fmt.Fprintf(w,returnError,err)
			return
		}
		cloud,err := expense.NewCloud(ctx,account)
		if err != nil {
			fmt.Fprintf(w,returnError,err)
			return
		}
		defer cloud.Close()
		cloud.Merge(exps,func(add,update int, err error){
			if err == nil {
				fmt.Fprintf(w,"新增%d条记录，更新%d条记录\n",add,update)
			} else {
				fmt.Fprintf(w,returnError,err)
			}
		})
	}

}
func jsonListAccount(w io.Writer) func(string){
	first := true
	return func(ac string){
		if first {
			first=false
			fmt.Fprintf(w,`["%s"`,ac)
			return
		}
		if ac == "" {
			fmt.Fprintf(w,"]")
			return
		}
		fmt.Fprintf(w,`,"%s"`,ac)
	}
}

const (
	returnError = `["error","%s"]`
	data_incoming       byte = 1
	data_outgoing       byte = 2
	account_request     byte = 3
	expense_request     byte = 4
	expense_request_new byte = 5
)
