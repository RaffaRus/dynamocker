package errormsg

import (
	"fmt"
	"runtime"
	"strings"
)

func DynaError(err error) error {
	pc, _, _, _ := runtime.Caller(1)
	res := strings.Split((runtime.FuncForPC(pc).Name()), ".")
	return fmt.Errorf(res[len(res)-1], " :")
}
