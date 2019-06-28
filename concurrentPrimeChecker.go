package main

import (
	"fmt"
	"math"
	"runtime"
	"sync"
	"time"
)

func primeTesterFromTo(n int, from int, to int, isNotPrimeChannel chan <- bool, done <-chan interface{}) {
	defer close(isNotPrimeChannel)
	for i := from; i < to; i++ {
		select{
		case <-done:
			return
		default:
			if n%i == 0 {
				isNotPrimeChannel <- true
				return
			}
		}
	}
}

func createChannels(numCPUS int) []chan bool {
	isNotPrimeChannels := make([]chan bool, 0)
	for i := 0; i < numCPUS; i++ {
		isNotPrimeChannels = append(isNotPrimeChannels, make(chan bool))
	}
	return isNotPrimeChannels
}

func merge(cs ...chan bool) <-chan bool {
	out := make(chan bool)
	var wg sync.WaitGroup
	wg.Add(len(cs))
	for _, c := range cs {
		go func(c <-chan bool) {
			for v := range c {
				out <- v
			}
			wg.Done()
		}(c)
	}
	go func() {
		wg.Wait()
		close(out)
	}()
	return out
}

func primeTesterConc(n int) bool {
	done := make(chan interface{})
	defer close(done)

	numCPUs := runtime.NumCPU()
	numChannels := numCPUs
	isNotPrimeChannels := createChannels(numChannels)
	merged := merge(isNotPrimeChannels...)

	to_max := int(math.Sqrt(float64(n))) + 1
	to_interval := to_max / len(isNotPrimeChannels)
	for i := 0; i < len(isNotPrimeChannels); i++ {
		from := to_interval * i + 2
		to := to_interval * (i+1) + 2
		primeTesterFromTo(n, from, to, isNotPrimeChannels[i], done)
	}


	for v := range merged {
		if v == true {
			return false
		}
	}
	return true
}


func primeTesterConcTwo(n int) bool {
	done := make(chan interface{})
	defer close(done)

	numChannels := 2
	isNotPrimeChannels := createChannels(numChannels)

	to_max := int(math.Sqrt(float64(n))) + 1
	to_interval := to_max / len(isNotPrimeChannels)
	for i := 0; i < len(isNotPrimeChannels); i++ {
		from := to_interval * i + 2
		to := to_interval * (i+1) + 2
		primeTesterFromTo(n, from, to, isNotPrimeChannels[i], done)
	}


	for _, v := range isNotPrimeChannels {
		isNotPrime := <-v
		if isNotPrime == true {
			return false
		}
	}
	return true
}

func primeTesterLinear(n int) bool {
	for i := 2; i < int(math.Sqrt(float64(n))) + 1; i++ {
		if n%i == 0 {
			return false
		}
	}
	return true
}

func main() {
	n := 100000004987 // * 1002523// 1002523 * 1000099 // * 24941317
	start := time.Now()
	fmt.Println(primeTesterLinear(n))
	fmt.Println("Linear Search took:", time.Since(start))

	start = time.Now()
	fmt.Println(primeTesterConc(n))
	fmt.Println("Conc Search took:", time.Since(start))

	start = time.Now()
	fmt.Println(primeTesterConcTwo(n))
	fmt.Println("Conc Two Search took:", time.Since(start))
}
