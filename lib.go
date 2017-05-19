package main

import (
	"time"
	"sync"
	"fmt"
	"log"
	"sync/atomic"
)

const max = 9999  // total number of data
const jump = 40   // square root of maximum returned LIMIT query, also this number*2 is the number of parallel SELECT
const source = 99 // number of parallel INSERT/UPDATE

// note: within for loop, goroutines will not executed until runtime.Gosched() called

func BenchmarkInsert(exec func(x int) error) {
	var wg sync.WaitGroup
	
	t := time.Now()
	wg.Add(source)
	for y := 0; y < source; y++ {
		go func(y int) {
			subops := 0
			t2 := time.Now()
			defer wg.Done()
			for x := 1 + max * y / source; x <= max * (y + 1) / source; x++ {
				subops++
				if err := exec(x); err != nil {
					log.Fatal(err)
					return
				}
			}
			dur2 := time.Now().Sub(t2)
			fmt.Printf("I-%02d: (%.2f ms/op: %d)\n", y, float64(dur2.Nanoseconds()) / 1000000 / float64(subops), subops)
		}(y)
	}
	wg.Wait()
	dur := time.Now().Sub(t)
	fmt.Printf("INSERT: %v (%.2f ms/op)\n", dur, float64(dur.Nanoseconds()) / 1000000 / max)
}

func BenchmarkUpdate(exec func(x int) error) {
	var wg sync.WaitGroup
	
	t := time.Now()
	wg.Add(source)
	for y := 0; y < source; y++ {
		go func(y int) {
			subops := 0
			t2 := time.Now()
			defer wg.Done()
			for x := 1 + max / source * y; x <= max / source * (y + 1); x++ {
				subops++
				if err := exec(x); err != nil {
					log.Fatal(err)
					return
				}
			}
			dur2 := time.Now().Sub(t2)
			fmt.Printf("U-%02d: (%.2f ms/op: %d)\n", y, float64(dur2.Nanoseconds()) / 1000000 / float64(subops), subops)
		}(y)
	}
	wg.Wait()
	dur := time.Now().Sub(t)
	fmt.Printf("UPDATE: %v (%.2f ms/op)\n", dur, float64(dur.Nanoseconds()) / 1000000 / max)
}

func BenchmarkSelect(exec1, exec2 func(x, lim int) error) {
	var wg sync.WaitGroup
	
	t := time.Now()
	ops := int64(0)
	for y := 2; y < jump; y++ {
		wg.Add(1)
		go func(y int) {
			defer wg.Done()
			subops := 0
			t2 := time.Now()
			for x := max - 1; x > 0; x -= y {
				atomic.AddInt64(&ops, 1)
				subops++
				if err := exec1(x,y*y); err != nil {
					log.Fatal(err)
					return
				}
			}
			dur2 := time.Now().Sub(t2)
			fmt.Printf("S-%02d: (%.2f ms/op: %d)\n", y, float64(dur2.Nanoseconds()) / 1000000 / float64(subops), subops)
		}(y)
		wg.Add(1)
		go func(y int) {
			defer wg.Done()
			subops := 0
			t2 := time.Now()
			for x := 1; x < max; x += y {
				atomic.AddInt64(&ops, 1)
				subops++
				if err := exec2(x,y*y); err != nil {
					log.Fatal(err)
					return
				}
			}
			dur2 := time.Now().Sub(t2)
			fmt.Printf("S-%02d: (%.2f ms/op: %d)\n", y + jump, float64(dur2.Nanoseconds()) / 1000000 / float64(subops), subops)
		}(y)
	}
	wg.Wait()
	dur := time.Now().Sub(t)
	fmt.Printf("SELECT: %v (%.2f ms/op: %d)\n", dur, float64(dur.Nanoseconds()) / 1000000 / float64(ops), ops)
}