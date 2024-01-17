package main

import (
	"fmt"
	"sync"
)

func main() {
	a := make(chan int, 2)
	b := make(chan int, 2)
	c := make(chan int, 1)

	a <- 5
	a <- 6
	b <- 7
	b <- 8
	c <- 9

	close(a)
	close(b)
	close(c)

	fmt.Println("Using wait group in a goroutine: ")
	for num := range mergeChannels(a, b, c) {
		fmt.Println(num)
	}

	d := make(chan int, 2)
	e := make(chan int, 2)
	f := make(chan int, 1)

	d <- 5
	d <- 6
	e <- 7
	e <- 8
	f <- 9

	close(d)
	close(e)
	close(f)

	fmt.Println("Using buffered channel:")
	for num := range mergeChannelsBuffered(5, d, e, f) {
		fmt.Println(num)
	}
}

func mergeChannels(n ...chan int) chan int {
	wg := &sync.WaitGroup{}
	result := make(chan int)

	for _, c := range n {
		wg.Add(1)
		go func(ch chan int) {
			for chanValue := range ch {
				result <- chanValue
			}
			wg.Done()
		}(c)
	}

	go func() {
		wg.Wait()
		close(result)
	}()

	return result
}

// using buffered channel to avoid blocking
func mergeChannelsBuffered(size int, n ...chan int) chan int {
	wg := &sync.WaitGroup{}
	result := make(chan int, size)

	for _, c := range n {
		wg.Add(1)
		go func(ch chan int) {
			for chanValue := range ch {
				result <- chanValue
			}
			wg.Done()
		}(c)
	}

	wg.Wait()
	close(result)

	return result
}
