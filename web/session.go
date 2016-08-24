package web

import (

	"golang.org/x/net/context"
	"google.golang.org/appengine/memcache"
	"google.golang.org/appengine/log"

	"github.com/liangx8/spark/session"
)
type (
	memSession struct{
		ctx context.Context
		id string
		login bool
	}
)

// SessionMaker implements
func GaeSessionMaker(ctx context.Context,id string) session.Session{
	var s *memSession
	var ok bool
	if id== "" {
		ok = false
	} else {
		item,err := memcache.Get(ctx,id)
		if err != nil {
			log.Errorf(ctx,"%v",err)
			ok = false
		} else {
			ok = true
			if item.Value[0]==0 {
				s= &memSession{ctx:ctx,id:id,login: false}
			} else {
				s= &memSession{ctx:ctx,id:id,login: true}
			}
		}

	}
	if !ok {
		id = session.UniqueId()
		item := &memcache.Item{
			Key: id,
			Value:[]byte{0},
		}
		if err := memcache.Set(ctx,item); err != nil {
			log.Errorf(ctx,"%v",err)
			return nil
		}
		s=&memSession{ctx:ctx,id:id,login:false}
	}
	return s
}
// Session implement
func (se *memSession)Get(key string, ptr interface{})bool{
	return se.login
}
func (se *memSession)Put(key string, ptr interface{}){
	v,ok := ptr.(*bool)
	if !ok {
		log.Errorf(se.ctx,"Error:Only accept *bool for Put")
		return
	}
	b :=make([]byte,1)
	se.login=*v
	if *v {
		b[0]=1
	} else {
		b[0]=0
	}
	
	item := &memcache.Item{
		Key: se.id,
		Value:b,
	}
	if err := memcache.Set(se.ctx,item); err != nil {
		log.Errorf(se.ctx,"%v",err)
	}
}
func (se *memSession)Id()string{
	return se.id
}

