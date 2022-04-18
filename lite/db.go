/*
 * 每个表的主键使用了rowid,因此在go中定义的主键类型必须是数字整型
 */

package lite

import (
	"database/sql"
	"fmt"
	"io"
	"reflect"
	"strings"

	_ "github.com/mattn/go-sqlite3"
	"rcgreed.bid/ics/entity"
	"rcgreed.bid/ics/mgr"
)

type (
	dbiImp struct {
		db *sql.DB
	}
)

func MakeDbError(sql string, err error) error {
	return fmt.Errorf("error [%v] at SQL:%s", err, sql)
}

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

	row := my.db.QueryRow("SELECT name,pwd FROM t_user WHERE rowid = ?", u.Seq)
	err := row.Scan(&u.Name, &u.Pwd)
	if err != nil {
		return err
	}
	return nil
}
func (my *dbiImp) GetUserByName(usr *entity.User, name string) error {
	row := my.db.QueryRow("SELECT name,pwd FROM t_user WHERE name = ? AND active = 1", name)
	if err := row.Scan(&usr.Name, &usr.Pwd); err != nil {
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
func (my *dbiImp) Init(admin entity.User) error {
	//my.db.Exec("CREATE TABLE ")
	sql := createSQL(reflect.TypeOf((*entity.User)(nil)).Elem())
	if _, err := my.db.Exec(sql); err != nil {
		return MakeDbError(sql, err)
	}
	sql = createSQL(reflect.TypeOf((*entity.Participate)(nil)).Elem())
	if _, err := my.db.Exec(sql); err != nil {
		return MakeDbError(sql, err)
	}
	sql = createSQL(reflect.TypeOf((*entity.Item)(nil)).Elem())
	if _, err := my.db.Exec(sql); err != nil {
		return MakeDbError(sql, err)
	}
	if err := my.Add(admin); err != nil {
		return err
	}

	return nil
}
func (my *dbiImp) Add(d any) error {
	switch obj := d.(type) {
	case entity.User:
		if err := addUser(my.db, obj); err != nil {
			return err
		}

	}
	return nil
}

func addUser(dB *sql.DB, usr entity.User) error {
	if _, err := dB.Exec(insertSQL(usr)); err != nil {
		return err
	}
	return nil
}
func (my *dbiImp) Save(d any) error {
	switch obj := d.(type) {
	case entity.User:
		if _, err := my.db.Exec(fmt.Sprintf("UPDATE t_user set pwd=x'%x',active = %t WHERE rowid = ?", obj.Pwd, obj.Active), obj.Seq); err != nil {
			return err
		}
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
	case reflect.Slice:
		// 期望类型是byte
		fmt.Fprintf(w, "x'%x'", v.Bytes())
	case reflect.String:
		fmt.Fprintf(w, "'%s'", v.String())
	}
}
func createSQL(t reflect.Type) string {
	var cr strings.Builder
	if tn, ok := tableByEntityName[t.Name()]; ok {
		fmt.Fprintf(&cr, "DROP TABLE IF EXISTS %s; CREATE TABLE %s (", tn, tn)
	} else {
		panic("no table coresponding to type")
	}
	first := true
	for ii := 0; ii < t.NumField(); ii++ {
		f := t.Field(ii)
		if _, ok := f.Tag.Lookup("primary"); ok {
			// 用rowid来代替
			continue
		}
		if fn, ok := f.Tag.Lookup("name"); ok {
			if first {
				first = false
			} else {
				fmt.Fprint(&cr, ",")
			}
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
				fmt.Fprintf(&cr, " BLOB")
			}
		}
	}
	cr.WriteRune(')')
	return cr.String()
}
func insertSQL(obj any) string {
	tn := tableNameByEntity(obj)
	ty := reflect.TypeOf(obj)
	va := reflect.ValueOf(obj)
	var fie, val strings.Builder
	fmt.Fprintf(&fie, "INSERT INTO %s (", tn)
	fmt.Fprintf(&val, " VALUES (")
	first := true
	for ix := 0; ix < ty.NumField(); ix++ {
		// 主键用rowid，因此无需建立主键字段
		if _, ok := ty.Field(ix).Tag.Lookup("primary"); ok {
			continue
		}
		if fn, ok := ty.Field(ix).Tag.Lookup("name"); ok {
			if first {
				first = false
			} else {
				fmt.Fprint(&fie, ",")
				fmt.Fprint(&val, ",")
			}
			fmt.Fprint(&fie, fn)
			setVal(&val, ty.Field(ix).Type, va.Field(ix))
		}

	}
	fmt.Fprintf(&fie, ")")
	fmt.Fprintf(&val, ")")
	fie.WriteString(val.String())
	return fie.String()
}
func TinsertSQL(obj any) string {
	return insertSQL(obj)
}
func SelectSQL(obj any) string {
	tn := tableNameByEntity(obj)
	var b strings.Builder
	b.WriteString("SELECT * FROM ")
	b.WriteString(tn)
	b.WriteString(" WHERE rowid = ?")
	return b.String()
}
