package main

// https://leetcode.com/problems/print-in-order/description/

import "fmt"

func firstCall(first, second chan string) {
	fmt.Print(<-first)
	second <- "second"
}

func secondCall(second, third chan string) {
	fmt.Print(<-second)
	third <- "third"
}

func thirdCall(third chan string) {
	fmt.Print(<-third)
	third <- "done"
}

func main() {
	nums := []int{1, 2, 3}
	first := make(chan string)
	second := make(chan string)
	third := make(chan string)

	for _, num := range nums {
		switch num {
		case 1:
			go firstCall(first, second)
		case 2:
			go secondCall(second, third)
		case 3:
			go thirdCall(third)
		}
	}

	first <- "first"
	<-third
}
