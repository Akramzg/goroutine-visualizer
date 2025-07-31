package main

import (
	"context"
	"fmt"
	"math"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

func crunchNumbers(n int) float64 {
	result := 0.0
	for i := 0; i < n; i++ {
		result += math.Sqrt(float64(i))
	}
	return result
}

func processWorker(id int, ch chan int, wg *sync.WaitGroup) {
	defer wg.Done()
	<-ch
}

func queryWorker(ctx context.Context) {
	select {
	case <-ctx.Done():
		return
	case <-time.After(10 * time.Minute):
	}
}

func dbWorker(ctx context.Context, id int, wg *sync.WaitGroup) {
	defer wg.Done()
	select {
	case <-ctx.Done():
		return
	case <-time.After(8 * time.Second):
	}
	crunchNumbers(100000)
}

func fanOutWorker(jobs <-chan int, results chan<- float64, wg *sync.WaitGroup) {
	defer wg.Done()
	for j := range jobs {
		results <- crunchNumbers(j * 1000)
	}
}

func main() {
	go func() {
		if err := http.ListenAndServe(":6060", nil); err != nil {
			fmt.Printf("pprof server error: %v\n", err)
		}
	}()

	fmt.Println("scratch app running — pprof on :6060")
	fmt.Println("Press Ctrl+C to gracefully shutdown.")

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

Loop:
	for {
		// --- Phase 1: Channel Waiters ---
		fmt.Println("[phase 1] spawning channel waiters...")
		var wg1 sync.WaitGroup
		ch := make(chan int)
		for i := 0; i < 15; i++ {
			wg1.Add(1)
			go processWorker(i, ch, &wg1)
		}

		select {
		case <-ctx.Done():
			close(ch)
			wg1.Wait()
			break Loop
		case <-time.After(3 * time.Second):
			close(ch) // Unblock phase 1 workers before moving to phase 2
			wg1.Wait()
		}

		// --- Phase 2: Select Waiters ---
		fmt.Println("[phase 2] spawning select waiters...")
		for i := 0; i < 20; i++ {
			go queryWorker(ctx)
		}

		select {
		case <-ctx.Done():
			break Loop
		case <-time.After(3 * time.Second):
		}

		// --- Phase 3: DB Workers ---
		fmt.Println("[phase 3] spawning db workers...")
		var wg2 sync.WaitGroup
		for i := 0; i < 8; i++ {
			wg2.Add(1)
			go dbWorker(ctx, i, &wg2)
		}

		select {
		case <-ctx.Done():
			wg2.Wait()
			break Loop
		case <-time.After(3 * time.Second):
			wg2.Wait()
		}

		// --- Phase 4: Fan-out Pipeline ---
		fmt.Println("[phase 4] fan-out pipeline...")
		jobs := make(chan int, 50)
		results := make(chan float64, 50)
		var wg3 sync.WaitGroup
		for i := 0; i < 10; i++ {
			wg3.Add(1)
			go fanOutWorker(jobs, results, &wg3)
		}
		for i := 1; i <= 50; i++ {
			jobs <- i
		}
		close(jobs)

		select {
		case <-ctx.Done():
			wg3.Wait()
			close(results)
			break Loop
		case <-time.After(3 * time.Second):
			wg3.Wait()
			close(results)
		}

    fmt.Println("[done] restarting cycle...")

		select {
		case <-ctx.Done():
			break Loop
		case <-time.After(1 * time.Second):
		}
	}

	fmt.Println("\n[shutdown] Graceful shutdown complete. Goodbye!")
}
