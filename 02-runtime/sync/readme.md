
## sync.Once

#### Source Code

```
type Once struct {
	done uint32
	m    Mutex
}

func (o *Once) Do(f func()) {
	if atomic.LoadUint32(&o.done) == 0 { 
		o.doSlow(f)
	}
}

func (o *Once) doSlow(f func()) {
	o.m.Lock()
	defer o.m.Unlock()
	if o.done == 0 {
		defer atomic.StoreUint32(&o.done, 1)
		f()
	}
}
```

#### Key Points

1. A Once must not be copied after first use. Since `defer atomic.StoreUint32(&o.done, 1)` would update `o.done` by address
2. Why `done` is the first?	
  * `It is first in the struct because it is used in the hot path. The hot path is inlined at every call site. Placing done first allows more compact instructions on some architectures (amd64/386), and fewer instructions (to calculate offset) on other architectures.
  * In short, `o.done` is high frequently accessed, put it first would make the performance better. 
3. Pay attention to `Deadlock`. Because no call to Do returns until the one call to f returns, if f causes Do to be called, it will deadlock.
4. About `panic`. If f panics, Do considers it to have returned; future calls of Do return without calling f. => How? 
5. Wrong version of `Once.Do`. Why?
  ```
  if atomic.CompareAndSwapUint32(&o.done, 0, 1) {
    f()
  }
  ```
  > Do guarantees that when it returns, f has finished. This implementation would not implement that guarantee:
	given two simultaneous calls, the winner of the cas would call f, and the second would return immediately, without waiting for the first's call to f to complete.
	This is why the slow path falls back to a mutex, and why the atomic.StoreUint32 must be delayed until after f returns.

  In short, `CAS` would return at immediately, but lock would block the call.
