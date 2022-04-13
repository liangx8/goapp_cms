package lite

import (
	"crypto/sha256"
	"database/sql"
	"fmt"
	"io"
	"reflect"
	"strings"

	_ "github.com/mattn/go-sqlite3"
	"rcgreed.bid/ics/entity"
	"rcgreed.bid/ics/mgr"
	"rcgreed.bid/ics/utils"
)

type (
	dbiImp struct {
		db *sql.DB
	}
)

var pwdKid = utils.NewPasswordKit(sha256.New())

func (my *dbiImp) Load(a any) error {
	//根据不同的数据对象定义数据库表
	switch b := a.(type) {
	case *entity.User:
		return my.loadUser(b)
	case *entity.Item:
		return my.loadItem(b)
	}
	return fmt.Errorf("Object %T unknow", a)
}
func (my *dbiImp) loadUser(u *entity.User) error {
	row := my.db.QueryRow("SELECT seq,name,pwd FROM t_user WHERE seq = ?", u.Seq)
	err := row.Scan(&u.Seq, &u.Name, &u.Password)
	if err != nil {
		return err
	}
	return nil
}
func (my *dbiImp) loadItem(it *entity.Item) error {
	return nil
}

func (my *dbiImp) Close() {
	my.db.Close()
}
func (my *dbiImp) Init() error {
	//my.db.Exec("CREATE TABLE ")

	if _, err := my.db.Exec(CreateSQL(reflect.TypeOf((*entity.User)(nil)).Elem())); err != nil {
		return err
	}
	if _, err := my.db.Exec(CreateSQL(reflect.TypeOf((*entity.Participate)(nil)).Elem())); err != nil {
		return err
	}
	if _, err := my.db.Exec(CreateSQL(reflect.TypeOf((*entity.Item)(nil)).Elem())); err != nil {
		return err
	}

	return nil
}
func NewDBI(dsn string) (mgr.DBI, error) {
	db, err := sql.Open("sqlite3", dsn)
	if err != nil {
		return nil, err
	}
	return &dbiImp{db}, nil
}
func tableNameByEntity(en any) string {
	switch en.(type) {
	case entity.User:
		return "t_user"
	case entity.Participate:
		return "t_participate"
	case entity.Item:
		return "t_item"
	}
	panic("no entity match")

}

var tableByEntityName = map[string]string{
	"User":        "t_user",
	"Participate": "t_participate",
	"Item":        "t_item",
}

func setVal(w io.Writer, t reflect.Type, v reflect.Value) {
	switch t.Kind() {
	default:
		fmt.Fprint(w, v)
	case reflect.Bool:
		if v.Bool() {
			fmt.Fprint(w, 1)
		} else {
			fmt.Fprint(w, 0)
		}

	case reflect.String:
		fmt.Fprintf(w, "'%s'", v.String())
	}
}
func CreateSQL(t reflect.Type) string {
	var cr strings.Builder
	if tn, ok := tableByEntityName[t.Name()]; ok {
		fmt.Fprintf(&cr, "DROP TABLE IF EXISTS %s; CREATE TABLE %s (", tn, tn)
	} else {
		panic("no table coresponding to type")
	}

	for ii := 0; ii < t.NumField(); ii++ {
		f := t.Field(ii)
		if ii > 0 {
			fmt.Fprint(&cr, ",")
		}
		if fn, ok := f.Tag.Lookup("name"); ok {
			fmt.Fprint(&cr, fn)
			switch f.Type.Kind() {
			case reflect.Uint64, reflect.Bool:
				fmt.Fprint(&cr, " INT")
			case reflect.String:
				fmt.Fprintf(&cr, " TEXT")
			case reflect.Struct:
				panic("struct")
			case reflect.Array:
				panic("Array")
			case reflect.Slice:

			}
		}
		_, ok := f.Tag.Lookup("primary")
		if ok {
			fmt.Fprint(&cr, " PRIMARY KEY")
		}
	}
	cr.WriteRune(')')
	return cr.String()
}
func InsertSQL(obj any) string {
	tn := tableNameByEntity(obj)
	ty := reflect.TypeOf(obj)
	va := reflect.ValueOf(obj)
	var fie, val strings.Builder
	fmt.Fprintf(&fie, "INSERT INTO %s (", tn)
	fmt.Fprintf(&val, " VALUES (")
	for ix := 0; ix < ty.NumField(); ix++ {
		if ix > 0 {
			fmt.Fprint(&fie, ",")
			fmt.Fprint(&val, ",")
		}
		fmt.Fprint(&fie, ty.Field(ix).Tag.Get("name"))
		setVal(&val, ty.Field(ix).Type, va.Field(ix))
	}
	fmt.Fprintf(&fie, ")")
	fmt.Fprintf(&val, ")")
	fie.WriteString(val.String())
	return fie.String()
}
func SelectSQL(obj any) string {
	tn := tableNameByEntity(obj)
	var b strings.Builder
	b.WriteString("SELECT * FROM ")
	b.WriteString(tn)
	b.WriteString(" WHERE seq = ?")
	return b.String()
}
