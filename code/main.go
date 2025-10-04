package main

import (
	"context"
	"golang_examples/var_examples"
)

// DebugTypeName reads a runtime internal type pointer (uintptr)
// and prints the Go type name using reflect.
/*
func DebugTypeName(ptr uintptr) {
	rtyp := (*reflect.rtype)(unsafe.Pointer(ptr))
	fmt.Println("Type is:", reflect.TypeOf(reflect.NewAt(rtyp, unsafe.Pointer(ptr)).Elem()).String())
}
*/

func main() {
	ctx := context.Background()
	//interface_examples.ShowExample(context.Background())
	//struct_examples.ShowExample(context.Background())
	(&var_examples.VarDeclearExample{}).Run(ctx)
}
