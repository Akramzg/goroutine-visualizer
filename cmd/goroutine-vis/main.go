package main
import (
  "fmt"
  "os"
)


 func main(){
  if len(os.Args) < 2 {
    fmt.Println("Usage: goroutine-vis <target-pid>")
    os.Exit(1)
  }

  pid := os.Args[1]
  fmt.Printf("Starting Goroutine Visualizer for PID: %s\n", pid)

}
