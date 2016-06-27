package dusl

import (
  "os"
  "io/ioutil"
)

// A Source object represents a source text loaded into memory for parsing.
// The Path and LineOffset fields are used in error reporting. For a normal
// source file LineOffset should be 0.
type Source struct {
  Path string
  LineOffset int
  Text []byte
}

// SourceFromString creates a source object from a given string. Useful for unit testing.
func SourceFromString(s string) *Source {
  return &Source{ Path: "str", Text: []byte(s) }
}

// Create a source object by reading in a given filepath.
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

// Compute the line and column values corresponding to a given position
// (byte-offset) into the source.  Beware that this operation is expensive as
// it potentially requires scanning the entire source. Calling this repeatedly may
// lead to quadratic time behaviour.
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

// Create an ambit that represent the entire source, minus the BOM (byte order mark)
// if present and minus the first byte if this byte corresponds to a '#'. The latter
// convention helps ignoring hashbang line and allows DUSL sources to be easily
// embedded inside a markdown document.
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
