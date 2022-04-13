package session

import (
	"log"
	"time"

	"rcgreed.bid/ics/utils"
)

type (
	Session struct {
		expired time.Time
		Data    map[string]any
	}
)

var sessions = make(map[string]*Session)

func Get(id string) (string, *Session) {
	var se *Session
	if se, ok := sessions[id]; ok {
		now := time.Now()
		if now.Before(se.expired) {
			se.expired = time.Now().Add(time.Hour * 24)
			return id, se
		} else {
			delete(sessions, id)
			log.Printf("Session %s is expired\n", id)
		}

	}
	id = utils.MakeID()
	se = &Session{
		expired: time.Now().Add(time.Hour * 24),
		Data:    make(map[string]any),
	}
	sessions[id] = se
	log.Printf("New session was created %s", id)
	return id, se
}
