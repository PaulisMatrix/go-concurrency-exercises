package main

import "fmt"

// https://leetcode.com/problems/fizz-buzz-multithreaded/

type FizzBuzz struct {
	input chan int
	done  chan bool
	end   chan bool

	fizz     chan int
	buzz     chan int
	fizzbuzz chan int
	num      chan int
}

func (fb FizzBuzz) Fizz() {
	for {
		select {
		case <-fb.end:
			break
		case item := <-fb.fizz:
			if (item % 3) == 0 {
				fmt.Println("fizz")
				fb.done <- true
			}
		}
	}
}

func (fb FizzBuzz) Buzz() {
	for {
		select {
		case <-fb.end:
			break
		case item := <-fb.buzz:
			if (item % 5) == 0 {
				fmt.Println("buzz")
				fb.done <- true
			}
		}

	}
}

func (fb FizzBuzz) FizzBuzz_() {
	for {
		select {
		case <-fb.end:
			break
		case item := <-fb.fizzbuzz:
			if item%3 == 0 && item%5 == 0 {
				fmt.Println("fizzbuzz")
				fb.done <- true
			}
		}

	}

}

func (fb FizzBuzz) Num() {
	for {
		select {
		case <-fb.end:
			break
		case item := <-fb.num:
			if item%3 != 0 && item%5 != 0 {
				fmt.Println(item)
				fb.done <- true
			}
		}
	}
}

func startFizzBuzzing() {
	var n = 15

	fb := FizzBuzz{
		input:    make(chan int),
		done:     make(chan bool),
		fizz:     make(chan int),
		buzz:     make(chan int),
		fizzbuzz: make(chan int),
		end:      make(chan bool),
		num:      make(chan int),
	}

	go fb.Fizz()
	go fb.Buzz()
	go fb.FizzBuzz_()
	go fb.Num()

	// send to all the channels
	go func(fb FizzBuzz) {
		for {
			select {
			case <-fb.end:
				break
			case x := <-fb.input:
				fb.fizz <- x
				fb.buzz <- x
				fb.fizzbuzz <- x
				fb.num <- x
			}
		}
	}(fb)

	for i := 1; i <= n; i++ {
		fb.input <- i
		<-fb.done
	}

	fb.end <- true
}
