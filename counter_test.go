package main

import (
	"fmt"
	"sync"
	"testing"
)

func TestCountUp(t *testing.T) {

	c := New()

	if c.Read() != 0 {
		t.Error("Counter did not initialise correctly")
	}

	for j := 0; j < 2; j++ {

		iterations := 1000

		for i := 0; i < iterations; i++ {
			c.Increment()
			if c.Read() != i+1 {
				t.Error("Counter did not increment correctly")
			}

		}
		if testing.Verbose() {
			fmt.Printf("Total count %d\n", c.Read())
		}

		c.Reset()
		if c.Read() != 0 {
			t.Error("Counter did not Reset correctly")
		}
	}

}

func TestCompetingWrites(t *testing.T) {

	c := New()

	iterations := 1000
	competingFuncs := 200

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

	if testing.Verbose() {
		fmt.Printf("Total count expected %d got %d\n", iterations*competingFuncs, c.Read())
	}

	if c.Read() != iterations*competingFuncs {
		t.Error("Locking failed, count was wrong")
	}

}
