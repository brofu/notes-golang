# Filed Alignment

## What's Filed Alignment

Let's check the following code 

```
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


s := &MyStruct{}
s2 := &MyStruct2{}

fmt.Printf("\tOffset:    %d\n", unsafe.Offsetof(s.A))
fmt.Printf("\tOffset:    %d\n", unsafe.Offsetof(s.B))
fmt.Printf("\tOffset:    %d\n", unsafe.Offsetof(s.C))

fmt.Printf("\tOffset:    %d\n", unsafe.Offsetof(s2.A))
fmt.Printf("\tOffset:    %d\n", unsafe.Offsetof(s2.B))
fmt.Printf("\tOffset:    %d\n", unsafe.Offsetof(s2.C))
```

>// Offsetof returns the offset within the struct of the field represented by x,
>
>// which must be of the form structValue.field. In other words, it returns the
>
>// number of bytes between the start of the struct and the start of the field.
>
>// The return value of Offsetof is a Go constant if the type of the argument x
>
>// does not have variable size.


As we can see, for the case of `MyStruct`, the `memory padding` is reflected in the offset. It's memory waist definitely. So, the solution of `MyStruct2` is better

## An Example From Golang 

```
// go1.23/src/internal/abi/iface.go

type ITab struct {
	Inter *InterfaceType
	Type  *Type
	Hash  uint32     // copy of Type.Hash. Used for type switches.
	Fun   [1]uintptr // variable sized. fun[0]==0 means Type does not implement Inter.
}

```

The `ITab.Hash` has only 4 bytes. So on the 64bit CPU platform, there would be `memory padding` for it for the memory alignment. But why NOT define `ITab.Func` before `ITab.Hash`, so that it can save around 4 bytes from padding? 

That's because `Fun [1]uintptr` is actually a `placeholder` for a `dynamically sized array of method function pointers`. And the `runtime` would actually append more memory space at the end of the `ITab` struct.

