//////////////////////////////////////////////////////////////////////
//
// Given is a mock process which runs indefinitely and blocks the
// program. Right now the only way to stop the program is to send a
// SIGINT (Ctrl-C). Killing a process like that is not graceful, so we
// want to try to gracefully stop the process first.
//
// Change the program to do the following:
//   1. On SIGINT try to gracefully stop the process using
//          `proc.Stop()`
//   2. If SIGINT is called again, just kill the program (last resort)
//

package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	// Create a process
	proc := &MockProcess{}

	exit := make(chan os.Signal)
	done := make(chan bool)

	signal.Notify(exit, syscall.SIGINT)

	// run the process
	go func(proc *MockProcess) {
		proc.Run()
	}(proc)

	// wait for the signal
	go func(proc *MockProcess) {
		<-exit
		go proc.Stop()
		<-exit
		done <- true
	}(proc)

	<-done
	fmt.Println("shutting down!")

}
