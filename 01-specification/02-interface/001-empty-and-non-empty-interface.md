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

### How Is the Method-Call for a `Non-Empty` Interface Processed

So, for an `non-empty interface`, what's workflow when we call an method based on it? Let's check this code 

```
type People interface {
	Name() string
	Work() string
}

type Teacher struct{}

func (t *Teacher) Name() string {
	return "teacher"
}

func (t *Teacher) Work() string {
	return "teach"
}

type Student struct{}

func (s *Student) Name() string {
	return "student"
}

func (s *Student) Work() string {
	return "study"
}

func showExampleMethodCallV2(ctx context.Context) {
	for _, p := range []People{&Teacher{}, &Student{}} {
		fmt.Println("method call:", p.Work())
		fmt.Println("method call:", p.Name())
	}
}
```

Please note, there are 2 key points of these code snippet. And we would explain this later.

1. **There is a slice of `People`, with 2 different concrete types, say that, `Student` and `Teacher`.**
2. **`Work()` method is defined secondly, but called firstly, while `Name()` is reverse.** 

OK. In short, the workflow would be (roughly) like this:
1. Get the method pointer (`Name()` or `Work()`) from the `iface.*table.Fun`
3. Get the concrete data pointer by `iface.data`
4. Call the method with concrete data as the 1st parameter

Let's verify this by debugging the above code snippet.

The debugging info showed as following, when calling `p.Work()`. Pay attention to the comments inline.

```
        interface_exmaples.go:106       0x1024c99c4     e18b40f9        MOVD 272(RSP), R1
        interface_exmaples.go:106       0x1024c99c8     220040f9        MOVD (R1), R2   // R2 get `itab` data, by pointer
        interface_exmaples.go:106       0x1024c99cc     200440f9        MOVD 8(R1), R0  // R0 get the `iface.data`, by pointer
        interface_exmaples.go:106       0x1024c99d0     e22300f9        MOVD R2, 64(RSP)
        interface_exmaples.go:106       0x1024c99d4     e02700f9        MOVD R0, 72(RSP)
        interface_exmaples.go:107       0x1024c99d8*    5b008039        MOVB (R2), R27
=>      interface_exmaples.go:107       0x1024c99dc     411040f9        MOVD 32(R2), R1 // R1 the method `Work()`, by pointer. Why it's (R2) + 32?
        interface_exmaples.go:107       0x1024c99e0     20003fd6        CALL (R1) // call the `Work()` method, with R0 as the 1st parameter. 
```

Note: 
> On Apple Silicon (ARM64), Go follows the Plan 9 ABI, where the first argument to a function is passed in R0.

> -- By ChatGPT


When call `p.Name()`. Pay attention to the comments inline.

```
> [Breakpoint 3] golang_examples/interface_examples.showExampleMethodCallV2() ./interface_examples/interface_exmaples.go:108 (hits goroutine(1):1 total:1) (PC: 0x1024c9ab0)
   103: }
   104:
   105: func showExampleMethodCallV2(ctx context.Context) {
   106:         for _, p := range []People{&Teacher{}, &Student{}} {
   107:                 fmt.Println("method call:", p.Work())
=> 108:                 fmt.Println("method call:", p.Name())
   109:         }
   110: }
   111:
   112: func showExampleMethodCall(ctx context.Context) {
   113:         var stu People = &Student{}
(dlv) si
> golang_examples/interface_examples.showExampleMethodCallV2() ./interface_examples/interface_exmaples.go:108 (PC: 0x1024c9ab4)
        interface_exmaples.go:107       0x1024c9aa0     e28300f9        MOVD R2, 256(RSP)
        interface_exmaples.go:107       0x1024c9aa4     e28700f9        MOVD R2, 264(RSP)
        interface_exmaples.go:107       0x1024c9aa8     e10302aa        MOVD R2, R1
        interface_exmaples.go:107       0x1024c9aac     15edff97        CALL fmt.Println(SB)
        interface_exmaples.go:108       0x1024c9ab0*    e32340f9        MOVD 64(RSP), R3
=>      interface_exmaples.go:108       0x1024c9ab4     7b008039        MOVB (R3), R27
        interface_exmaples.go:108       0x1024c9ab8     630c40f9        MOVD 24(R3), R3  // R3 now holds the pointer of `Name()` method. Why the address is (R3)+24? 
        interface_exmaples.go:108       0x1024c9abc     e02740f9        MOVD 72(RSP), R0 // R0 get the `iface.data` from (RSP) + 72, which is set when process `p.Work()`
        interface_exmaples.go:108       0x1024c9ac0     60003fd6        CALL (R3) // call the `Name()` method, with R0 as the 1st parameter.
        interface_exmaples.go:108       0x1024c9ac4     e04700f9        MOVD R0, 136(RSP)
        interface_exmaples.go:108       0x1024c9ac8     e14b00f9        MOVD R1, 144(RSP)
(dlv)
```

OK, let's back to the 2 import points mentioned previously.

* **Why there is a slice of `People`, with 2 different concrete types, say that, `Student` and `Teacher`?**

The main purpose for this, is to let the code `common` (for example, to seek the method address by address offset), instead of by `function symbol`.

For example, if we write code like this, the code would be presented by `function symbol`, and as a result, we cannot observe the `address offset` clearly.

```
func showExampleMethodCall(ctx context.Context) {
	var stu People = &Student{} 
	fmt.Println("method call:", stu.Name()) // the compiler can get the address of name statically 
}
```

And if we debug this line, we would get this. Pay attention to the `CALL` instruction. 

```
kpoint 3] golang_examples/interface_examples.showExampleMethodCall() ./interface_examples/interface_exmaples.go:114 (hits goroutine(1):1 total:1) (PC: 0x102199310)
   109:         }
   110: }
   111:
   112: func showExampleMethodCall(ctx context.Context) {
   113:         var stu People = &Student{}
=> 114:         fmt.Println("method call:", stu.Name())
   115: }
   116:
   117: func showExample(ctx context.Context) {
   118:
   119:         examples := []example.GoExample{&NilInterfaceExample{}}
(dlv) si
> golang_examples/interface_examples.showExampleMethodCall() ./interface_examples/interface_exmaples.go:114 (PC: 0x102199314)
        interface_exmaples.go:113       0x102199300     a10100d0        ADRP 221184(PC), R1
        interface_exmaples.go:113       0x102199304     21801691        ADD $1440, R1, R1
        interface_exmaples.go:113       0x102199308     e11700f9        MOVD R1, 40(RSP)
        interface_exmaples.go:113       0x10219930c     e01b00f9        MOVD R0, 48(RSP)
        interface_exmaples.go:114       0x102199310*    01000014        JMP 1(PC)
=>      interface_exmaples.go:114       0x102199314     e02f00f9        MOVD R0, 88(RSP)
        interface_exmaples.go:114       0x102199318     caffff97        CALL golang_examples/interface_examples.(*Student).Name(SB) // Show by function symbol.
        interface_exmaples.go:114       0x10219931c     e04700f9        MOVD R0, 136(RSP)
        interface_exmaples.go:114       0x102199320     e14b00f9        MOVD R1, 144(RSP)
        interface_exmaples.go:114       0x102199324     ffff06a9        STP (ZR, ZR), 104(RSP)
        interface_exmaples.go:114       0x102199328     ffff07a9        STP (ZR, ZR), 120(RSP)
```

As we can see, the address of `Name()` is presented by the `a Function Symbol`. And there is NO procedure of `getting the address` triggered.

* **`Work()` method is defined secondly, but called firstly, while `Name()` is reverse**. 

This is key reason why the address of `Name()` is `24(R3)`, while that of `Work()` is `32(R2)`. Let's back to the definition of `runtime.itab`

```
type ITab struct {
	Inter *InterfaceType
	Type  *Type
	Hash  uint32     // copy of Type.Hash. Used for type switches.
	Fun   [1]uintptr // variable sized. fun[0]==0 means Type does not implement Inter.
}
```

After `memory padding`, the actual layout of `ITab` is like this 

```
type ITab struct {
	Inter *Interfacetype // 8 bytes, offset 0
	Type *Type           // 8 bytes, offset 8
	hash  uint32         // 4 bytes, offset 16
	_     [4]byte        // 4 bytes, offset 20, padding for alignment
	Fun   [1]uintptr     // offset 24, start of method pointers, its length depends on the method number of the interface, and each pointer has the size of 8 bytes
}
```

So, the offset of the 1st method (in the order of definition in the interface, here, it's `People`) `Name()` is `start_address` + 24, and 

The offset of the 2nd method `Work()` is `start_address` + 32.

### Another Story For Method Call to `Non-Empty` Interface

We just talked about the way of method call on `non-empty interface` like this 

```
func showExampleMethodCallV2(ctx context.Context) {
	for _, p := range []People{&Teacher{}, &Student{}} {
		fmt.Println("method call:", p.Work())
		fmt.Println("method call:", p.Name())
	}
}
```

Now, let's check a different one, say that, 

```
func showExampleMethodCallV3(ctx context.Context) {
	for _, p := range []interface{}{&Teacher{}, &Student{}} {
        people := p.(People)
		fmt.Println("method call:", people.Work())
		fmt.Println("method call:", people.Name())
	}
}
```

The difference of this version is that, the slice type is actually `[]interface{}`, and before the method `Name()` or`Work()` is called, a type assertion is necessary.

Let's debug this to find out the under layer story. (To do this, we need to setup a break point at `runtime.getitab`)

```
(dlv) b runtime.typeAssert
Breakpoint 3 set at 0x102fee27c for runtime.typeAssert() /Users/jeff_shao/.gvm/gos/go1.23/src/runtime/iface.go:467
(dlv) b runtime.assertE2I
Breakpoint 4 set at 0x102fee17c for runtime.assertE2I() /Users/jeff_shao/.gvm/gos/go1.23/src/runtime/iface.go:449
(dlv) b runtime.assertE2I2
Breakpoint 5 set at 0x102fee21c for runtime.assertE2I2() /Users/jeff_shao/.gvm/gos/go1.23/src/runtime/iface.go:457
(dlv) b 107
Breakpoint 6 set at 0x1030999d8 for golang_examples/interface_examples.showExampleMethodCallV3() ./interface_examples/interface_exmaples.go:107
(dlv) c
> [Breakpoint 6] golang_examples/interface_examples.showExampleMethodCallV3() ./interface_examples/interface_exmaples.go:107 (hits goroutine(1):1 total:1) (PC: 0x1030999d8)
   102:         showExampleMethodCallV3(ctx)
   103: }
   104:
   105: func showExampleMethodCallV3(ctx context.Context) {
   106:         for _, p := range []interface{}{&Teacher{}, &Student{}} {
=> 107:                 people := p.(People)
   108:                 fmt.Println("method call:", people.Work())
   109:                 fmt.Println("method call:", people.Name())
   110:         }
   111: }
   112:
(dlv) c
> [Breakpoint 3] runtime.typeAssert() /Users/jeff_shao/.gvm/gos/go1.23/src/runtime/iface.go:467 (hits goroutine(1):1 total:1) (PC: 0x102fee27c)
Warning: debugging optimized function
   462: }
   463:
   464: // typeAssert builds an itab for the concrete type t and the
   465: // interface type s.Inter. If the conversion is not possible it
   466: // panics if s.CanFail is false and returns nil if s.CanFail is true.
=> 467: func typeAssert(s *abi.TypeAssert, t *_type) *itab {  // Key Point 1. `runtime.typeAssert` is involved.
   468:         var tab *itab
   469:         if t == nil {
   470:                 if !s.CanFail {
   471:                         panic(&TypeAssertionError{nil, nil, &s.Inter.Type, ""})
   472:                 }
(dlv) c
> [Breakpoint 2] runtime.getitab() /Users/jeff_shao/.gvm/gos/go1.23/src/runtime/iface.go:44 (hits goroutine(1):1 total:1) (PC: 0x10304604c)
Warning: debugging optimized function
    39: //
    40: // Do not remove or change the type signature.
    41: // See go.dev/issue/67401.
    42: //
    43: //go:linkname getitab
=>  44: func getitab(inter *interfacetype, typ *_type, canfail bool) *itab { // Key Point 2. `runtime.getitab` is involved.
    45:         if len(inter.Methods) == 0 {
    46:                 throw("internal error - misuse of itab")
    47:         }
    48:
    49:         // easy case

(dlv) p inter
("*internal/abi.InterfaceType")(0x102b5b9e0)
*internal/abi.InterfaceType {
        Type: internal/abi.Type {Size_: 16, PtrBytes: 16, Hash: 2321956729, TFlag: TFlagUncommon|TFlagExtraStar|TFlagNamed (7), Align_: 8, FieldAlign_: 8, Kind_: Interface (20), Equal: runtime.interequal, GCData: *2, Str: 22216, PtrToThis: 25056},
        PkgPath: internal/abi.Name {Bytes: *0},
        Methods: []internal/abi.Imethod len: 2, cap: 2, [
                (*"internal/abi.Imethod")(0x102b5ba40),
                (*"internal/abi.Imethod")(0x102b5ba48),
        ],}
(dlv) p firstmoduledata.types + 22216
4340389576
(dlv) x -fmt hex -count 1 -size 1 4340389577
0x102b516c9:   0x1a
(dlv) x -fmt raw -count 26 -size 1 4340389578 // Key Point 3. check the type name of `inter`, the 1st parameter of `runtime.getitab`. It's actually `People`
*interface_examples.People(dlv)
(dlv) p typ
("*internal/abi.Type")(0x102b5a660)
*internal/abi.Type {Size_: 8, PtrBytes: 8, Hash: 2452649489, TFlag: TFlagUncommon|TFlagRegularMemory (9), Align_: 8, FieldAlign_: 8, Kind_: Int|Int16|Complex128|KindDirectIface (54), Equal: runtime.memequal64, GCData: *1, Str: 22497, PtrToThis: 0}
(dlv) p firstmoduledata.types + 22497
4340389857
(dlv) x -fmt hex -count 1 -size 1 4340389858
0x102b517e2:   0x1b
(dlv) x -fmt raw -count 27 -size 1 4340389859 // Key Point 4. check the type name of `typ`, the 2nd parameter of `runtime.getitab`. It's actually `Teacher`
*interface_examples.Teacher(dlv)
```

We can find the following facts from the debugging info. (The key point are marked as `Key Point x` in the comments)

1. **When the `people := p.(People)` is called, the `runtime.typeAssert()` function is triggered. And following the `runtime.getitab()` function.** (Refer to Key Point 1 and Key Point 2)
2. **Within the process of `runtime.getitab`, if we check the `type` of the 2 parameters `inter` and `typ`, we can tell that, 
    * The `inter` is actually the type of interface `People` and
    * the `typ` is actually the type of `Teacher`**

Yes, that's the 2nd story, when there is type assertion relevant call, such as `people := p.(People)`. 

In this story function of `runtime.getitab(inter *interfacetype, typ *_type, canfail bool) *itab`, would be involved, to determine if the assertion succeed, and get the `itab` object, with which, the method based on the interface `People` can be executed. 

And during this process, `itab.Hash` would be helpful. Actually, there is a global itab cache named `itabTable`, which holds the `itab`s. And the key of it, is `itab.Hash`

There is one more question need to pay attention. That's the offset of the type name. For example, when we check the info about `typ` parameter, we do like this:

```
(dlv) p firstmoduledata.types + 22497
4340389857
(dlv) x -fmt hex -count 1 -size 1 4340389858
0x102b517e2:   0x1b
(dlv) x -fmt raw -count 27 -size 1 4340389859
```

We get the address by offset, and check the `index 1th` for the `length` of the name, and get the real name from `index 2nd` for the real name.  

That's how the runtime handle the `Type.Str`. Refer to the code

```
type rtype struct {
	*abi.Type // embedding is okay here (unlike reflect) because none of this is public
}

func (t rtype) string() string {
	s := t.nameOff(t.Str).Name()
	if t.TFlag&abi.TFlagExtraStar != 0 {
		return s[1:]
	}
	return s
}

// Name returns the tag string for n, or empty if there is none.
func (n Name) Name() string {
	if n.Bytes == nil {
		return ""
	}
	i, l := n.ReadVarint(1)
	return unsafe.String(n.DataChecked(1+i, "non-empty string"), l)
}
```

The magic lies in line `i, l := n.ReadVarint(1)`. In our case, `i` would be 1, and `l` is the really length of the type name. For `Teacher` and `Student`, is 27 and `People`, 26.

In short, the layout would be like this:

```
[0] byte: flags
[1] byte: length
[2:]     : UTF-8 string of that length
```


### References

1. Golang Source Code V1.23
