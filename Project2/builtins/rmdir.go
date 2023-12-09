package builtins

import (
	"errors"
	"fmt"
	"os"
)

var (
	ModeDir4 = 0
	ErrInvalidArgCount2 = errors.New("invalid argument count")
)

func RemoveDirectory(args ...string) error {
	switch len(args) {
	case 1:
                return os.Remove(args[0])    
	default:
		return fmt.Errorf("%w: expected one or two arguments (directory)", ErrInvalidArgCount2)
	}
}
