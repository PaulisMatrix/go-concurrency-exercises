# General notes about Semaphores


1.  From OSTEP: Initializing the value of a semaphore: 

    ```
    When we set the value to 1(binary semaphore as a locking primitive), we can use the semaphore as a lock; 
    in the second, when to 0, we can use the semaphore for ordering(ex: parent, child ordering). 
    So what’s the general rule for semaphore initialization?
    
    One simple way to think about it, thanks to Perry Kivolowitz, is to consider the number of resources you are willing to give away immediately after initialization. With the lock, it was 1, because you are willing to have the lock locked (given away) immediately after initialization. 
    
    With the ordering case, it was 0, because there is nothing to give away at the
    start; only when the child thread is done is the resource created, at which
    point, the value is incremented to 1. Try this line of thinking on future
    semaphore problems, and see if it helps.
    ```

    ```    
    But in Golang, when we use a binary semaphore for mutual exclusion(sync.Mutext), the semantics are reversed. 
    Canonical definition being, init value 1=unlocked and 0=locked. 
    
    But in go, init value 0=unlocked and 1=locked, "mutexLocked = 1 << iota // mutex is locked" -> default 
    So when you do mutext.Lock() it does a CAS op and replaces value to 1 ergo locked and on mutext.UnLock(), 
    0 is swapped back, ergo unlocked. 
    
    This is opposite to the original semantics when you do sem_wait(Lock()), 1 value is decremented to 0 ergo locked
    and on sem_signal(UnLock()), 0 value is incremented to 1 ergo unlocked. This makes more sense cause when value is 0,
    no one can access the critical section cause then it would be negative, so you just keep spinning. 

    One possible reason for the initialisations in golang:
    The default state of a mutex should be unlocked when you do sync.Muxtex{}, you get an unlocked mutex.
    Generally we want to always have the 0 value be useful
    ```

2.  How semaphores can be used in golang: 

    
    *   semaphores can be used for both, for `locking/guarding a critical section (mutext.Lock() and mutext.Unlock())`
	    as well as for `ordering of process(sync.Waitgroup or chan)`. Both these pkgs use semaphores as synchornization primitive.
    

	*   Gotchas:

        sync.Mutex docs says that: 
        `A locked mutex is not associated with particular goroutine and it can be unlocked by an another goroutine as well`.

        Example: 
        <details>
        <summary>Click me</summary>
        
        ```golang
            func gotchas() {
                var wg sync.WaitGroup
                mutexChan := make(chan *sync.Mutex)

                wg.Add(1)
                go func() {
                    // create a lock and pass the lock in the chan
                    var mu sync.Mutex
                    fmt.Println("goroutine1 taking the lock")

                    mu.Lock()
                    mutexChan <- &mu

                    fmt.Println("goroutine2 will release the lock not held by it technically")
                    wg.Done()
                }()

                wg.Add(1)
                go func() {
                    mu := <-mutexChan
                    mu.Unlock()

                    fmt.Println("releasing the lock originally held by goroutine1")
                    wg.Done()
                }()

                wg.Wait()
            }
        ```    
        </details>
        
        This is correct since remember, `semaphores are just signally mechanism`.
        There is nothing something inherently stopping other goroutines to access the critical section. like goroutines being held back in shackles by something.<br><br>
	    semaphores just tell the current state(set, unset) of the critical section
	    more explanation here: https://www.reddit.com/r/golang/comments/1797dtu/comment/k54ckx2/<br><br>
	    But do remember that unlocking a mutex which doesn't even have a lock in the first place will obv not work
	    and also recursive locking will not work since its basically `hold and wait(already holding a lock and waiting to acquir a new one)` and that's one of the contenders for deadlocks.

3.  sync.Mutext source:

    **mutex** source:   https://go.dev/src/sync/mutex.go
    
    **semaphores** source:  https://go.dev/src/runtime/sema.go

    `which above mutex pkg uses to schedule/deschedule GOROUTINES, 
    similar to kernel doing it for THREADS using futex operation`

4.  References:
    
    *   More about semaphores:
    
        a. https://pages.cs.wisc.edu/~remzi/OSTEP/threads-sema.pdf

        b.  https://swtch.com/semaphore.pdf
    
    *   semaphore implementations(calling underlying C code or using syscalls): 

        a.  https://github.com/tmthrgd/go-sem

        b.  https://github.com/aka-mj/go-semaphore

        c.  https://github.com/shubhros/drunkendeluge/blob/master/semaphore/semaphore.go
        

5.  