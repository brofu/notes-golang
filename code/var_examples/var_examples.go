package var_examples

import (
	"context"
	"fmt"
	"golang_examples/example"
)

func ShowExample(ctx context.Context) {
}

type VarDeclearExample struct{}

func (v *VarDeclearExample) Run(ctx context.Context) error {

	var (
		a  int
		b  string
		c  VarDeclearExample
		a1 *int
		b1 *string
		c1 *VarDeclearExample
	)

	fmt.Println("=======")
	fmt.Println("case 1 - declare int", a, &a)
	fmt.Println("case 1 - declare string", b, &b)
	fmt.Println("case 1 - declare struct", c, &c)
	fmt.Printf("case 1 - declare struct %p\n", &c)

	fmt.Println("=======")
	fmt.Println("case 2 - declare *int", a1, &a1)
	fmt.Println("case 2 - declare *string", b1, &b1)
	fmt.Println("case 2 - declare *struct", c1, &c1)

	fmt.Println("=======")
	example.PanicWrapper(func() {
		fmt.Println("case 3 - assignment *int - *a1 = 1", a1, &a1)
		*a1 = 1
	})
	a1 = &a
	fmt.Println("case 3 - assignment *int - a1 = &a", a1, &a1)
	example.PanicWrapper(func() {
		fmt.Println("case 3 - assignment *struct - *c1 = VarDeclearExample{}", c1, &c1)
		*c1 = VarDeclearExample{}
	})
	c1 = &c
	fmt.Printf("case 3 - assignment *struct - c1 = &c %p, %p\n", c1, &c1)

	return nil
}

func (v *VarDeclearExample) Notes() string {
	notes := `For reference type vars, need to Assign the memory, before assign value to them.
	Reference Types:	
		* Pointer
		* Map
		* Channel
		* Slice
	Value Types
		* Basic Types. Int, String, Rune, Byte
		* Struct
		* Array
	`

	return notes
}
