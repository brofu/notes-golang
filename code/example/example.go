package example

import (
	"context"
	"fmt"
)

type GoExample interface {
	Run(context.Context) error
	Notes() string
}

func PanicWrapper(f func()) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("panic: ", err)
		}
	}()
	f()
}
