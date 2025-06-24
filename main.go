package main

import (
  "os"
  "fmt"
  "bufio"
  "github.com/minosys-jp/tiny-mino-script/internal/myScanner"
)

func main() {
  f, err := os.Open("sample.tmscr")
  if err != nil {
    fmt.Println("file error")
    return
  }
  defer f.Close()
  scanner := myScanner.CreateScanner(f)
  for scanner.Scan() {
    fmt.Println(scanner.Text())
  }
}
