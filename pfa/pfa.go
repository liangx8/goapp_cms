package pfa

import (
	"fmt"
	"strings"
	"io/ioutil"
	"encoding/json"
	"io"
	"unicode/utf8"
	"golang.org/x/net/context"
	"google.golang.org/appengine/log"

	"github.com/liangx8/spark"
	"gopkg.in/yaml.v2"
)
/*
 * 第一个字节
 *   1 数据上传，下一个字节未帐套都名称的长度,后面的指定长度内容必须是帐套的名称(utf-8编码)
 *   2 数据下载，下一个字节
 *      3 请求帐套的列表
 *      4 请求帐套内容，下一个字节为帐套的名称长度,后面的指定长度内容必须是帐套的名称(utf-8编码)(未实现)
 */
func PFA(ctx context.Context){
	w,r,err:=spark.ReadHttpContext(ctx)
	if err != nil {
		log.Errorf(ctx,"%v",err)
		return
	}
	switch r.Method{
	case "PUT","POST":
	default:
		index(ctx)
		return
	}
	buf :=make([]byte,2)
	num,err :=r.Body.Read(buf)
	if num < 2{
		fmt.Fprintln(w,"世界真美好!")
		log.Errorf(ctx,"%v",err)
		return
	}
	// trim "/pfa/"
	fileName := r.URL.RequestURI()
	var unmarshal func([]byte,interface{})error
	var listAccountRender func(string)
	if strings.HasSuffix(fileName,".json"){
		unmarshal = json.Unmarshal
		listAccountRender=jsonListAccount(w)
	}
	if strings.HasSuffix(fileName,".yaml"){
		unmarshal = yaml.Unmarshal
		listAccountRender=yamlListAccount(w)
	}
	if buf[0] == data_outgoing {
		if buf[1] == 3 {
			if listAccountRender == nil {
				fmt.Fprint(w,"必须告诉我要用什么格式?")
				log.Errorf(ctx,"必须告诉我要用什么格式?")
				return
			}
			bucket,err := myBucket(ctx)
			if err != nil {
				fmt.Fprint(w,err)
				return
			}
			defer bucket.Close()
			err = allAccount(bucket,listAccountRender)
			if err != nil {
				log.Errorf(ctx,"%v",err)
			}
			return
		}
		fmt.Fprint(w,"你想知道什么?")
		return
	}
	if buf[0] != data_incoming {
		fmt.Fprint(w,"我不知道你想要什么?")
		return
	}
	length := int(buf[1])
	if length==0{
		fmt.Fprint(w,"必须提供帐号名称")
		return
	}
	buf = make([]byte,length)
	_,err = r.Body.Read(buf)
	if err != nil {
		fmt.Fprintln(w,err)
		log.Errorf(ctx,"%v",err)
		return
	}
	account := string(buf)
	if !utf8.ValidString(account){
		fmt.Fprintln(w,"必须提供一个UTF-8编码的账户名")
		return
	}
	// 上传数据
	
	buf,err = ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintln(w,err)
		log.Errorf(ctx,"%v",err)
		return
	}

	exps := make([]Expense,0)
	if unmarshal == nil {
		fmt.Fprintf(w,"你要做什么?")
		return
	}
	err = unmarshal(buf,&exps)
	if err != nil {
		log.Errorf(ctx,"%v",err)
		fmt.Fprintf(w,result_template,false,err,"解析传入数据错误")
		return
	}
	dao,err := NewYamlDao(ctx,account)
	if err != nil {
		log.Errorf(ctx,"%v",err)
		fmt.Fprintf(w,result_template,false,err,"读取帐套错误")
		return
	}
	defer dao.Close()
	
	err=dao.Merge(exps,func(c1,c2 int){
		fmt.Fprintf(w,"更新了%d个项目，添加了 %d 个项目\n",c1,c2)
		fmt.Fprintf(w,"上传账套:%s\n",account)
	})
	if err != nil {
		log.Errorf(ctx,"%v",err)
		fmt.Fprintf(w,result_template,false,err,"保存帐套错误")
		return
	}
}
func yamlListAccount(w io.Writer) func(string){
	return func(a string){
		if a == "" { return }
		fmt.Fprint(w,"- ")
		fmt.Fprintln(w,a)
	}
}
func jsonListAccount(w io.Writer) func(string){
	first := true
	return func(ac string){
		if first {
			first=false
			fmt.Fprintf(w,"[\"%s\"",ac)
			return
		}
		if ac=="" {
			fmt.Fprintln(w,"]")
			return
		}
		fmt.Fprintf(w,",\"%s\"",ac)
	}
}
const (
	data_incoming byte = 1
	data_outgoing byte = 2
	result_template =`
ok: %v
error: %v
message: %s
`
)

