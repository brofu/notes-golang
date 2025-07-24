package interface_examples

import "errors"

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
