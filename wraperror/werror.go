package wraperror

import (
	"errors"
	"fmt"
)
type werror struct{
	err error
	msg string
}
func (w werror)Error() string{
	return fmt.Sprintf("%s\n\t[%s]",w.msg,w.err.Error())
}
func WrapperError(inner error,message string) error{
	if inner == nil {
		return errors.New(message)
	}
	if message == "" {
		return inner
	}
	return &werror{err:inner,msg:message}
}

func Printf(inner error,format string, objs ...interface{})error{
	if inner == nil {
		return fmt.Errorf(format, objs...)
	}
	return &werror{err:inner,msg:fmt.Sprintf(format,objs...)}
}

type MultiError []error
func (me MultiError)Error()string{
	return fmt.Sprintf("total %d error(s)",len(me))
}
