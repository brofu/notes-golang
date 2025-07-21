package example

import "context"

type GoExample interface {
	Run(context.Context) error
	Notes() string
}
