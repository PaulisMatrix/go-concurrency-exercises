package main

import (
	"fmt"
	"sync"
)

// https://leetcode.com/problems/print-foobar-alternately/description/

func fooCall(n int, foo, bar chan bool, wg *sync.WaitGroup) {
	defer wg.Done()
	for i := 0; i < n; i++ {
		fmt.Print("foo")
		<-foo
		bar <- true
	}

}

func barCall(n int, bar, foo chan bool, wg *sync.WaitGroup) {
	defer wg.Done()
	for i := 0; i < n; i++ {
		foo <- true
		fmt.Print("bar")
		<-bar
	}

}

func alternateFooBar() {
	var wg sync.WaitGroup

	var n = 8
	var foo = make(chan bool)
	var bar = make(chan bool)

	wg.Add(1)
	go fooCall(n, foo, bar, &wg)

	wg.Add(1)
	go barCall(n, bar, foo, &wg)

	wg.Wait()
	close(foo)
	close(bar)

}
