package main
import (
  "fmt"
  "os"
  "github.com/Akramzg/goroutine-vis/internal/proc"
  "log"
  "strconv"

)


 func main(){

  if len(os.Args) < 2 {
    fmt.Println("Usage: goroutine-vis <target-pid>")
    os.Exit(1)
  }

  pidStr := os.Args[1]
  pid, err := strconv.Atoi(pidStr)
  if err!=nil{
    log.Fatalf("Invalid PID '%s': must be a number",pidStr)
  }

  fmt.Printf("Starting Goroutine Visualizer for PID: %s\n", pidStr)
  
  status, err := proc.ReadStatus(pid)
  if err != nil{
    log.Fatalf("Failed to read status: %v", err)
  }

  fmt.Printf("Process Name: %s\n", status["Name"])
  fmt.Printf("State:      %s\n", status["State"])
  fmt.Printf("Parent Process ID %s\n", status["PPid"])
  fmt.Printf("Number of Threads %s\n", status["Threads"])
  fmt.Printf("Virtual memory allocated %s\n", status["VmPeak"])
  fmt.Printf("Physical memory currently occupied by this process %s\n", status["VmRSS"])








} 
