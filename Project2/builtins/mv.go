package builtins

import (
	"errors"
	"fmt"
	"os"
)

var (
	ModeDir6 = 0
	ErrInvalidArgCount6 = errors.New("invalid argument count")
)

func Renamefile(args ...string) error {
	switch len(args) {
	case 1:
               return fmt.Errorf("%w: expected two arguments (directory)", ErrInvalidArgCount6) 
        case 2:
              return os.Rename(args[0], args[1]) 
	default:
		return fmt.Errorf("%w: expected two arguments (directory)", ErrInvalidArgCount6)
	}
}
