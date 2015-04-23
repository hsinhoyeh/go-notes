package test

import (
	"sync"
	"testing"
)

// Q: should we call defer mu.Unlock before/after we call mu.Lock?
// check the benchmark result below:
//--- master* ?› » go test -bench=Defer
//PASS
//BenchmarkDefer1 2000000000               0.09 ns/op
//BenchmarkDefer2 2000000000               0.07 ns/op

func f1(mu *sync.Mutex, val *int) int {
	mu.Lock()
	defer mu.Unlock()

	(*val) = (*val) + 1
	return (*val)
}

func f2(mu *sync.Mutex, val *int) int {
	defer mu.Unlock()
	mu.Lock()

	(*val) = (*val) + 1
	return (*val)
}

func BenchmarkDefer1(b *testing.B) {
	run(f1, 5)
}

func BenchmarkDefer2(b *testing.B) {
	run(f2, 5)
}

func run(f func(mu *sync.Mutex, val *int) int, numRT int) {
	acc := 0
	goal := 1000000
	mu := &sync.Mutex{}
	wg := &sync.WaitGroup{}

	wg.Add(numRT)
	for i := 0; i < numRT; i++ {

		go func() {
			defer wg.Done()
			for {
				v := f(mu, &acc)
				if v > goal {
					return
				}
			}
		}()
	}
	wg.Wait()
}
