package web

import (
	"fmt"
	"crypto/md5"
	"cloud.google.com/go/storage"
	"golang.org/x/net/context"
	"google.golang.org/appengine/log"
)

func authorized(ctx context.Context,pass string)bool{
	LoadConfig(ctx)
	log.Infof(ctx,"encoding: %s",defaultConfig.Passphase)
	return equal(pass,defaultConfig.Passphase)
	
}
func equal(plain,enc string)bool{
	l := len(enc)
	if l < 4 {
		return false 
	}
	salt := enc[:4]
	m :=md5.New()
	fmt.Fprint(m,plain)
	b:=m.Sum([]byte(salt))
	bs := fmt.Sprintf("%x",b)

	return bs==enc[4:]
}
func pwdEncrypt(salt,pwd string)string{
	m := md5.New()
	fmt.Fprint(m,pwd)
	b:=m.Sum([]byte(salt))
	return fmt.Sprintf("%s%x",salt,b)
}
func saveConfig(ctx context.Context,c *Config){
		bucketName := defaultConfig.BucketName
	client,err:=storage.NewClient(ctx)
	if err != nil {
		log.Errorf(ctx,"%v",err)
		return
	}
	defer client.Close()
	
	bucket := client.Bucket(bucketName)
	cfgObj:=bucket.Object("config.yaml")
	w:=cfgObj.NewWriter(ctx)
	cfgBuf,err:=c.Yaml()
	if err != nil {
		log.Errorf(ctx,"%v",err)
		return
	}
	_,err = w.Write(cfgBuf)
	if err != nil {
		log.Errorf(ctx,"%v",err)
		return
	}
	defer w.Close()
}
