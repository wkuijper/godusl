package dusl

import (
  "os"
  "io/ioutil"
)

type Source struct {
  Path string
  LineOffset int
  Text []byte
}

func StringSource(s string) *Source {
  return &Source{ Path: "str", Text: []byte(s) }
}

func SourceFromPath(path string) (*Source, error) {
  in, err := os.Open(path)
  if err != nil {
    return nil, err
  }
  text, err := ioutil.ReadAll(in)
  if err != nil {
    return nil, err
  }
  return &Source{ Path : path, Text : text }, nil
}

func (this *Source) LineColumn(pos int) (int, int) {
  lineNum := 1 + this.LineOffset
  colNum := 0
  text := this.Text
  justReadCR := false
  for i := 0; i < pos; i++ {
    c := text[i]
    if c == '\r' {
      lineNum++
      colNum = 0
      justReadCR = true
    } else if c == '\n' {
      if !justReadCR {
        lineNum++
        colNum = 0
      }
      justReadCR = false
    } else {
      colNum++
      justReadCR = false
    }
  }
  return lineNum, colNum
}

func (this *Source) FullAmbit() *Ambit {
  text := this.Text
  start := 0
  end := len(text)
  // Check for Byte Order Mark (BOM)
  if end >= 3 &&  text[0] == 0xEF && text[1] == 0xBB && text[2] == 0xBF { 
    start += 3
  }
  ambit := &Ambit{ Source: this, Start: start, End: end }
  // Check for hash (to avoid shebang and/or markdown header)
  if start < end && text[start] == '#' {
    _, ambit = ambit.SplitLine()
  }
  return ambit
}
