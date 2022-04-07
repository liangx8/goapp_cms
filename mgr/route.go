package mgr

import (
	"log"
	"net/http"
	"time"
)

type (
	hdl struct {
		act  Action
		view View
	}
	Action func(r *http.Request) any
)

const SESSIONID = "SID"

func (h *hdl) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	cookies, err := r.Cookie(SESSIONID)
	if err != nil {
		log.Printf("%s fetch cookies error %v", r.RequestURI, err)
		http.SetCookie(w, &http.Cookie{
			Name:       SESSIONID,
			Value:      "akdksjk34k",
			Expires:    time.Now().Add(24 * time.Hour),
			RawExpires: "",
			Path:       "/",
		})
	} else {
		log.Print(cookies)
		log.Print(time.Now().Add(time.Duration(cookies.Expires.Unix())))
	}
	h.view(w, h.act(r))

}
func Filter(hf Action, view View) http.Handler {
	return &hdl{hf, view}
}
