package ctrl

import (
	"log"
	"net/http"
	"time"

	"rcgreed.bid/ics/mgr/session"
)

type (
	hdl struct {
	}
	Action func(r *http.Request) any
)

const SESSIONID = "SID"

func (h *hdl) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var sid string
	var ses *session.Session
	cookies, err := r.Cookie(SESSIONID)
	if r.RequestURI == "/favicon.ico" {
		http.NotFound(w, r)
		return
	}
	log.Println(r.RequestURI)
	if err != nil {
		log.Printf("%s fetch cookies error %v", r.RequestURI, err)
		sid = ""
	} else {
		log.Print(cookies)
		sid = cookies.Value
	}
	sid, ses = session.Get(sid)
	http.SetCookie(w, &http.Cookie{
		Name:       SESSIONID,
		Value:      sid,
		Expires:    time.Now().Add(24 * time.Hour),
		RawExpires: "",
		Path:       "/",
	})
	action := r.FormValue("ACTION")
	if action == "" {
		if _, ok := ses.Data["user"]; !ok {
			http.ServeFile(w, r, "view/login.html")
		}
	} else {
		http.ServeFile(w, r, "view/index.html")
	}

}
func Route() http.Handler {
	return &hdl{}
}
