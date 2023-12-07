# Slice

### Under Layer of Slice


**Must Read**

* https://go.dev/blog/slices-intro
* In short, it's like:

  ![slice internal](../images/slice-internal.png)

**Key Points**  

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


**Two Addresses**

Let's use this piece code to understand these 2 addresses

```
func main() {   
    a := make([]string, 0, 2)
	fmt.Printf("a: %p, %p\n", a, &a)
	a = append(a, "a1")
	fmt.Printf("a: %p, %p\n", a, &a)

	b := make([]string, 2, 2)
	b[0] = "b1"
	fmt.Printf("b: %p, %p\n", b, &b)
	b = append(b, "b2")
	fmt.Printf("b: %p, %p\n", b, &b)
}
```

In this code snipest, `fmt.Println("%p, %p\n", a, &a)`, would print the addresses of the `array under layer` and `the slice header`

Let's use the `dlv` tool to verify this.

* Launch the main.go with `dlv` tool. You may do this by:

```
dlv debug main.go
b main.main // set breakpoint at entry point
```
* And then, check the addresses of `a` 

```
> github.com/brofu/deepingo/slice.Addresses() ./slice/examples.go:10 (PC: 0x10aea9f)
     5: )
     6:
     7: func Addresses() {
     8:
     9:         a := make([]string, 0, 2)
=>  10:         fmt.Printf("a: %p, %p\n", a, &a)
    11:         a = append(a, "a1")
    12:         fmt.Printf("a: %p, %p\n", a, &a)
    13:
    14:         b := make([]string, 2, 2)
    15:         b[0] = "b1"
(dlv) n
a: 0xc0000be000, 0xc0000ac018   // NOTE: here, we got 2 addresses
> github.com/brofu/deepingo/slice.Addresses() ./slice/examples.go:11 (PC: 0x10aeb8c)
     6:
     7: func Addresses() {
     8:
     9:         a := make([]string, 0, 2)
    10:         fmt.Printf("a: %p, %p\n", a, &a)
=>  11:         a = append(a, "a1")
    12:         fmt.Printf("a: %p, %p\n", a, &a)
    13:
    14:         b := make([]string, 2, 2)
    15:         b[0] = "b1"
    16:         fmt.Printf("b: %p, %p\n", b, &b)
(dlv) p *(*string)(0xc0000be000)       //NOTE: we print the 1st evelement in the array, and we don't have any element yet
""
(dlv) p *(*unsafeheader.Slice)(0xc0000ac018)
internal/unsafeheader.Slice {Data: unsafe.Pointer(0xc0000be000), Len: 0, Cap: 2} //NOTE: we print the `SliceHeader`, by the address of a 
(dlv) n
> github.com/brofu/deepingo/slice.Addresses() ./slice/examples.go:12 (PC: 0x10aec3d)
     7: func Addresses() {
     8:
     9:         a := make([]string, 0, 2)
    10:         fmt.Printf("a: %p, %p\n", a, &a)
    11:         a = append(a, "a1")
=>  12:         fmt.Printf("a: %p, %p\n", a, &a)
    13:
    14:         b := make([]string, 2, 2)
    15:         b[0] = "b1"
    16:         fmt.Printf("b: %p, %p\n", b, &b)
    17:         b = append(b, "b2")
(dlv) n
a: 0xc0000be000, 0xc0000ac018   //NOTE: the address of the array under layer doens't change.
> github.com/brofu/deepingo/slice.Addresses() ./slice/examples.go:14 (PC: 0x10aed2c)
     9:         a := make([]string, 0, 2)
    10:         fmt.Printf("a: %p, %p\n", a, &a)
    11:         a = append(a, "a1")
    12:         fmt.Printf("a: %p, %p\n", a, &a)
    13:
=>  14:         b := make([]string, 2, 2)
    15:         b[0] = "b1"
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
> github.com/brofu/deepingo/slice.Addresses() ./slice/examples.go:15 (PC: 0x10aed87)
    10:         fmt.Printf("a: %p, %p\n", a, &a)
    11:         a = append(a, "a1")
    12:         fmt.Printf("a: %p, %p\n", a, &a)
    13:
    14:         b := make([]string, 2, 2)
=>  15:         b[0] = "b1"
    16:         fmt.Printf("b: %p, %p\n", b, &b)
    17:         b = append(b, "b2")
    18:         fmt.Printf("b: %p, %p\n", b, &b)
    19:
    20: }
(dlv) n
> github.com/brofu/deepingo/slice.Addresses() ./slice/examples.go:16 (PC: 0x10aedd2)
    11:         a = append(a, "a1")
    12:         fmt.Printf("a: %p, %p\n", a, &a)
    13:
    14:         b := make([]string, 2, 2)
    15:         b[0] = "b1"
=>  16:         fmt.Printf("b: %p, %p\n", b, &b)
    17:         b = append(b, "b2")
    18:         fmt.Printf("b: %p, %p\n", b, &b)
    19:
    20: }
(dlv) n
b: 0xc0000be020, 0xc0000ac060
> github.com/brofu/deepingo/slice.Addresses() ./slice/examples.go:17 (PC: 0x10aeece)
    12:         fmt.Printf("a: %p, %p\n", a, &a)
    13:
    14:         b := make([]string, 2, 2)
    15:         b[0] = "b1"
    16:         fmt.Printf("b: %p, %p\n", b, &b)
=>  17:         b = append(b, "b2")
    18:         fmt.Printf("b: %p, %p\n", b, &b)
    19:
    20: }
(dlv) p *(*string)(0xc0000be020)
"b1"  //NOTE: we have "b1"
(dlv) p *(*unsafeheader.Slice)(0xc0000ac060)
internal/unsafeheader.Slice {Data: unsafe.Pointer(0xc0000be020), Len: 2, Cap: 2} //NOTE: the `Slice Header` data

(dlv) n
> github.com/brofu/deepingo/slice.Addresses() ./slice/examples.go:18 (PC: 0x10aef7f)
    13:
    14:         b := make([]string, 2, 2)
    15:         b[0] = "b1"
    16:         fmt.Printf("b: %p, %p\n", b, &b)
    17:         b = append(b, "b2")
=>  18:         fmt.Printf("b: %p, %p\n", b, &b)
    19:
    20: }
(dlv) n
b: 0xc0000b8040, 0xc0000ac060  //NOTE: the address of the array under layer changed. Since there is capbility extension. This makes sense
> github.com/brofu/deepingo/slice.Addresses() ./slice/examples.go:20 (PC: 0x10af06c)
    15:         b[0] = "b1"
    16:         fmt.Printf("b: %p, %p\n", b, &b)
    17:         b = append(b, "b2")
    18:         fmt.Printf("b: %p, %p\n", b, &b)
    19:
=>  20: }

```

   
   
   
   
   
   
   
   
   
   
   
   
   
   
   
   
   
   
   
   
   
   
   
   
   
   
   
   
   
   
   
   
   
   
   
   
   
   
   
   
   
   
   
   
   
   
   
   
   
   
   
   
   
   
   
   
   
   
   
   
   
   
   
   
   
   
   
   
   
   
   
   
   
   
   
   
   
   
   
   
   
   
   
   
   
   
   
   
   
   

   




