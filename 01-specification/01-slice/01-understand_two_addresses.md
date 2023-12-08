# Understand Two Addresses

#### Must Read

* https://go.dev/blog/slices-intro
* https://go.dev/blog/slices
* In short, it's like:

  ![slice internal](../../images/slice-internal.png)

#### Key Points  

So, it's easy to get the idea that, there are 2 addresses we may need to pay attention to
* The `address of the array under layer ` and
* The `address of the slice header`

The 1st address is straight forward, but what's the `slice header` mean? Actually, the `slice header` is the struct of slice itself above mentioned.

```
// SliceHeader is the runtime representation of a slice.
// It cannot be used safely or portably and its representation may
// change in a later release.
// Moreover, the Data field is not sufficient to guarantee the data
// it references will not be garbage collected, so programs must keep
// a separate, correctly typed pointer to the underlying data.
type SliceHeader struct {
	Data uintptr
	Len  int
	Cap  int
}

```

Ok, let's understand these 2 addresses with `dlv` debugger.


#### Two Addresses - string slice

Let's use this piece code to understand these 2 addresses. First, let's look at `string slice`

```
func StringSliceAddresses() {

	a := make([]string, 0, 2)
	fmt.Printf("a: %p, %p\n", a, &a)
	a = append(a, "a1")
	fmt.Printf("a: %p, %p\n", a, &a)

	b := make([]string, 2, 2)
	b[0] = "b0"
	fmt.Printf("b: %p, %p\n", b, &b)
	b = append(b, "b2")
	fmt.Printf("b: %p, %p\n", b, &b)

}

```

In this code gist, `fmt.Println("%p, %p\n", a, &a)`, would print the addresses of the `array under layer` and `the slice header`

Let's use the `dlv` tool to verify this.

* Launch the main.go with `dlv` tool. You may do this by:

```
dlv debug main.go
b main.main // set breakpoint at entry point
```
* And then, let's exam the memory of `a` 

```
> github.com/brofu/deepingo/slice.StringSliceAddresses() ./slice/examples.go:10 (PC: 0x10aea9f)
     5: )
     6:
     7: func StringSliceAddresses() {
     8:
     9:         a := make([]string, 0, 2)
=>  10:         fmt.Printf("a: %p, %p\n", a, &a)
    11:         a = append(a, "a1")
    12:         fmt.Printf("a: %p, %p\n", a, &a)
    13:
    14:         b := make([]string, 2, 2)
    15:         b[0] = "b0"
(dlv) n
a: 0xc0000be000, 0xc0000ac018   // NOTE: here, we got 2 addresses
> github.com/brofu/deepingo/slice.StringSliceAddresses() ./slice/examples.go:11 (PC: 0x10aeb8c)
     6:
     7: func StringSliceAddresses() {
     8:
     9:         a := make([]string, 0, 2)
    10:         fmt.Printf("a: %p, %p\n", a, &a)
=>  11:         a = append(a, "a1")
    12:         fmt.Printf("a: %p, %p\n", a, &a)
    13:
    14:         b := make([]string, 2, 2)
    15:         b[0] = "b0"
    16:         fmt.Printf("b: %p, %p\n", b, &b)
(dlv) p *(*string)(0xc0000be000)       //NOTE: we print the 1st evelement in the array, and we don't have any element yet
""
(dlv) p *(*unsafeheader.Slice)(0xc0000ac018)
internal/unsafeheader.Slice {Data: unsafe.Pointer(0xc0000be000), Len: 0, Cap: 2} //NOTE: we print the `SliceHeader`, by the address of a 
(dlv) n
> github.com/brofu/deepingo/slice.StringSliceAddresses() ./slice/examples.go:12 (PC: 0x10aec3d)
     7: func StringSliceAddresses() {
     8:
     9:         a := make([]string, 0, 2)
    10:         fmt.Printf("a: %p, %p\n", a, &a)
    11:         a = append(a, "a1")
=>  12:         fmt.Printf("a: %p, %p\n", a, &a)
    13:
    14:         b := make([]string, 2, 2)
    15:         b[0] = "b0"
    16:         fmt.Printf("b: %p, %p\n", b, &b)
    17:         b = append(b, "b2")
(dlv) n
a: 0xc0000be000, 0xc0000ac018   //NOTE: the address of the array under layer doens't change.
> github.com/brofu/deepingo/slice.StringSliceAddresses() ./slice/examples.go:14 (PC: 0x10aed2c)
     9:         a := make([]string, 0, 2)
    10:         fmt.Printf("a: %p, %p\n", a, &a)
    11:         a = append(a, "a1")
    12:         fmt.Printf("a: %p, %p\n", a, &a)
    13:
=>  14:         b := make([]string, 2, 2)
    15:         b[0] = "b0"
    16:         fmt.Printf("b: %p, %p\n", b, &b)
    17:         b = append(b, "b2")
    18:         fmt.Printf("b: %p, %p\n", b, &b)
    19:
(dlv) p *(*string)(0xc0000be000)
"a1"   //NOTE: now, we have the 1st element "a1"
(dlv) p *(*unsafeheader.Slice)(0xc0000ac018)
internal/unsafeheader.Slice {Data: unsafe.Pointer(0xc0000be000), Len: 1, Cap: 2}  //NOTE: the length of the array under layer changed.
```

* Now, let's continue with `b`

```
(dlv) n
> github.com/brofu/deepingo/slice.StringSliceAddresses() ./slice/examples.go:15 (PC: 0x10aed87)
    10:         fmt.Printf("a: %p, %p\n", a, &a)
    11:         a = append(a, "a1")
    12:         fmt.Printf("a: %p, %p\n", a, &a)
    13:
    14:         b := make([]string, 2, 2)
=>  15:         b[0] = "b0"
    16:         fmt.Printf("b: %p, %p\n", b, &b)
    17:         b = append(b, "b2")
    18:         fmt.Printf("b: %p, %p\n", b, &b)
    19:
    20: }
(dlv) n
> github.com/brofu/deepingo/slice.StringSliceAddresses() ./slice/examples.go:16 (PC: 0x10aedd2)
    11:         a = append(a, "a1")
    12:         fmt.Printf("a: %p, %p\n", a, &a)
    13:
    14:         b := make([]string, 2, 2)
    15:         b[0] = "b0"
=>  16:         fmt.Printf("b: %p, %p\n", b, &b)
    17:         b = append(b, "b2")
    18:         fmt.Printf("b: %p, %p\n", b, &b)
    19:
    20: }
(dlv) n
b: 0xc0000be020, 0xc0000ac060
> github.com/brofu/deepingo/slice.StringSliceAddresses() ./slice/examples.go:17 (PC: 0x10aeece)
    12:         fmt.Printf("a: %p, %p\n", a, &a)
    13:
    14:         b := make([]string, 2, 2)
    15:         b[0] = "b0"
    16:         fmt.Printf("b: %p, %p\n", b, &b)
=>  17:         b = append(b, "b2")
    18:         fmt.Printf("b: %p, %p\n", b, &b)
    19:
    20: }
(dlv) p *(*string)(0xc0000be020)
"b0"  //NOTE: we have "b0"
(dlv) p *(*unsafeheader.Slice)(0xc0000ac060)
internal/unsafeheader.Slice {Data: unsafe.Pointer(0xc0000be020), Len: 2, Cap: 2} //NOTE: the `Slice Header` data

(dlv) n
> github.com/brofu/deepingo/slice.StringSliceAddresses() ./slice/examples.go:18 (PC: 0x10aef7f)
    13:
    14:         b := make([]string, 2, 2)
    15:         b[0] = "b0"
    16:         fmt.Printf("b: %p, %p\n", b, &b)
    17:         b = append(b, "b2")
=>  18:         fmt.Printf("b: %p, %p\n", b, &b)
    19:
    20: }
(dlv) n
b: 0xc0000b8040, 0xc0000ac060  //NOTE: the address of the array under layer changed. Since there is capbility extension. This makes sense
> github.com/brofu/deepingo/slice.StringSliceAddresses() ./slice/examples.go:20 (PC: 0x10af06c)
    15:         b[0] = "b0"
    16:         fmt.Printf("b: %p, %p\n", b, &b)
    17:         b = append(b, "b2")
    18:         fmt.Printf("b: %p, %p\n", b, &b)
    19:
=>  20: }
(dlv) p *(*string)(0xc0000b8040) //NOTE: The 1st element of new array
"b0"
(dlv) p *(*string)(0xc0000b8040+16) //NOTE: The 2nd one
""
(dlv) p *(*string)(0xc0000b8040+32) //NOTE: The 3rd one. Got "b2" now
"b2"

```

There is a question here. **When get the 2nd and 3rd element (from the under layer array of a string slice), why the offset is 16 and 32? Considering that, the length of pointer (to string) on a 64bit system should be 8?** 

Let's discuss this in next section.

#### Two Addresses - string

The answer to above mentioned question is, the data stored in the under layer array, **is the `String Header` itself, instead of the pointer to the string (presented by a `String Header`)** 

How to understand this?
1. `string` in golang is just a **read-only slice of bytes**
2. It's presented by `reflect.StringHeader`, which contains a pointer to the data and length, each takes 8 bytes on 64 bit system
3. In the under layer array of string slice, the `StringHeader` itself is stored. (Instead of a pointer pointing to it)

The `refelct.StringHeader`: 
```
// StringHeader is the runtime representation of a string.
// It cannot be used safely or portably and its representation may
// change in a later release.
// Moreover, the Data field is not sufficient to guarantee the data
// it references will not be garbage collected, so programs must keep
// a separate, correctly typed pointer to the underlying data.
type StringHeader struct {
	Data uintptr
	Len  int
}
```

So actually **there are 2 addresses for string** also. And the hidden address is pointing to the under layer read-only bytes array

```
(dlv) n
> github.com/brofu/deepingo/slice.StringAddresses() ./slice/examples.go:25 (PC: 0x10aea7a)
    20:
    21: }
    22:
    23: func StringAddresses() {
    24:         a := "test"
=>  25:         fmt.Printf("a: %p\n", &a)
    26: }
    27:
    28: func IntegerSliceAddresses() {
    29:
    30:         a := make([]int, 2, 2)                                                                                                                 (dlv) n                                                                                                                                                a: 0xc000010250
> github.com/brofu/deepingo/slice.StringAddresses() ./slice/examples.go:26 (PC: 0x10aeae5)
    21: }                                                                                                                                                  22:                                                                                                                                                    23: func StringAddresses() {                                                                                                                           24:         a := "test"
    25:         fmt.Printf("a: %p\n", &a)
=>  26: }
    27:
    28: func IntegerSliceAddresses() {
    29:
    30:         a := make([]int, 2, 2)
    31:         fmt.Printf("a: %p, %p\n", a, &a)
(dlv) p *(*unsafeheader.String)(0xc000010250)
internal/unsafeheader.String {Data: unsafe.Pointer(0x10c67b7), Len: 4}  //NOTE: the address of bytes array is 0x10c67b7
```

#### Two Addresses - integer slice

Similar mechanism for integer slice. But the difference from string slice is that, the data stored in the under layer array is `*int`, so elements should be accessed by offset 8 * n (on 64 bit system). 

Debug info:

```
> github.com/brofu/deepingo/slice.IntegerSliceAddresses() ./slice/examples.go:35 (PC: 0x10aed17)
    30:         a := make([]int, 2, 2)
    31:         fmt.Printf("a: %p, %p\n", a, &a)
    32:         a[0] = 1
    33:         fmt.Printf("a: %p, %p\n", a, &a)
    34:         a = append(a, 3)
=>  35:         fmt.Printf("a: %p, %p\n", a, &a)
    36:
    37: }
(dlv) n
a: 0xc00001a100, 0xc00000c030   // NOTE: 2 addresses
> github.com/brofu/deepingo/slice.IntegerSliceAddresses() ./slice/examples.go:37 (PC: 0x10aee09)
    32:         a[0] = 1
    33:         fmt.Printf("a: %p, %p\n", a, &a)
    34:         a = append(a, 3)
    35:         fmt.Printf("a: %p, %p\n", a, &a)
    36:
=>  37: }
(dlv) p *(*unsafeheader.Slice)(0xc00000c030)
internal/unsafeheader.Slice {Data: unsafe.Pointer(0xc00001a100), Len: 3, Cap: 4}
(dlv) p *(*int)(0xc00001a100)
1
(dlv) p *(*int)(0xc00001a100+8)
0
(dlv) p *(*int)(0xc00001a100+16)  //NOTE: access 3rd element by offset 8 * 2
3
```

#### Two Addresses - struct slice

Similar machanism  for struct slice. 
* `[]StructType`, the data stored in under layer array, it's the value (an instance) of `StructType`. But
* `[]*StructType`, the data stored is the `pointer` pointing to the instance of `StructType`

Code for debug

```
type Test struct {
	A string
	B int
}

func StructSliceAddress() {

	a := make([]Test, 2, 2)
	b := make([]*Test, 2, 2)

	s1 := Test{
		A: "s1",
		B: 1,
	}

	s2 := Test{
		A: "s2",
		B: 2,
	}

	a[0], a[1] = s1, s2
	b[0], b[1] = &s1, &s2

	fmt.Printf("%p, %p\n", a, &a)
	fmt.Printf("%p, %p\n", b, &b)
}
```
Debug info

```
> github.com/brofu/deepingo/slice.StructSliceAddress() ./slice/examples.go:72 (PC: 0x10aed71)
    67:         }
    68:
    69:         a[0], a[1] = s1, s2
    70:         b[0], b[1] = &s1, &s2
    71:
=>  72:         fmt.Printf("%p, %p\n", a, &a)
    73:         fmt.Printf("%p, %p\n", b, &b)
    74: }
(dlv) n
0xc00007c180, 0xc00000c030
> github.com/brofu/deepingo/slice.StructSliceAddress() ./slice/examples.go:73 (PC: 0x10aee57)
    68:
    69:         a[0], a[1] = s1, s2
    70:         b[0], b[1] = &s1, &s2
    71:
    72:         fmt.Printf("%p, %p\n", a, &a)
=>  73:         fmt.Printf("%p, %p\n", b, &b)
    74: }
(dlv) p *(*unsafeheader.Slice)(0xc00000c030)
internal/unsafeheader.Slice {Data: unsafe.Pointer(0xc00007c180), Len: 2, Cap: 2}
(dlv) p *(*slice.Test)(0xc00007c180)
github.com/brofu/deepingo/slice.Test {A: "s1", B: 1}
(dlv) p *(*slice.Test)(0xc00007c180+24)  //NOTE, offset as 24. 16 for string header + 8 for integer (or pointer)
github.com/brofu/deepingo/slice.Test {A: "s2", B: 2}

(dlv) l
> github.com/brofu/deepingo/slice.StructSliceAddress() ./slice/examples.go:73 (PC: 0x10aee57)
    68:
    69:         a[0], a[1] = s1, s2
    70:         b[0], b[1] = &s1, &s2
    71:
    72:         fmt.Printf("%p, %p\n", a, &a)
=>  73:         fmt.Printf("%p, %p\n", b, &b)
    74: }
(dlv) n
0xc000010250, 0xc00000c048
> github.com/brofu/deepingo/slice.StructSliceAddress() ./slice/examples.go:74 (PC: 0x10aef3e)
    69:         a[0], a[1] = s1, s2
    70:         b[0], b[1] = &s1, &s2
    71:
    72:         fmt.Printf("%p, %p\n", a, &a)
    73:         fmt.Printf("%p, %p\n", b, &b)
=>  74: }
(dlv) p *(*unsafe.Pointer)(0xc000010250)  //NOTE: Get pointer first
unsafe.Pointer(0xc00000c060)
(dlv) p *(*slice.Test)(0xc00000c060)  //NOTE: get struct from pointer
github.com/brofu/deepingo/slice.Test {A: "s1", B: 1}
(dlv) p *(*unsafe.Pointer)(0xc000010250+8)  //NOTE: offset as 8
unsafe.Pointer(0xc00000c078)
(dlv) p *(*slice.Test)(0xc00000c078)
github.com/brofu/deepingo/slice.Test {A: "s2", B: 2}
```

#### Two Addresses - map slice

Similar mechanism for map slice. The data stored in under layer array is pointer to the map.

```
(dlv) n
> github.com/brofu/deepingo/slice.MapSliceAddresses() ./slice/examples.go:45 (PC: 0x10af096)
    40:
    41:         a := make([]map[int]string, 2, 2)
    42:         fmt.Printf("a: %p, %p\n", a, &a)
    43:         a[0] = map[int]string{0: "0"}
    44:         a[1] = map[int]string{1: "1"}
=>  45:         fmt.Printf("a: %p, %p\n", a, &a)
    46:
    47: }
(dlv) n
a: 0xc000010230, 0xc00000c030  // NOTE: 2 addresses
> github.com/brofu/deepingo/slice.MapSliceAddresses() ./slice/examples.go:47 (PC: 0x10af185)
    42:         fmt.Printf("a: %p, %p\n", a, &a)
    43:         a[0] = map[int]string{0: "0"}
    44:         a[1] = map[int]string{1: "1"}
    45:         fmt.Printf("a: %p, %p\n", a, &a)
    46:
=>  47: }
(dlv) p *(*unsafeheader.Slice)(0xc00000c030)
internal/unsafeheader.Slice {Data: unsafe.Pointer(0xc000010230), Len: 2, Cap: 2}
(dlv) p *(*map[int]string)(0xc000010230)
map[int]string [
        0: "0",
]
(dlv) p *(*map[int]string)(0xc000010230+8)  //NOTE: offset as 8 * n
map[int]string [
        1: "1",
]

```
