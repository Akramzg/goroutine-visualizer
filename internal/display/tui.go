package display

import (
	"cmp"
	"fmt"
	"slices"
	"strings"

	"github.com/Akramzg/goroutine-vis/internal/parser"
)


func Render(pid int ,port int, status map[string]string, goroutines []parser.Goroutine){
  states := make(map[string]int)
  topFrames := make(map[string]int)

  for _, g := range goroutines{
    states[g.State]++
    topFrames[g.TopFrame]++
  }

  //sort state keys
  keyStates := make([]string, 0, len(states))
  for k := range states{
    keyStates = append(keyStates, k)
  }

  slices.SortFunc(keyStates, func(a, b string)int{
    return cmp.Compare(states[b],states[a])
  })
  

  // sort topframe keys
  keyTopFrames := make([]string, 0, len(topFrames))
  for k := range topFrames{
    keyTopFrames = append(keyTopFrames,k)
  }
  slices.SortFunc(keyTopFrames, func(a, b string) int {
    return cmp.Compare(topFrames[b],topFrames[a])
  })

  // Header output
  fmt.Printf("Starting Goroutine Visualizer for PID %d on PORT %d\n", pid, port)
  fmt.Printf("PID: %d - %s - Threads %s - VmRSS %s\n",pid, status["Name"],status["Threads"],status["VmRSS"])
  fmt.Printf("Total goroutines: %d\n\n", len(goroutines))
  if len(goroutines)==0{
    fmt.Println("No goroutine to display")
  }

  // Bar chart output
  maxCount := states[keyStates[0]]
  maxBarWidth := 20
  fmt.Println("--- STATES ---")
  for _, k := range keyStates{
    count := states[k]
    filledBlocks := (count*maxBarWidth)/maxCount
    emptyBlocks := maxBarWidth - filledBlocks
    bar := strings.Repeat("█",filledBlocks)+strings.Repeat("░",emptyBlocks)
    fmt.Printf("%-15s |%s| %d\n", k, bar, count)
  }

  // Top frames output
  fmt.Println("\n--- Top Frames ---")
  for _, k := range keyTopFrames{
    count := topFrames[k]
    fmt.Printf("%-50s [%dx]\n", k,count)
  }
  


}
