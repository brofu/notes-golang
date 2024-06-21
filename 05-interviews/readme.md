
## Code Design

#### How to implement `singleton` in golang

1. Create instance at package importing

```
package singleton   

type singletonObject struct {}
var singletonObj = &singletonObject{}

func GetSingletonObject () *singletonObject {
    return singletonObj
}

```


* Instance created when the package is imported. Alternative way is to put the code in `init` function
* May lead to slow start of the program, as well as some memory waist 

2. With `lock`

```
package singleton   

import sync.Mutex

type singletonObject struct {}
var mutex sync.Mutext
var singletonObj *singletonObject  

func GetSingletonObject () *singletonObject {
    mutex.Lock()
    defer mutex.Unlock()
    if singletonObj == nil {
        singletonObj = &singletonObject{}
    }
    return singletonObj 
}
```

* Concurrency is not good enough.

3. `Lock` with pre-check

```
package singleton   

import sync.Mutex

type singletonObject struct {}
var mutex sync.Mutext
var singletonObj *singletonObject  

func GetSingletonObject () *singletonObject {
    if singletonObj == nil {
        mutex.Lock()
        defer mutex.Unlock()
        if singletonObj == nil {
            singletonObj = &singletonObject{}
        }
    }
    return singletonObj 
}
```

* Only if the instance is NOT initialized, race the lock
* Can reduce the load after the instance is initialized.

4. With `sync.Once`

```
var once sync.Once
func GetSingletonObject() *singletonObject {
    once.Do(func() {
        instance = &singletonObject{}
    })
    return instance
}

```

#### Follow Up. What's the implementation of `sync.Once`





## References:

1. https://juejin.cn/post/7160327827131203592
2. 
