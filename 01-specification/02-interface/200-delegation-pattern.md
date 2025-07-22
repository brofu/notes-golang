## Delegation Pattern

### Concepts

>`Delegation` means: "I'm not implementing this method myself, but instead I'm passing the call to another object that knows how."

In Go, delegation is commonly implemented through:

  * `Struct embedding` (automatic delegation)
  * Explicit forwarding (manual delegation)

But with `interface`, we may have better flexible and extensible

### An Example 

Let's say that,  we have an `set` of integers, we can add or remove number from it.

```
type IntSet struct {
	data map[int]struct{}
}

func NewIntSet() IntSet {
	return IntSet{make(map[int]struct{})}
}

func (set *IntSet) Add(x int) {
	set.data[x] = struct{}{}
}

func (set *IntSet) Delete(x int) {
	delete(set.data, x)
}
```

And now, we need an `Undo` feature to `IntSet`. That's to say, remove the added number with `Undo`, or add the removed number with it. We may do that like this

Define the `Undoer` interface and implement it 

```
type Undoer interface {
	Trace(func())
	Undo() error
}

type DefaultUndoer []func()

func (du *DefaultUndoer) Trace(f func()) {
	*du = append(*du, f)
}

func (du *DefaultUndoer) Undo() error {
	if len(*du) == 0 {
		return errors.New("no function traced")
	}

	f := (*du)[len(*du)-1]
	if f != nil {
		f()
	}

	*du = (*du)[:len(*du)-1]
	return nil
}
```

Utilize the new interface
```
type IntSetWithUndo struct {
	IntSet
	Undoer
}

func NewIntSetWithUndo() IntSetWithUndo {
	return IntSetWithUndo{
		NewIntSet(),
		new(DefaultUndoer),
	}
}

func (set *IntSetWithUndo) Add(x int) {
	set.IntSet.Add(x)
	set.Undoer.Trace(func() {
		set.IntSet.Delete(x)
	})
}

func (set *IntSetWithUndo) Delete(x int) {
	set.IntSet.Delete(x)
	set.Undoer.Trace(func() {
		set.IntSet.Add(x)
	})
}

func (set *IntSetWithUndo) Undo() error {
	return set.Undoer.Undo()
}
```

With this we can gain the benefits, 

1. No modification to the `IntSet`. Closed to the modification
2. Implement the feature with `IntSetWithUndo`. Open to extension.
3. `Undoer` can ALSO be used by other components, with `Delegation Pattern`.


### Reference

1. [昨耳听风 - 编程范式](https://time.geekbang.org/column/article/2748) 

