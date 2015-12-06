package one
import (
	"io"
	"io/ioutil"
	"errors"

	"appengine"
	"appengine/datastore"
)
type File struct{
	Name,MimeType string
	Content []byte
}
// for use transaction, a root key must be define
func getRootKey(c appengine.Context) *datastore.Key{
	return datastore.NewKey(c,"File","root",0,nil)
}
func deleteAll(c appengine.Context,between func(appengine.Context) error) error{
	keys,err := datastore.NewQuery("File").KeysOnly().GetAll(c,nil)
	if err != nil { return err}
	err = datastore.RunInTransaction(c,between,nil)
	if err != nil { return err}
	err=datastore.DeleteMulti(c,keys)
	if err != nil {	return err	}
	return nil
}
func store(c appengine.Context,r io.Reader, name,mimeType string) error{
	var file File


	file.Name=name
	file.MimeType=mimeType
	buf,err := ioutil.ReadAll(r)
	if err != nil {
		return err
	}
	file.Content=buf
	k := datastore.NewIncompleteKey(c,"File",getRootKey(c))
	_,err = datastore.Put(c,k,&file)
	if err != nil { return err }

	return nil
}

func getByName(c appengine.Context,name string) (*File,error){
	var f File
	q:=datastore.NewQuery("File").Filter("Name =",name)
	cur:=q.Run(c)

	_,err :=cur.Next(&f)
	if nil==err {
		return &f,nil
	}
	if err == datastore.Done{
		return nil,ErrNotFound
	}
	return nil,err
}

var ErrNotFound=errors.New("Resouce not found");
