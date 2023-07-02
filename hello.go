package main

import (
	"context"
	"time"

	"fmt"
	"sync"
)

func main() {

	arr := make([]int, 1000000)
	for i := 0; i < 1000000; i++ {
		arr = append(arr, i)
	}

	start := time.Now()

	jobs := make(chan int, len(arr))
	res := make(chan int, len(arr))
	ctx := context.Background()
	ctxt, cancel := context.WithCancel(ctx)
	go distributor(jobs, res, cancel, ctxt)

	for _, v := range arr {
		jobs <- v
	}

Loop:
	for {
		select {
		case re := <-res:
			if re == 52345 {
				fmt.Println("Found it")
				cancel()
				break Loop
			}

		case <-ctxt.Done():
			fmt.Println("Not Found")
			break Loop

		}
	}

	elapsed := time.Since(start)
	fmt.Printf("Elapsed time: %s\n", elapsed)

}

func worker(jobs chan int, res chan int, wg *sync.WaitGroup, ctx context.Context) {
	defer wg.Done()
	for {
		select {
		case job, ok := <-jobs:
			if !ok {
				return
			}

			res <- job
		case <-ctx.Done():
			return

		}
	}

}

func distributor(jobs chan int, res chan int, cancel context.CancelFunc, ctx context.Context) {
	var wg sync.WaitGroup
	for i := 0; i < 4; i++ {
		wg.Add(1)
		go worker(jobs, res, &wg, ctx)

	}
	wg.Wait()
	close(res)
	cancel()

}
