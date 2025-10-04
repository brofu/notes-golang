package defer_examples

import "context"

type example1 struct{}

func (e *example1) Run(ctx context.Context) error {
	return nil
}
func (e *example1) Notes() {

}
