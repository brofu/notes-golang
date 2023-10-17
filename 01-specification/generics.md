# Generics

### Usage Scenarios 

#### Case 1

**The Situation**

Let's see one example 

```
// An interface Human
type Human interface {
    Name() string
}

// An Implementation to Human interface 
type Hero struct {
    name string
}
func (h *Hero) Name() string {
    return h.name
}

// A function would say hello to Human interface
func SayHello(humans []Human) {
    for _, human := range humans {
        fmt.Prinln("Hello, ", human.Name())
    }
}

// Call the function with Hero lists
func main() {
    heroList := []*Hero{&Hero{name: "thor"}, &Hero{name: "apolo"}}
    // This doesn't work
    // cannot use heroList (variable of type []*Hero) as []Human value in argument to SayHello
    SayHello(heroList) // This doesn't work
}
```

**The Reason**

> There is a general rule that syntax should not hide complext/costly operations

**Options**

One option for this it to use the `Generics` which is introduced from 1.18

```
// new version with Go Generics
func SayHelloV2[T Human](humans []T) {
	for _, human := range humans {
		fmt.Println("Hello, ", human.Name())
	}
}

// call the new function
SayHelloV2[*Hero](heroList)

```

**Pons and Cons of Go Generics** 

**Reference**

https://dusted.codes/using-go-generics-to-pass-struct-slices-for-interface-slices
https://go.dev/blog/intro-generics



