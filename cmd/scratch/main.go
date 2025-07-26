package main

import (
  "net/http"
  _ "net/http/pprof"
  "time"
)

func main(){
  ch := make(chan int)
  for i :=0; i<5; i++{
    go func(){
      <-ch
    }()
  }


  for i:=0;i<3;i++{
    go func(){ time.Sleep(time.Hour)}()
  }

  http.ListenAndServe(":6060",nil)
} 
