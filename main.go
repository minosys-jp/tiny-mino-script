package main

import (
  "os"
  "fmt"
  "bufio"
)

func isspace(c byte) bool {
  ref := []byte(" \r\n\t")
  for _, s := range ref {
    if s == c {
      return true
    }
  }
  return false
}

func isalpha(c byte) bool {
  ref := []byte("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz_")
  for _, s := range ref {
    if s == c {
      return true
    }
  }
  return false
}

func isnumber(c byte) bool {
  ref := []byte("0123456789")
  for _, s := range ref {
    if s == c {
      return true
    }
  }
  return false
}

func splitter(data []byte, isEof bool) (int, []byte, error) {
  if data != nil {
    // space ならば無視する
    idxs := 0
    for idx, c := range data {
      if (!isspace(c)) {
        idxs = idx
        break
      }
    }
    mode := 0
    trans := map[byte]int{0x2e: 3, 0x2d: 4,
      0x2b: 5, 0x3c: 6, 0x3d: 7, 0x3e: 8, 0x33: 9, 0x2a: 10,
      0x22: 0x22, 0x27: 0x27, 0x60: 0x60}
    stripped := data[idxs:]

    for idx, c := range stripped {
      if mode == 0 {
        if isalpha(c) {
          // keyword or name
          mode = 1
        } else if isnumber(c) {
          // integer or fraction
          mode = 2
        } else {
          m, e := trans[c]
          if e == true {
            mode = m
          }
        }
        if mode == 0 {
          return idxs + idx + 1, stripped[0:(idx + 1)], nil
        }
      } else if mode == 1 {
        // name
        if !(isalpha(c) || isnumber(c)) {
          return idxs + idx, stripped[0:idx], nil
        }
      } else if mode == 2 {
        // number
        if !isnumber(c) && c != 0x2e {
          return idxs + idx, stripped[0:idx], nil
        }
        if (c == 0x2e) {
          // fraction
          mode = 3
        }
      } else if mode == 3 {
        // fraction
        if !isnumber(c) {
          return idxs + idx, stripped[0:idx], nil
        }
      } else if mode == 4 || mode == 5 {
        if !isnumber(c) {
          if c == 0x2e {
            // fraction
            mode = 3
          } else if c == 0x3d {
            // -=
            return idxs + idx + 1, stripped[0:(idx + 1)], nil
          } else if isspace(c) {
            // mono +/-
            return idxs + idx, stripped[0:idx], nil
          }
        } else {
          mode = 2
        }
      } else if mode >= 6 && mode <= 10 {
        if c != 0x3d {
          return idxs + idx, stripped[0:idx], nil
        }
        return idxs + idx + 1, stripped[0:(idx + 1)], nil
      } else if mode == 0x22 {
        if c == 0x22 {
          return idxs + idx + 1, stripped[0:(idx + 1)], nil
        }
      } else if mode == 0x27 {
        if c == 0x27 {
          return idxs + idx + 1, stripped[0:(idx + 1)], nil
        }
      } else if mode == 0x60 {
        if c == 0x60 {
          return idxs + idx  + 1, stripped[0:(idx + 1)], nil
        }
      }
    }
  }
  return 0, nil, nil
}

func main() {
  f, err := os.Open("sample.tmscr")
  if err != nil {
    fmt.Println("file error")
    return
  }
  defer f.Close()
  scanner := bufio.NewScanner(f)
  scanner.Split(splitter)
  for scanner.Scan() {
    fmt.Println(scanner.Text())
  }
}
