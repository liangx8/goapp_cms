package one

import(
	"net/http"
	"gcs"
)

func init(){
	http.HandleFunc("/showdoor",backEnd)

	http.HandleFunc("/gcsb",gcs.Backend)
	http.HandleFunc("/",gcs.Front)

//	http.HandleFunc("/",front)

}
