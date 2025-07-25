package proc
import(
  "bufio"
  "fmt"
  "os"
  "strings"
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
