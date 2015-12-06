package zpack

import (
	"mime/multipart"
	"archive/zip"

)
func ZipForEach(mpf multipart.File,zc zCallback) error {
	size,err :=mpf.Seek(0,2)
	if err != nil {
		return err
	}
	zr,err := zip.NewReader(mpf,size)
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
