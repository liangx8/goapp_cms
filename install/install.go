package install

import (
	"errors"
	"net/http"
)

type (
	Config interface {
		GetInt(key string) (int, bool)
		GetString(key string) (string, bool)
	}
	cfgimp struct {
		getInt    func(key string) (int, bool)
		getString func(key string) (string, bool)
	}
)

func (c cfgimp) GetInt(key string) (int, bool) {
	return 0, false
}
func (c cfgimp) GetString(key string) (string, bool) {
	return "", false
}

func LoadCfg() (Config, error) {
	return nil, errors.New("ics.cfg is not exists")
}

func HttpHand(w http.ResponseWriter, r *http.Request) {

}
