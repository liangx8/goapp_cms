package one
import (
	"os"
	"io"
	"io/ioutil"
	"mime/multipart"

	"golang.org/x/net/context"

	"google.golang.org/appengine/log"
	"google.golang.org/appengine/datastore"

	"zpack"
	"utils"
	"wraperror"
)
// for reduc the cost on appengin, we don't use index for the kind, only access them by sha1 of name
type File struct{
	Name  string `datastore:",noindex"`
	MimeType string `datastore:",noindex"`
	Content []byte
}
// for use transaction, a root key must be define
func getRootKey(c context.Context) *datastore.Key{
	return datastore.NewKey(c,FILE_KIND,"root",0,nil)
}
func store(c context.Context,r io.Reader, name,mimeType string) (*datastore.Key,error){
	var file File
	var err error
	file.Name=name
	file.MimeType=mimeType
	file.Content,err = ioutil.ReadAll(r)
	if err != nil { return nil,err }
	k := datastore.NewKey(c,FILE_KIND,ShaStr(name),0,getRootKey(c))
	k,err = datastore.Put(c,k,&file)
	if err != nil { return nil,err }

	return k,nil
}
func builderEach(c context.Context,col utils.Collection,fl *[]string) zpack.Zcallback{
	return func(r io.Reader,fi os.FileInfo,fullname string) error{
		k,err:=store(c,r,fullname,guessMimeType(fullname))
		if err != nil { return wraperror.Printf(err,"the error was probobly cause by that file '%s' in package is too large", fullname) }
		*fl=append(*fl,fullname)
		it := col.Iterate()
		for it.Next() != utils.EOC {
			tmp,e := it.Get()
			if e != nil { return e }
			nk := tmp.(*datastore.Key)
			if nk.Equal(k) {
				_,e=it.Evict()
				if e != nil { return e }

				break;
			}
		}
		return nil
	}
}
type z struct{
	io.ReaderAt
	size int64
}
func (zz *z)Size() int64{
	return zz.size
}
func saveZip(c context.Context,f multipart.File,zip bool)([]string,error) {
	keys,err := datastore.NewQuery(FILE_KIND).KeysOnly().GetAll(c,nil)
	filelist :=make([]string,0,10)
	if err != nil { return nil,err}
	col := utils.NewCollection()
	for _,o := range keys{
		col.Add(o)
	}
	//	type zCallback func(io.Reader,os.FileInfo,string) error
	size,err :=f.Seek(0,2)
	if err != nil {
		return nil,err
	}

	err = datastore.RunInTransaction(c,func(cc context.Context) error{
		each := builderEach(cc,col,&filelist)
		if zip {
			return zpack.ZipForEach(&z{f,size},each)
		} else {
			return zpack.TarForEach(f,each)
		}
	},nil)
	if err != nil { return nil,err}
	keys=make([]*datastore.Key, 0, 20)
	it:=col.Iterate()
	for it.Next() != utils.EOC {
		tmp,_ :=it.Get()
		keys = append(keys,tmp.(*datastore.Key))
	}
	if len(keys) > 0 {
		for _,k := range keys{
			log.Infof(c,"delete:%s\n",k.StringID())

		}
		err=datastore.DeleteMulti(c,keys)
		if err != nil {	return nil,err	}
	}
	return filelist,nil
}

func getByName(c context.Context,name string) (*File,error){
	var f File
	sk := ShaStr(name)
	k := datastore.NewKey(c,FILE_KIND,sk,0,getRootKey(c))
	log.Infof(c,"%s",sk)
	err:=datastore.Get(c,k,&f)
	if err != nil { return nil,err}
	return &f,nil
}

const (
	FILE_KIND = "FileBlob"
)
