package pfa

import (
	"fmt"
	"strings"
	"io/ioutil"
	"encoding/json"
	"golang.org/x/net/context"
	"google.golang.org/appengine/log"

	"github.com/liangx8/spark"
	"gopkg.in/yaml.v2"
	"unicode/utf8"
)
func PFA(ctx context.Context){
	w,r,err:=spark.ReadHttpContext(ctx)
	if err != nil {
		log.Errorf(ctx,"%v",err)
		return
	}
	switch r.Method{
	case "PUT","POST":
	default:
		fmt.Fprintf(w,"Welcome!")
		return
	}
	buf :=make([]byte,2)
	_,err =r.Body.Read(buf)
	if err != nil {
		fmt.Fprintln(w,"世界真美好!")
		log.Errorf(ctx,"%v",err)
		return
	}
	if buf[0] == data_outgoing {
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
	// trim "/pfa/"
	fileName := r.URL.RequestURI()
	
	buf,err = ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintln(w,err)
		log.Errorf(ctx,"%v",err)
		return
	}

	exps := make([]Expense,0)
	var unmarshal func([]byte,interface{})error
	if strings.HasSuffix(fileName,".json"){
		fmt.Fprintln(w,"JSON")
		unmarshal = json.Unmarshal
	}
	if strings.HasSuffix(fileName,".yaml"){
		unmarshal = yaml.Unmarshal
	}
	if unmarshal == nil {
		fmt.Fprintf(w,"你要做什么?")
		return
	}
	err = unmarshal(buf,&exps)
	if err != nil {
		log.Errorf(ctx,"%v",err)
		return
	}
	dao,err := NewDao(ctx,account)
	if err != nil {
		log.Errorf(ctx,"%v",err)
		return
	}
	defer dao.Close()
	err=dao.Save(exps)
	if err != nil {
		log.Errorf(ctx,"%v",err)
		return
	}
	fmt.Fprintf(w,"上传了%d个项目\n",len(exps))
	fmt.Fprintf(w,"上传账套:%s\n",account)
}
const (
	data_incoming byte = 1
	data_outgoing byte = 2
)
