package mgr

import (
	"errors"
	"io"

	"rcgreed.bid/ics/entity"
)

type (
	DBI interface {
		Load(d any) error
		Init() error
		Close()
	}
	View    func(w io.Writer, data any) error
	Manager struct {
		DBI
	}
)

func NewManager(dbi DBI) *Manager {
	return &Manager{dbi}
}
func (mgr *Manager) Login(usr entity.User) error {
	db := entity.User{Seq: usr.Seq}
	if err := mgr.Load(db); err != nil {
		return err
	}
	if db.Password == usr.Password {
		return nil
	}
	return loginFail
}

var loginFail = errors.New("User name or password does not match")
