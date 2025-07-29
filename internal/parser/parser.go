package parser

import(
  "strings"
  "bufio"
  "fmt"
  "strconv"
)

type Goroutine struct{
  ID  int
  State string
  TopFrame string
}

func Parse(raw string) ([]Goroutine, error){

  if strings.TrimSpace(raw) == ""{
    return nil, fmt.Errorf("empty raw dump")
  }

  var goroutines []Goroutine
  scanner := bufio.NewScanner(strings.NewReader(raw))

  var currentGoroutine *Goroutine
  captureNextLineAsTopFrame := false
  lineNum := 0

  for scanner.Scan(){
    lineNum++
    line := scanner.Text()
    //detecting goroutine header
    if strings.HasPrefix(line, "goroutine "){
      if currentGoroutine != nil{
        goroutines = append(goroutines, *currentGoroutine)
      }
      fields := strings.Fields(line)
      if len(fields)<2{
        return nil, fmt.Errorf("malformed goroutine header at line %d: %s", lineNum, line)
      }
      id, err := strconv.Atoi(fields[1])
      if err != nil{
        return nil, fmt.Errorf("invalid goroutine ID at line %d: %v", lineNum, err)
      }
      state := extractState(line)
      if state == ""{
        return nil, fmt.Errorf("could not extract state from header at line %d: %s",lineNum, line)
      }
      currentGoroutine = &Goroutine{
        ID: id,
        State: state,
      }
      captureNextLineAsTopFrame = true
      continue
    }

    if captureNextLineAsTopFrame {
      trimmed := strings.TrimSpace(line)
      if trimmed != ""{
        topFrame := strings.Split(trimmed, "(")[0]
        if currentGoroutine != nil{
          currentGoroutine.TopFrame = topFrame
        }
        captureNextLineAsTopFrame = false
      }
    }
  }
  if err := scanner.Err(); err !=nil{
    return nil, fmt.Errorf("error reading dump: %v", err)
  }
  if currentGoroutine != nil{
    goroutines = append(goroutines, *currentGoroutine)
  }
  return goroutines, nil
}

func extractState(line string) string {
  start := strings.Index(line, "[")
  end := strings.LastIndex(line,"]")
  if start == -1 || end == -1 || start >=end{
    return ""
  }
  return line[start+1 : end]
}
