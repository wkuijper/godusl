package dusl

import (
  "fmt"
)

// An Ambit represent a region inside some Source, identified by a Start
// (byte-offset) and an end (byte-offset).
type Ambit struct {
  Source *Source
  Start int
  End int
}

// StringAmbit creates an Ambit representing a given string. Useful for unit testing.
func StringAmbit(s string) *Ambit {
  return StringSource(s).FullAmbit()
}

func (this *Ambit) String() string {
  return fmt.Sprintf("%s[%d:%d]", this.Source.Path, this.Start, this.End)
}

// ToString returns the literal source fragment that this ambit represents as a string.
func (this *Ambit) ToString() string {
  return string(this.Source.Text[this.Start:this.End])
}

// Location returns a string containing the path, start-line and -column and
// end-line and -column for this Ambit. Beware that this operation is expensive as
// it potentially requires scanning the entire source. Calling this repeatedly may
// lead to quadratic time behaviour.
func (this *Ambit) Location() string {
  source := this.Source
  startLine, startColumn := source.LineColumn(this.Start)
  endLine, endColumn := source.LineColumn(this.End)
  if startLine != endLine {
    return fmt.Sprintf("%s:%d:%d:%d:%d", source.Path, startLine, startColumn, endLine, endColumn)
  }
  return fmt.Sprintf("%s:%d:%d:%d", source.Path, startLine, startColumn, endColumn)
}

// Merge returns the smallest ambit that extends from this ambit to the given ambit.
func (this *Ambit) Merge(that *Ambit) *Ambit {
  return &Ambit{Source: this.Source, Start: min(this.Start, that.Start), End: max(this.End, that.End)}
}

// SubtractLeft returns the largest ambit that starts at this ambit start and does not
// overlap in any character position with the given ambit.
func (this *Ambit) SubtractLeft(that *Ambit) *Ambit {
  return &Ambit{Source: this.Source, Start: min(this.End, max(this.Start, that.End)), End: this.End}
}

// SubtractRight returns the largest ambit that end at this ambit end and does not
// overlap in any character position with the given ambit.
func (this *Ambit) SubtractRight(that *Ambit) *Ambit {
  return &Ambit{Source: this.Source, Start: this.Start, End: max(this.Start, min(this.End, that.Start))}
}

// SplitAbs returns two ambits by splitting at the given absolute byte offset.
func (this *Ambit) SplitAtAbs(i int) (*Ambit, *Ambit) {
  return this.To(i), this.From(i)
}

// SplitLine returns two ambits by splitting at the first encountered end-of-line for this ambit.
func (this *Ambit) SplitLine() (*Ambit, *Ambit) {
  text := this.Source.Text
  for i := this.Start; i < this.End; i++ {
    c := text[i]
    if c == '\r' {
      if i+1 < this.End && text[i+1] == '\n' {
        return this.To(i+2), this.From(i+2)
      }
      return this.To(i+1),this.From(i+1)
    } else if c == '\n' {
      return this.To(i+1), this.From(i+1)
    }
  }
  return this, this.From(this.End)
}

// StripIndent returns the number of leading spaces and the remaining ambit for this ambit.
func (this *Ambit) StripIndent() (int, *Ambit) {
  text := this.Source.Text
  for i := this.Start; i < this.End; i++ {
    c := text[i]
    if c != ' ' {
      return i-this.Start, this.From(i)
    }
  }
  return 0, this
}

// IsEmpty returns true iff start equals end for this ambit.
func (this *Ambit) IsEmpty() bool {
  return this.Start == this.End
}

// To returns the ambit from the start of this ambit to the given absolute byte offset.
// The given offset must lie within this ambit.
func (this *Ambit) To(i int) *Ambit {
  return &Ambit{ Source: this.Source, Start: this.Start, End: max(i, this.Start) }
}

// From returns the ambit from the given absolute byte offset to the end of this ambit.
// The given offset must lie within this ambit.
func (this *Ambit) From(i int) *Ambit {
  return &Ambit{ Source: this.Source, Start: min(i, this.End), End: this.End }
}

// CollapseLeft returns the empty ambit that starts and end where this ambit starts.
func (this *Ambit) CollapseLeft() *Ambit {
  return &Ambit{ Source: this.Source, Start: this.Start, End: this.Start }
}

// CollapseRight returns the empty ambit that starts and end where this ambit ends.
func (this *Ambit) CollapseRight() *Ambit {
  return &Ambit{ Source: this.Source, Start: this.End, End: this.End }
}

// Shift returns this ambit with both start and end shifted by the given amount.
func (this *Ambit) Shift(offset int) *Ambit {
  return &Ambit{ Source: this.Source, Start: this.Start + offset, End: this.End + offset }
}

// Return true iff this ambit consists entirely out of ASCII whitespace characters.
func (this *Ambit) IsWhitespace() bool {
  text := this.Source.Text
  for i := this.Start; i < this.End; i++ {
    c := text[i]
    if !(c == ' ' || c == '\t' || c == '\n' || c == '\r') {
      return false
    }
  }
  return true
}

// Returns true iff the first byte of this ambit equals the given byte.
func (this *Ambit) FirstByteIs(b byte) bool {
  if this.Start >= this.End {
    return false
  }
  return this.Source.Text[this.Start] == b
}