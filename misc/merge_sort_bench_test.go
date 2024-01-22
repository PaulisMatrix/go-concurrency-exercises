package main

import "testing"

var num = 1000

func BenchmarkMergeSort(b *testing.B) {
	for i := 0; i < num; i++ {
		mergeSort()
	}
}

func BenchmarkMergeSorConcurrent(b *testing.B) {
	for i := 0; i < num; i++ {
		mergeSort()
	}
}
