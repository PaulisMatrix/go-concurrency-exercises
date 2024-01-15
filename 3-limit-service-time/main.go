//////////////////////////////////////////////////////////////////////
//
// Your video processing service has a freemium model. Everyone has 10
// sec of free processing time on your service. After that, the
// service will kill your process, unless you are a paid premium user.
//
// Beginner Level: 10s max per request
// Advanced Level: 10s max per user (accumulated)
//

package main

import (
	"sync/atomic"
	"time"
)

const MAX_SECONDS = 10

// User defines the UserModel. Use this to check whether a User is a
// Premium user or not
type User struct {
	ID        int
	IsPremium bool
	TimeUsed  int64 // in seconds
}

func (u *User) AddTime(seconds int64) int64 {
	return atomic.AddInt64(&u.TimeUsed, seconds)
}

// HandleRequest runs the processes requested by users. Returns false
// if process had to be killed
func HandleRequest(process func(), u *User) bool {
	// before processing the request, check if the TimeUsed by an non premium UserID > 10secs?
	completed := make(chan bool)

	// skip throttling is user is premium
	if u.IsPremium {
		process()
		return true
	}

	if atomic.LoadInt64(&u.TimeUsed) >= MAX_SECONDS {
		return false
	}

	// use ticker to tick for every sec and check if >= 10secs
	throttle := time.Tick(1 * time.Second)

	go func() {
		process()
		completed <- true
	}()

	for {
		select {
		case <-completed:
			return true
		case <-throttle:
			if u.AddTime(1) >= MAX_SECONDS {
				return false
			}
		}
	}

	/*
		funName := strings.Split(runtime.FuncForPC(reflect.ValueOf(process).Pointer()).Name(), ".")[1]
		if funName == "shortProcess" {
			u.AddTime(int64(6 * time.Second))
		} else if funName == "longProcess" {
			u.AddTime(int64(11 * time.Second))
		}

		process()
		return true
	*/
}

func main() {
	RunMockServer()
}
