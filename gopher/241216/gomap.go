package main

import (
	"fmt"
	"sync"
)

// 使用锁或者channel实现一个线程安全的map
func main() {
	func2()
}

func func1() {
	var (
		mutex   sync.Mutex
		mapList = make(map[int]struct{})
		wg      = &sync.WaitGroup{}
	)

	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func(wg *sync.WaitGroup, i int) {
			mutex.Lock()
			mapList[i] = struct{}{}
			mutex.Unlock()
			defer func() { wg.Done() }()
		}(wg, i)
	}
	wg.Wait()
	fmt.Println("num:", len(mapList))
}

func func2() {
	var (
		numChan  = make(chan int, 1000)
		numCount int
		wg       = &sync.WaitGroup{}
		endChan  = make(chan struct{}, 1)
	)
	defer func() {
		close(numChan)
		close(endChan)
	}()

	for i := 0; i < 1000; i++ {
		go func(wg *sync.WaitGroup) {
			numChan <- i
			wg.Add(1)
		}(wg)

	}

	go func(wg *sync.WaitGroup, numChan chan int, endChan chan struct{}) {
		for {
			select {
			case i, ok := <-numChan:
				if ok {
					numCount++
					fmt.Println("num:", numCount, i)
					wg.Done()
				}
			case <-endChan:
				return
			}
		}
	}(wg, numChan, endChan)

	wg.Wait()
	endChan <- struct{}{}
	fmt.Println("numCount:", numCount)
}
