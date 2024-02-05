package main

import (
	"math/rand"
	"sync"
	"time"
)

const MAX_SIZE = 1000

func MergeSortConcurrent(mylist []int) []int {
	// divide and conquer algorithm

	if len(mylist) > 1 {
		mid := int((len(mylist)) / 2)

		left_array := make([]int, len(mylist[:mid]))
		right_array := make([]int, len(mylist[mid:]))

		var wg sync.WaitGroup

		// wait for these merging to end
		// each one has to wait for its child routine to finish merging
		wg.Add(1)
		go func() {
			defer wg.Done()
			left_array = MergeSortConcurrent(mylist[:mid])
		}()

		wg.Add(1)
		go func() {
			defer wg.Done()
			right_array = MergeSortConcurrent(mylist[mid:])
		}()

		wg.Wait()

		// start merging
		var i = 0
		var j = 0
		var k = 0

		result := make([]int, len(left_array)+len(right_array))

		for i < len(left_array) && j < len(right_array) {
			if left_array[i] <= right_array[j] {
				result[k] = left_array[i]
				i++
			} else {
				result[k] = right_array[j]
				j++
			}
			k++

		}

		for i < len(left_array) {
			result[k] = left_array[i]
			k++
			i++
		}

		for j < len(right_array) {
			result[k] = right_array[j]
			k++
			j++
		}

		return result
	}
	return mylist
}

func MergeSort(mylist []int) []int {
	// divide and conquer algorithm

	if len(mylist) > 1 {
		mid := int((len(mylist)) / 2)

		left_array := MergeSort(mylist[:mid])
		right_array := MergeSort(mylist[mid:])

		// start merging
		var i = 0
		var j = 0
		var k = 0

		result := make([]int, len(left_array)+len(right_array))

		for i < len(left_array) && j < len(right_array) {
			if left_array[i] <= right_array[j] {
				result[k] = left_array[i]
				i++
			} else {
				result[k] = right_array[j]
				j++
			}
			k++

		}

		for i < len(left_array) {
			result[k] = left_array[i]
			k++
			i++
		}

		for j < len(right_array) {
			result[k] = right_array[j]
			k++
			j++
		}

		return result
	}
	return mylist
}

func mergeSort() {
	rand.Seed(time.Now().UnixNano())

	arr := make([]int, MAX_SIZE)

	for i := 0; i < MAX_SIZE; i++ {
		arr[i] = rand.Intn(1000000)
	}

	//MergeSort(arr)
	MergeSortConcurrent(arr)

}
