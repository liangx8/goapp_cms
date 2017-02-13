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
		if buf[1]==expense_request{
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
			expenses:=expense.OldData(ctx)
			if expenses != nil {
				buf,_=json.Marshal(expenses)
				w.Write(buf)
			}
		}
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
	data_incoming   byte = 1
	data_outgoing   byte = 2
	account_request byte = 3
	expense_request byte = 4
)
