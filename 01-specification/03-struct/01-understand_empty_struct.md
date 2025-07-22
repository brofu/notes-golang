## Empty Struct

### Usage

There are 3 typical usage of empty struct `struct{}`.

* **Used as a place holder.** 

```
// for example, implementation of HashSet

type Set map[int]struct{}

func NewSet() Set {
    return map[int]struct{}{}
}

func (s Set) Add(item int) {
    s[item] = struct{}{}
}

func (s Set) Remove(item int) {
    if _, ok := s[item]; ok {
        delete(s, item)
    }
}
```

* **Use as signal without data**

```
// Used to pass signal

func WorkControl() {
	ch := make(chan struct{}, 1)

	go func(chan struct{}) {
		// do some work
		fmt.Println("doing work...")
		fmt.Println("doing done")
		ch <- struct{}{}
	}(ch)

	select {
	case <-ch:
		fmt.Println("ok")
	}
}
```

* **Define simple interface implementations**

```
type InterfaceImplementation struct{}

func(i *InterfaceImplementation) Method() {
}
```

But why should we use it like this? Mainly reason is **struct{} uses NO memory space**


### zerobase

Before we talk about **zerobase**, let check this code piece first

```
func ExploreZerobase() {
	a := struct{}{}
	b := struct{}{}
	fmt.Println(a == b)
	fmt.Println(&a == &b)
	//fmt.Println(&a, &b)
	fmt.Println(&a == &b)
}
```
   



### References

1. https://mp.weixin.qq.com/s/KaAFRLKlWrefXQxRtliUjw


