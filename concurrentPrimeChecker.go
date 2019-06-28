package main

import (
	"fmt"
	"math"
	"runtime"
	"sync"
	"time"
)

func primeTesterFromTo_CheckDoneLessOften(n int, from int, to int, isNotPrimeChannel chan <- bool, done <-chan interface{}) {
	defer close(isNotPrimeChannel)
	interval := int((to - from) / 100) + 1
	for i := 0; i < 100; i++ {
		select {
		case <-done:
			return
		default:
		}
		for j := 0; j <= interval; j++ {
			if n % (from + interval * i + j) == 0 {
				isNotPrimeChannel <- true
				return
			}
		}
	}
}

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

func createChannels(numChannels int) []chan bool {
	isNotPrimeChannels := make([]chan bool, 0)
	for i := 0; i < numChannels; i++ {
		isNotPrimeChannels = append(isNotPrimeChannels, make(chan bool))
	}
	return isNotPrimeChannels
}

func merge(cs ...chan bool) <-chan bool {
	// https://medium.com/justforfunc/two-ways-of-merging-n-channels-in-go-43c0b57cd1de
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
		go primeTesterFromTo(n, from, to, isNotPrimeChannels[i], done)
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

	isNotPrimeChannels := createChannels(2)

	max := int(math.Sqrt(float64(n))) + 1
	mid := int(max / 2) + 1
	go primeTesterFromTo(n, 2, mid, isNotPrimeChannels[0], done)
	go primeTesterFromTo(n, mid, max, isNotPrimeChannels[1], done)

	if <-isNotPrimeChannels[0] {
		return false
	}
	return !<-isNotPrimeChannels[1]
}


func primeTesterConcTwo_CheckDoneLessOften(n int) bool {
	done := make(chan interface{})
	defer close(done)

	isNotPrimeChannels := createChannels(2)

	max := int(math.Sqrt(float64(n))) + 1
	mid := int(max / 2) + 1
	go primeTesterFromTo_CheckDoneLessOften(n, 2, mid, isNotPrimeChannels[0], done)
	go primeTesterFromTo_CheckDoneLessOften(n, mid, max, isNotPrimeChannels[1], done)

	if <-isNotPrimeChannels[0] {
		return false
	}
	return !<-isNotPrimeChannels[1]
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


	start = time.Now()
	fmt.Println(primeTesterConcTwo_CheckDoneLessOften(n))
	fmt.Println("Conc Two Search took:", time.Since(start))
}