package main

import (
	"fmt"
	"sync"
)

type Counter struct {
	mux   sync.Mutex
	count int
}

func (c *Counter) Increment() {
	c.mux.Lock()
	c.count++
	c.mux.Unlock()
}

func (c *Counter) Read() int {
	c.mux.Lock()
	count := c.count
	c.mux.Unlock()
	return count
}

func (c *Counter) Reset() {
	c.mux.Lock()
	c.count = 0
	c.mux.Unlock()
}

func New() *Counter {
	return &Counter{count: 0}
}

func main() {

	iterations := 1000
	competingFuncs := 200

	nmExpected, nmGot := demoNonMux(iterations, competingFuncs)
	mExpected, mGot := demoMux(iterations, competingFuncs)
	fmt.Println("\u250C----------------------------------------\u2510")
	fmt.Println("| method | expected |    got   |  ok?    |")
	fmt.Println("|--------|----------|----------|---------|")
	fmt.Printf("|non-mux |%8d  |%8d  |  %v  |\n", nmExpected, nmGot, nmExpected == nmGot)
	fmt.Printf("|  mux   |%8d  |%8d  |  %v   |\n", mExpected, mGot, mExpected == mGot)
	fmt.Println("\u2514----------------------------------------\u2518")
}

func demoNonMux(iterations int, competingFuncs int) (int, int) {

	c := 0

	var wg sync.WaitGroup
	wg.Add(competingFuncs)

	for j := 0; j < competingFuncs; j++ {
		go func() {
			defer wg.Done()
			for i := 0; i < iterations; i++ {
				c++
			}
		}()
	}
	wg.Wait()

	return iterations * competingFuncs, c

}

func demoMux(iterations int, competingFuncs int) (int, int) {

	c := New()

	var wg sync.WaitGroup
	wg.Add(competingFuncs)

	for j := 0; j < competingFuncs; j++ {
		go func() {
			defer wg.Done()
			for i := 0; i < iterations; i++ {
				c.Increment()
			}
		}()
	}
	wg.Wait()

	return iterations * competingFuncs, c.Read()

}
