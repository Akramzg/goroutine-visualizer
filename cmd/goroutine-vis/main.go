package main
import (
  "fmt"
  "github.com/Akramzg/goroutine-vis/internal/proc"
  "log"
  "flag"

)


 func main(){

  var port int
  var pid int

  flag.IntVar(&pid, "pid", 0, "Process Id to visualize")
  flag.IntVar(&port, "port",8080, "Port number to listen on")
  flag.Parse()


  if pid <= 0{
    fmt.Println("Usage: goroutine-vis -pid=<target-pid>")
    flag.PrintDefaults()
    log.Fatalf("Missing or invalid PID")
  }
  if(port<1 || port>65535){
    log.Fatalf("Invalid port number %d", port)
  }
  fmt.Printf("Starting Goroutine Visualizer for PID %d on PORT %d\n", pid, port)


  status, err := proc.ReadStatus(pid)
  if err != nil{
    log.Fatalf("Failed to read status: %v", err)
  }

  fmt.Printf("PID: %d - %s - Threads %s - VmRSS %s\n",pid, status["Name"],status["Threads"],status["VmRSS"])
  
  // fetch and print raw goroutine dump

  portStr := fmt.Sprintf("%d", port)
  dump, err := proc.FetchGor(portStr)
  if err != nil{
    log.Fatalf("Failed to fetch goroutine dump:", err)
  }

  fmt.Println("\n--- Raw Goroutine Dump ---")
  fmt.Println(dump)



} 
