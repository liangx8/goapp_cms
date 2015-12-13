package utils

import (
	"strings"
	"regexp"
	"crypto/sha1"
	"fmt"
)

var pat *regexp.Regexp
func init(){
	var err error
	pat,err=regexp.Compile("\\.[0-9a-zA-Z]*?$")
	if err !=nil {panic(err)}
}
func GuessMimeType(name string)string{
	str:=strings.ToLower(name)


	a:=pat.FindStringIndex(str)
	if len(a)==0{
		return "application/octet-stream"
	}
	switch(str[a[0]:a[1]]){
	case ".svg":
		return "image/svg+xml"
	case ".txt":
		return "text/plain"
	case ".bmp":
		return "image/bmp"
	case ".png":
		return "image/png"
	case ".gif":
		return "image/gif"
	case ".js":
		return "application/javascript"
	case ".css":
		return "text/css"
	case ".html",".htm":
		return "text/html"
	case ".jpg":
		return "image/jpeg"
	case ".json":
		return "application/json"
	case ".xml":
		return "applicaton/xml"
	case ".ico":
		return "image/x-icon"
	default:
		return "application/octet-stream"
	}

}

func ShaStr(s string) string{
	sha1.Sum([]byte(s))
	return fmt.Sprintf("%x",sha1.Sum([]byte(s)))
}
