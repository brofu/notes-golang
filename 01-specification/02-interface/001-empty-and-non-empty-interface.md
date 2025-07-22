## Empty V.S. Non-Empty Interface

### Concepts

```
// empty interface
interface {}

// non-empty interface
interface {
  Method() 
}

```

As we know, in Go, an interface variable is a two-part structure:

1. `Type`: A pointer to the concrete type of the value it holds.
2. `Value`: A pointer to the actual data.

But, what's the story under layer? How does an `empty` and `non-empty` store such above mentioned 2 parts?


### The Under Layer Struct for Empty and Non-Empty Interface

Let's check the source code with version `v1.23`

```
// go1.23/src/runtime/runtime2.go

type itab = abi.ITab
type _type = abi.Type

type iface struct { // presentation of non-empty interface
	tab  *itab
	data unsafe.Pointer
}

type eface struct { // presentation of empty interface
	_type *_type
	data  unsafe.Pointer
}

// go1.23/src/internal/abi/type.go

// Type is the runtime representation of a Go type.
//
// Be careful about accessing this type at build time, as the version
// of this type in the compiler/linker may not have the same layout
// as the version in the target binary, due to pointer width
// differences and any experiments. Use cmd/compile/internal/rttype
// or the functions in compiletype.go to access this type instead.
// (TODO: this admonition applies to every type in this package.
// Put it in some shared location?)
type Type struct {
	Size_       uintptr
	PtrBytes    uintptr // number of (prefix) bytes in the type that can contain pointers
	Hash        uint32  // hash of type; avoids computation in hash tables
	TFlag       TFlag   // extra type information flags
	Align_      uint8   // alignment of variable with this type
	FieldAlign_ uint8   // alignment of struct field with this type
	Kind_       Kind    // enumeration for C
	// function for comparing objects of this type
	// (ptr to object A, ptr to object B) -> ==?
	Equal func(unsafe.Pointer, unsafe.Pointer) bool
	// GCData stores the GC type data for the garbage collector.
	// If the KindGCProg bit is set in kind, GCData is a GC program.
	// Otherwise it is a ptrmask bitmap. See mbitmap.go for details.
	GCData    *byte
	Str       NameOff // string form
	PtrToThis TypeOff // type for pointer to this type, may be zero
}

// go1.23/src/internal/abi/iface.go

// The first word of every non-empty interface type contains an *ITab.
// It records the underlying concrete type (Type), the interface type it
// is implementing (Inter), and some ancillary information.
//
// allocated in non-garbage-collected memory
type ITab struct {
	Inter *InterfaceType
	Type  *Type
	Hash  uint32     // copy of Type.Hash. Used for type switches.
	Fun   [1]uintptr // variable sized. fun[0]==0 means Type does not implement Inter.
}
```

From the code snippet, we can tell that,

1. For an `empty` and `non-empty` interface, there is a `data` filed, which is used to store the pointer to the `Value` 
2. For an `empty` interface, there is a field `_type` to keep the concrete `Type` of the data it holds    
2. For an `non-empty` interface, there is a `ITab` field, in which there are some important fields:
  * `Inter`. Keeps the specific `interface type`
  * `Type`. Same as the `Type` in `eface`, keep the concrete `Type` of the data it holds
  * `Fun`. Keeps the functions (an array of functions) implemented by the data. Yes, you are right, all the methods defined by the interface should be present here.

### More Details about `Non-Empty` Interface




### References

1. Golang Source Code V1.23
