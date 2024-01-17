package main

import "fmt"

// https://leetcode.com/problems/print-zero-even-odd/description/

type ZeroEvenOdd struct {
	input chan int
	zero  chan int
	even  chan int
	odd   chan int

	done chan bool
	end  chan bool
}

func (z ZeroEvenOdd) zeroCall() {
	for {
		select {
		case <-z.end:
			break
		case item := <-z.zero:
			fmt.Print(item)
			z.done <- true
		}
	}

}

func (z ZeroEvenOdd) evenCall() {
	for {
		select {
		case <-z.end:
			break
		case item := <-z.even:
			if item&1 == 0 {
				fmt.Print(item)
				z.zero <- 0
			}
		}
	}
}

func (z ZeroEvenOdd) oddCall() {
	for {
		select {
		case <-z.end:
			break
		case item := <-z.odd:
			if item&1 != 0 {
				fmt.Print(item)
				z.zero <- 0
			}
		}
	}
}

func startZeroEvenOdding() {
	var n = 10

	z := ZeroEvenOdd{
		input: make(chan int),
		zero:  make(chan int),
		even:  make(chan int),
		odd:   make(chan int),
		end:   make(chan bool),
		done:  make(chan bool),
	}

	go z.zeroCall()
	go z.oddCall()
	go z.evenCall()

	go func(z ZeroEvenOdd) {
		for {
			select {
			case <-z.end:
				break
			case x := <-z.input:
				z.even <- x
				z.odd <- x
			}
		}
	}(z)

	z.zero <- 0
	for i := 1; i <= n; i++ {
		<-z.done
		z.input <- i
	}

	// wait for all done's to complete.
	<-z.done

	z.end <- true
}
