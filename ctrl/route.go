package ctrl

import (
	"log"
	"net/http"
	"time"

	"rcgreed.bid/ics/lite"
	"rcgreed.bid/ics/mgr"
	"rcgreed.bid/ics/mgr/session"
	"rcgreed.bid/ics/view"
)

type (
	hdl struct {
	}
	Action func(r *http.Request, se *session.Session) (mgr.View, error)
)

const SESSIONID = "SID"

var dbm *mgr.Manager

func Login(r *http.Request, s *session.Session) (mgr.View, error) {
	name := r.FormValue("name")
	pwd := r.FormValue("password")
	ok, err := dbm.Login(name, pwd)
	data := make(map[string]any)
	data["login"] = []string{name, pwd}
	if err != nil {
		data["error"] = err
	} else {
		if ok {
			data["result"] = true
		} else {
			data["result"] = false
		}
	}
	return view.JsonView(data), nil
}
func (h *hdl) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var sid string
	var ses *session.Session
	cookies, err := r.Cookie(SESSIONID)
	if r.RequestURI == "/favicon.ico" {
		http.NotFound(w, r)
		return
	}
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
	if anws, ok := actionMap[action]; ok {
		log.Print(action)
		v, err := anws(r, ses)
		if err != nil {
			http.NotFound(w, r)
		}
		if err = v(w); err != nil {
			log.Fatal(err)
		}
	} else {
		http.ServeFile(w, r, "web/login.html")
		return
	}

}
func Route() http.Handler {
	dbi, err := lite.NewDBI("home.db")
	if err != nil {
		panic(err)
	}
	dbm = mgr.NewManager(dbi)
	return &hdl{}
}
