package main

import (
	"fmt"
	"sync"
)

// https://leetcode.com/problems/building-h2o/

var wg sync.WaitGroup

func barrier(hchan, ochan <-chan string, doneHchan chan<- bool, doneOChan chan<- bool) {
	for {
		fmt.Print(<-hchan)
		fmt.Print(<-hchan)
		fmt.Print(<-ochan)

		doneHchan <- true
		doneOChan <- true
		fmt.Println()
	}
}

func buildHydrogen(hchan chan<- string, doneHchan <-chan bool) {
	defer wg.Done()

	// buffered hchan. gets blocked when its full otherwise keeps pushing.
	// but we need to wait for O to be pushed after every O
	hchan <- "H"
	hchan <- "H"
	<-doneHchan

}
func buildOxygen(ochan chan<- string, doneOChan <-chan bool) {
	defer wg.Done()
	ochan <- "O"
	<-doneOChan

}

func buildH2O() {
	hchan := make(chan string, 2)
	ochan := make(chan string)
	doneHChan := make(chan bool)
	doneOChan := make(chan bool)

	go barrier(hchan, ochan, doneHChan, doneOChan)

	for i := 0; i < 5; i++ {
		wg.Add(2)

		go buildHydrogen(hchan, doneHChan)
		go buildOxygen(ochan, doneOChan)
	}

	wg.Wait()
}
