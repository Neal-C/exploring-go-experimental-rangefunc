package main

import (
	"context"
	"fmt"
	"sync"
)

// Read the docs : https://go.dev/wiki/RangefuncExperiment
// don't mind the linting errors on the IDE => it does compile with the following command
// Requires Go 1.22.x
// GOEXPERIMENT=rangefunc go run .

func backwards(xs []int) func(func(i int, x int) bool) {
	return func(yield func(i int, x int) bool) {
		for i := len(xs) - 1; i >= 0; i-- {
			if !yield(i, xs[i]) {
				return
			}
		}
	}
}

func filter(array []int, filterFn func(item int) bool) func(func(index int, item int) bool) {
	return func(yield func(index int, item int) bool) {
		for index, item := range array {
			if !filterFn(item) {
				// if the iterator receives a 'break' statement, it will receive 'false'
				// and should no longer continue
				if !yield(index, item) {
					return
				}
			} 
		}
	}
}


func Parallel[E any](events []E) func(func (int, E) bool) {
	return func (yield func(int, E) bool ) {
		ctx, cancelFn := context.WithCancel(context.Background())
		defer cancelFn()

		var waitGroup sync.WaitGroup
		waitGroup.Add(len(events))

		for index, event := range events {
			go func(){
				defer waitGroup.Done()
				select {
				case <-ctx.Done():
					return
				default:
					if !yield(index, event){
						return
					}
				}
			}()
		}
	}
}

func isOdd(n int) bool {
	return n % 2 == 0
}

func main() {
	nums := []int{1, 2, 3, 4, 5}

	fmt.Println("backwards numbers")
	for _, x := range backwards(nums) {
		fmt.Println(x)
	}

	fmt.Println("filtered numbers")
	for _, x := range filter(nums, isOdd) {
		fmt.Println(x)
	}
}

