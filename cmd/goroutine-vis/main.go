package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Akramzg/goroutine-vis/internal/display"
	"github.com/Akramzg/goroutine-vis/internal/parser"
	"github.com/Akramzg/goroutine-vis/internal/proc"
)

func main() {
	// flags — all required except interval which defaults to 2s
	var pid, port, interval int
	flag.IntVar(&pid, "pid", 0, "target process pid")
	flag.IntVar(&port, "port", 6060, "pprof port on target process")
	flag.IntVar(&interval, "interval", 2, "refresh interval in seconds")
	flag.Parse()

	// basic validation before we even try to connect
	if pid <= 0 {
		fmt.Println("usage: goroutine-vis -pid=<pid> [-port=6060] [-interval=2]")
		flag.PrintDefaults()
		log.Fatalf("missing or invalid pid")
	}
	if port < 1 || port > 65535 {
		log.Fatalf("invalid port: %d", port)
	}
	if interval <= 0 {
		log.Fatalf("interval must be > 0")
	}

	// intercept ctrl-c and sigterm so we exit cleanly
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sigChan
		fmt.Print("\033[H\033[2J")
		fmt.Println("stopped. bye!")
		os.Exit(0)
	}()

	portStr := fmt.Sprintf("%d", port)

	// renderFrame does the full fetch → parse → render cycle
	// on any error it warns and bails early — next tick will retry
	renderFrame := func() {
		status, err := proc.ReadStatus(pid)
		if err != nil {
			fmt.Printf("\033[H\033[2J\n[warning] /proc read failed: %v\n", err)
			return
		}

		dump, err := proc.FetchGor(portStr)
		if err != nil {
			// target process might be restarting — just wait it out
			fmt.Printf("\033[H\033[2J\n[warning] pprof unreachable: %v\n", err)
			return
		}

		goroutines, err := parser.Parse(dump)
		if err != nil {
			fmt.Printf("\033[H\033[2J\n[warning] parse failed: %v\n", err)
			return
		}

		fmt.Print("\033[H\033[2J")
		display.Render(pid, port, status, goroutines)
	}

	// render once immediately so there's no blank screen at startup
	renderFrame()

	// tick every N seconds and re-render
	ticker := time.NewTicker(time.Duration(interval) * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		renderFrame()
	}
}
