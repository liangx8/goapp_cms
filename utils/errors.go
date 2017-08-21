package utils

import (
	"fmt"
)

func WrappingError(parent error,msg string) error{
	return fmt.Errorf("%s\nparent's error:%v",msg,parent)
}
