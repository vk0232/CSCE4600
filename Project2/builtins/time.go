package builtins

import (
        "errors"
	"fmt"
        "time"
)

var (
	ModeDir5 = 0
	ErrInvalidArgCount5 = errors.New("invalid argument count")
)

func Gettime(args ...string)  {
    // using the function 
    mytime := time.Now() 
    fmt.Println(mytime)  
}
