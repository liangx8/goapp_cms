package mgr

import (
	"crypto/sha1"
	"io"

	"rcgreed.bid/ics/entity"
	"rcgreed.bid/ics/utils"
)

type (
	DBI interface {
		Load(d any) error
		Init(entity.User) error
		Add(o any) error
		Save(o any) error
		// name 查询条件，结果填充回user
		GetUserByName(user *entity.User, name string) error
		Close()
	}
	Condition interface {
	}
	View    func(w io.Writer) error
	Manager struct {
		DBI
	}
)

var pwdKid = utils.NewPasswordKit(sha1.New())

func NewManager(dbi DBI) *Manager {
	return &Manager{dbi}
}
func (mgr *Manager) Login(name, password string) (bool, error) {
	var ur entity.User
	if err := mgr.GetUserByName(&ur, name); err != nil {
		return false, err
	}
	return pwdKid.Verify(ur.Pwd, password), nil
}
