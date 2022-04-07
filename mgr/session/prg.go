package session

import "context"

type (
	Session struct {
		ctx context.Context
	}
)

var sessions = make(map[string]map[string]Session)

func Get(id string) *Session {
	return nil
}
