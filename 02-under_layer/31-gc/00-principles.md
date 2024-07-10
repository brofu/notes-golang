# GC Principles

### Workflows

* **Marking**. The marking phase identifies which objects are still reachable from the root set (global variables, stack variables, etc.). Objects that are still in use are marked as reachable.
* **Sweeping**. After marking, the sweeping phase begins. This phase reclaims the memory occupied by unreachable objects and returns it to the heap for future allocations.  
* **Pacing**. The GC dynamically adjusts its pace based on the application's allocation rate. This adaptive pacing ensures that garbage collection keeps up with the application's memory usage without introducing significant pauses


### Marking

####  Tri-color Marking Algorithm

**Roles**

Objects are marked with 3 colors, 
  * white (unreachable),  
  * gray (reachable but not fully processed), and 
  * black (reachable and fully processed). 

**Workflow**
1. Inital State 
  * Initially, all objects in the heap are conceptually **white**. 
  * The GC starts with a small number of `roots`, which are marked **grepy** and added to the grey set
  * `Roots` are the entry points to the object graph (**global variables, stack variables, and registers**).
2. Marking Phase
  * The algorithm processes each **gray** node by marking it **black** and processing its references. 
  * If a referenced node is **white**, it is marked **gray** and added to the gray set.
  * Repeat this process.
  * This process is running concurrently and incrementally.
3. Marking Completion
  * When the **gray** set is empty, all reachable nodes are **black**.
  * Any remaining **white** nodes are unreachable and can be reclaimed by the garbage collector.

**Features**
1. The `STW` time is very short. How?
  * **Concurrent Marking**.  The Go GC is concurrent, meaning it can run alongside the application (mutator). 
  * **Incremental Marking**. The marking work is broken down into small chunks, and the objects are being marked incrementally. 

**Questions**
1. How make sure all the objects are correctly marked? (Given that the marking progress is running concurrently with the application)
  * The GC uses a `write barrier to ensure that any changes to object references made by the application are correctly tracked during the marking phase.
  * Tri-Color Invariant
    * No **black** object points to a **white** object. This ensures that all reachable objects are eventually marked.
    * If a **gray** object points to a **white** object, the **white** object is also marked **gray**.
  
2. Why Golang GC can NOT remove the `STW` event totally?

  There are several scenarios the GC need to STW to make sure the process correctly.

  a. `Roots` Scan.
    * Roots of All goroutines need to scan, which includes `global variables, local variables and registers`. In a concurrent system like Go, where multiple goroutines are running, each goroutine’s stack and registers may contain references to heap objects. These are part of the root set that needs to be scanned.
    * Ensuring Consistency: To get a consistent snapshot of the roots, the GC must pause all goroutines at safe points. This ensures that no references are modified while the GC is scanning the roots, preventing inconsistencies.
    
  b. `Write Barriers` and `Tri-color Invariant`
    * During certain phases, the GC may need to apply these write barriers globally, requiring a brief pause to update references and maintain consistency.
    
  c. Concurrent Marking Completion
    * STW for another cycle of `Roots` scan.
    * Prepare for the Transition to Sweep Phase. (To make sure the marking is completely finished also. <span style="color:red">**TODO**</span>: More checking)

3. An example of STW event caused by `Write Barriers` and `Tri-color Invariant`?

  Let's check this code.
  ```
    package main

    import (
        "runtime"
        "time"
    )

    type Node struct {
        value int
        next  *Node
    }

    func main() {
        root := &Node{value: 1}
        root.next = &Node{value: 2}
        root.next.next = &Node{value: 3}

        // Start a goroutine that mutates the object graph
        go func() {
            for {
                root.next.next = &Node{value: 4}
                time.Sleep(100 * time.Millisecond)
            }
        }()

        // Trigger periodic garbage collections
        for i := 0; i < 10; i++ {
            runtime.GC() // This triggers a GC cycle with potential STW pauses
            time.Sleep(500 * time.Millisecond)
        }
    }
  ```
  * Before Mutation: root -> black, root.next -> black, root.next.next -> white 
  * During Mutation: The application changes root.next.next to point to a new node, which is white. The write barrier ensures that the new node is marked gray if it is white, maintaining the tri-color invariant.
  * After Mutation: The GC can then safely continue its marking phase.
  <span style="color:red">**TODO**</span>: More checking

4. In the question of 3. Is it OK to stop the goroutine which is updating the `root.next.next` ONLY? Why need to `STW`?
  <span style="color:red">**TODO**</span>: More checking




### Sweeping 
