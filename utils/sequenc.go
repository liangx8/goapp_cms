package utils
import (
	"errors"
)
type node struct{
	e interface{}
	next *node
}
type Iterator interface{
	Next() error
	Get() (interface{},error)
	Evict() (interface{},error)
}
type Collection interface{
	Add(interface{})
	Iterate() Iterator
}

type coll struct{
	head *node
}
type iterator struct{
	boc bool
	phead **node
	prev,current *node
}
func (itr *iterator)Next() error{
	itr.boc=false
	if itr.current == nil {
		return EOC
	}
	itr.prev=itr.current
	itr.current=itr.current.next
	if itr.current == nil {
		return EOC
	}
	return nil
}
func (itr *iterator)Get() (interface{},error){
	if itr.boc {
		return nil,BOC
	}
	if itr.current == nil {
		return nil,EOC
	}
	return itr.current.e,nil
}
func (itr *iterator)Evict() (interface{},error){
	if itr.boc {
		return nil,BOC
	}
	if itr.current == nil {
		return nil,EOC
	}
	n := itr.current
	itr.current=itr.current.next
	itr.prev.next=itr.current
	if *itr.phead == n {
		*itr.phead = itr.current
	}
	return n.e,nil
}
func (co *coll)Add(e interface{}){
	n := &node{e:e,next:co.head}
	co.head=n
}
func (co *coll)Iterate() Iterator{
	bof := &node{next:co.head}
	return &iterator{phead:&co.head,boc:true,prev:nil,current:bof}
}
func NewCollection() Collection{
	return &coll{head:nil}
}

var EOC = errors.New("End of Collection")
var BOC = errors.New("Begin of Collection")
