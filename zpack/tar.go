package zpack
import (
	"io"
	"os"
	"archive/tar"
	"compress/gzip"
)
type zCallback func(io.Reader,os.FileInfo,string) error


func TarForEach(r io.Reader,zc zCallback) error {
	zr,err := gzip.NewReader(r)
	if err != nil {
		return err
	}
	tr := tar.NewReader(zr)
	for {
		hdr,err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		err=zc(tr,hdr.FileInfo(),hdr.Name)
		if err != nil { return err }
	}
	return nil
}
