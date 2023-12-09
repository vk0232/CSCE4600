package builtins

import (
	"errors"
	"fmt"
	"os"
)

var (
	ModeDir = 0
	ErrInvalidArgCount1 = errors.New("invalid argument count")
)

func MakeDirectory(args ...string) error {
	switch len(args) {
	case 1:
                return os.Mkdir(args[0], 0777)    
	default:
		return fmt.Errorf("%w: expected one or two arguments (directory)", ErrInvalidArgCount1)
	}
}
