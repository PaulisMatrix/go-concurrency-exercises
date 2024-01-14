//////////////////////////////////////////////////////////////////////
//
// Given is a producer-consumer scenario, where a producer reads in
// tweets from a mockstream and a consumer is processing the
// data. Your task is to change the code so that the producer as well
// as the consumer can run concurrently
//

package main

import (
	"fmt"
	"sync"
	"time"
)

func producer(stream Stream, tweetChan chan<- *Tweet, wg *sync.WaitGroup) {
	defer wg.Done()

	for {
		tweet, err := stream.Next()
		if err == ErrEOF {
			close(tweetChan)
			return
			//return tweets
		}
		tweetChan <- tweet
		//tweets = append(tweets, tweet)
	}
}

func consumer(tweetChan <-chan *Tweet, wg *sync.WaitGroup) {
	defer wg.Done()

	for {
		select {
		case t, ok := <-tweetChan:
			if !ok {
				fmt.Println("channel closed. no more tweets. returning!")
				return
			}
			if t.IsTalkingAboutGo() {
				fmt.Println(t.Username, "\ttweets about golang")
			} else {
				fmt.Println(t.Username, "\tdoes not tweet about golang")
			}

		}
	}

	/*
		for _, t := range tweets {
			if t.IsTalkingAboutGo() {
				fmt.Println(t.Username, "\ttweets about golang")
			} else {
				fmt.Println(t.Username, "\tdoes not tweet about golang")
			}
		}
	*/
}

func main() {
	var wg sync.WaitGroup

	start := time.Now()
	stream := GetMockStream()

	// use a chan for producer to add tweet to it and consumer to consume from

	tweetChan := make(chan *Tweet)

	// Producer
	wg.Add(1)
	go producer(stream, tweetChan, &wg)

	// Consumer
	wg.Add(1)
	go consumer(tweetChan, &wg)

	wg.Wait()
	fmt.Printf("Process took %s\n", time.Since(start))

}
