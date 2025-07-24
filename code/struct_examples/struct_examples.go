package struct_examples

import (
	"context"
	"fmt"
	"golang_examples/code/example"
	"unsafe"
)

type MyStruct struct {
	A byte
	B int64
	C byte
}

type MyStruct2 struct {
	A int64
	B byte
	C byte
}

type StructFieldAlignmentExample struct{}

func (sa *StructFieldAlignmentExample) Notes() string {

	return ""
}

func (sa *StructFieldAlignmentExample) Run(ctx context.Context) error {

	s := &MyStruct{}

	fmt.Printf("Case need padding\n")
	fmt.Printf("Field A (byte):\n")
	fmt.Printf("\tSize:      %d\n", unsafe.Sizeof(s.A))
	fmt.Printf("\tOffset:    %d\n", unsafe.Offsetof(s.A))
	fmt.Printf("\tAlignment: %d\n\n", unsafe.Alignof(s.A))

	fmt.Printf("Field B (int64):\n")
	fmt.Printf("\tSize:      %d\n", unsafe.Sizeof(s.B))
	fmt.Printf("\tOffset:    %d\n", unsafe.Offsetof(s.B))
	fmt.Printf("\tAlignment: %d\n\n", unsafe.Alignof(s.B))

	fmt.Printf("Field C (byte):\n")
	fmt.Printf("\tSize:      %d\n", unsafe.Sizeof(s.C))
	fmt.Printf("\tOffset:    %d\n", unsafe.Offsetof(s.C))
	fmt.Printf("\tAlignment: %d\n\n", unsafe.Alignof(s.C))

	s2 := &MyStruct2{}

	fmt.Printf("Case NO need padding\n")
	fmt.Printf("Field A (byte):\n")
	fmt.Printf("\tSize:      %d\n", unsafe.Sizeof(s2.A))
	fmt.Printf("\tOffset:    %d\n", unsafe.Offsetof(s2.A))
	fmt.Printf("\tAlignment: %d\n\n", unsafe.Alignof(s2.A))

	fmt.Printf("Field B (int64):\n")
	fmt.Printf("\tSize:      %d\n", unsafe.Sizeof(s2.B))
	fmt.Printf("\tOffset:    %d\n", unsafe.Offsetof(s2.B))
	fmt.Printf("\tAlignment: %d\n\n", unsafe.Alignof(s2.B))

	fmt.Printf("Field C (byte):\n")
	fmt.Printf("\tSize:      %d\n", unsafe.Sizeof(s2.C))
	fmt.Printf("\tOffset:    %d\n", unsafe.Offsetof(s2.C))
	fmt.Printf("\tAlignment: %d\n\n", unsafe.Alignof(s2.C))

	return nil
}

func ShowExample(ctx context.Context) {

	examples := []example.GoExample{&StructFieldAlignmentExample{}}

	for _, exam := range examples {

		fmt.Println(exam.Notes())
		if err := exam.Run(ctx); err != nil {
			fmt.Println(err)
		}
	}

}
