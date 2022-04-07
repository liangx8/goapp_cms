package lite

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
	"rcgreed.bid/ics/mgr"
)

type (
	dbiImp struct {
		db *sql.DB
	}
)

func (*dbiImp) Load(a any) error {
	return nil
}
func (my *dbiImp) Close() {
	my.db.Close()
}
func NewDBI(dsn string) (mgr.DBI, error) {
	db, err := sql.Open("sqlite3", dsn)
	if err != nil {
		return nil, err
	}
	return &dbiImp{db}, nil
}
