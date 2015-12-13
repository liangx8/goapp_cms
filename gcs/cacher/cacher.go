package cacher


import (
	"net/http"
	"errors"

	"golang.org/x/net/context"

	"google.golang.org/appengine"
	"google.golang.org/appengine/memcache"

	"utils"
)

type CacheContext struct {
	ctx context.Context

}
type CacheFile struct{
	ContentType string
	Content []byte
}

func NewCacheContext(r *http.Request) *CacheContext{
	return &CacheContext{ctx:appengine.NewContext(r)}
}

func (cc *CacheContext)Load(name string,defaultFile func(*CacheFile) error)(*CacheFile,error){
	var file CacheFile
	key := utils.ShaStr(name)
	var codec memcache.Codec
	codec.Marshal=marshal
	codec.Unmarshal=unmarshal
	_,err := codec.Get(cc.ctx,key,&file)
	if err == memcache.ErrCacheMiss {
		err = defaultFile(&file)
		if err != nil { return nil,err }
		it := &memcache.Item{Key:key,Object:&file}
		err = codec.Set(cc.ctx,it)
	}
	if err != nil {
		return nil,err
	}
	return &file,nil
}
func (cc *CacheContext)Flush()error{
	return memcache.Flush(cc.ctx)
}
func marshal(obj interface{})([]byte,error){
	if obj == nil {
		return nil,errors.New("marshal(): nil point exception")
	}
	file,ok := obj.(*CacheFile)
	if !ok {
		return nil,errors.New("cacher.marshal(interface{}) must be privded a *CacheFile object")
	}
	typesize := len(file.ContentType)
	size := typesize+2+len(file.Content)
	buf := make([]byte,size)
	buf[0]=byte(typesize % 256)
	buf[1]=byte(typesize / 256)
	i := 2
	for _,c := range []byte(file.ContentType){
		buf[i]=c
		i++
	}
	for _,c := range file.Content{
		buf[i]=c
		i++
	}
	return buf,nil
}
func unmarshal(buf []byte,obj interface{})error{
	if obj == nil || buf == nil {
		return errors.New("unmarshal(): nil point exception")
	}
	file,ok := obj.(*CacheFile)
	if !ok {
		return errors.New("cacher.unmarshal must be privded a *CacheFile object")
	}
	typesize := int(buf[0]) + int(buf[1]) * 256
	file.ContentType=string(buf[2:typesize+2])
	file.Content=buf[typesize+2:]
	return nil
}
