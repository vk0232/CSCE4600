package builtins

import (
	"errors"
	"fmt"
	"os"
)

var (
	ModeDir2 = 0
	ErrInvalidArgCount3 = errors.New("invalid argument count")
)

func GetworkDirectory(args ...string)  {
    mydir, err := os.Getwd() 
    if err != nil { 
        fmt.Println(err) 
    }
        fmt.Println(mydir)  
}
