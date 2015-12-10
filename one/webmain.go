package one

import(
	"net/http"

)

func init(){
	http.HandleFunc("/showdoor",backEnd)

	http.HandleFunc("/",front)
}
