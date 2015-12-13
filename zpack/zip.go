package zpack

import (
	"io"
	"archive/zip"



)
//type Zcallback func(io.Reader,os.FileInfo,string) error
type ZipReaderAdaptor interface{
	io.ReaderAt
	Size() int64
}
// ZipReaderAdaptor implement
type zra struct{
	size int64
	bytes [][]byte
}
func (z zra) ReadAt(p []byte, off int64) (int,error){
	offidx := off / blocksize
	offtail := off % blocksize
	i :=0
	for i,_ = range p {
		if offidx * blocksize + offtail >= z.size {
			return i,io.EOF
		}
		p[i]=z.bytes[offidx][offtail]
		offtail ++
		if offtail >= blocksize {
			offtail = 0
			offidx ++
		}
	}
	return i+1,nil
}
func (z zra) Size() int64{
	return z.size
}

// src source
// size the size of ZipReaderAdapotr
func NewZipReaderAdaptor(src io.Reader,size int64) (ZipReaderAdaptor,error) {

	bytes := make([][]byte,0)
	count := size / blocksize
	tail := size % blocksize
	var idx int64
	idx = 0

	for {
//		fmt.Printf("size:%d,count %d,tail %d,idx: %d,len(bytes): %d\n",size,count,tail,idx,len(bytes))
		if idx == count {
			buf := make([]byte,tail)
			n,err := src.Read(buf)


			if err != nil && err != io.EOF {
				return nil,err
			}
			size = count * blocksize + int64(n)
			bytes = append(bytes,buf)
			break
		}
		buf := make([]byte,blocksize)
		n,err := src.Read(buf)
		if err == io.EOF {
			bytes = append(bytes,buf)
			size = idx * blocksize + int64(n)
			break
		}

		if err != nil {
			return nil,err
		}
		bytes = append(bytes,buf)
		//if idx >= count {break}
		idx ++

	}
//	fmt.Printf("len(bytes) %d,size: %d\n",len(bytes),size)
//	for i,v := range bytes {
//		fmt.Printf("%d=>%d\n",i,len(v))
//	}
	return &zra{bytes:bytes,size:size},nil
}
func ZipForEach(mpf ZipReaderAdaptor,zc Zcallback) error {
	zr,err := zip.NewReader(mpf,mpf.Size())
	if err != nil {
		return err
	}
	for _,f := range zr.File {
		rc, err := f.Open()
		defer rc.Close()
		if err != nil {
			return err
		}

		err=zc(rc,f.FileInfo(),f.Name)
		if err != nil {return err}

	}
	return nil
}
const (
	blocksize int64 = 1024
)
