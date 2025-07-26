package proc
import(
  "bufio"
  "fmt"
  "os"
  "strings"
  "net/http"
  "io"
)



func ReadStatus(pid int) (map[string]string, error){
  stat := make(map[string]string)
  //prep & open the file
  path := fmt.Sprintf("/proc/%d/status",pid)
  file, err := os.Open(path)
  if err != nil{
    return nil, err
  }
  defer file.Close()

  //scan line by line
  scanner := bufio.NewScanner(file)
  for scanner.Scan(){
    line := scanner.Text()
    slices := strings.SplitN(line,":",2)

    // guard against empty lines & trim extra whitespace
    if len(slices)==2{
      stat[slices[0]]=strings.TrimSpace(slices[1])
    }
  }

  if err := scanner.Err(); err != nil{
    return nil,err
  }

  return stat, nil

}

func FetchGor(port string) (string, error){
  url := fmt.Sprintf("http://localhost:%s/debug/pprof/goroutine?debug=2",port)

  // making GET req
  resp, err := http.Get(url)
  if err != nil{
    return "", fmt.Errorf("failed to connect to target app: %v",err)
  }
  defer resp.Body.Close()
  if resp.StatusCode != http.StatusOK{
    return "", fmt.Errorf("bad status code received: %d", resp.StatusCode)
  }
  body, err := io.ReadAll(resp.Body)
  if err != nil{
    return "", fmt.Errorf("failed to read response body: %v", err)
  }
  return string(body), nil
}
