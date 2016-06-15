package dusl

import (
  "fmt"
)

type Ambit struct {
  Source *Source
  Start int
  End int
}

func StringAmbit(s string) *Ambit {
  return StringSource(s).FullAmbit()
}

func (this *Ambit) String() string {
  return fmt.Sprintf("%s[%d:%d]", this.Source.Path, this.Start, this.End)
}

func (this *Ambit) ToString() string {
  return string(this.Source.Text[this.Start:this.End])
}

func (this *Ambit) Location() string {
  source := this.Source
  startLine, startColumn := source.LineColumn(this.Start)
  endLine, endColumn := source.LineColumn(this.End)
  if startLine != endLine {
    return fmt.Sprintf("%s:%d:%d:%d:%d", source.Path, startLine, startColumn, endLine, endColumn)
  }
  return fmt.Sprintf("%s:%d:%d:%d", source.Path, startLine, startColumn, endColumn)
}

func (this *Ambit) Merge(that *Ambit) *Ambit {
  return &Ambit{Source: this.Source, Start: min(this.Start, that.Start), End: max(this.End, that.End)}
}

func (this *Ambit) SubtractLeft(that *Ambit) *Ambit {
  return &Ambit{Source: this.Source, Start: min(this.End, max(this.Start, that.End)), End: this.End}
}

func (this *Ambit) SubtractRight(that *Ambit) *Ambit {
  return &Ambit{Source: this.Source, Start: this.Start, End: max(this.Start, min(this.End, that.Start))}
}

func (this *Ambit) SplitAtAbs(i int) (*Ambit, *Ambit) {
  return this.To(i), this.From(i)
}

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

func (this *Ambit) IsEmpty() bool {
  return this.Start == this.End
}

func (this *Ambit) To(i int) *Ambit {
  return &Ambit{ Source: this.Source, Start: this.Start, End: max(i, this.Start) }
}

func (this *Ambit) From(i int) *Ambit {
  return &Ambit{ Source: this.Source, Start: min(i, this.End), End: this.End }
}

func (this *Ambit) CollapseLeft() *Ambit {
  return &Ambit{ Source: this.Source, Start: this.Start, End: this.Start }
}

func (this *Ambit) CollapseRight() *Ambit {
  return &Ambit{ Source: this.Source, Start: this.End, End: this.End }
}

func (this *Ambit) Shift(offset int) *Ambit {
  return &Ambit{ Source: this.Source, Start: this.Start + offset, End: this.End + offset }
}

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

func (this *Ambit) FirstByteIs(b byte) bool {
  if this.Start >= this.End {
    return false
  }
  return this.Source.Text[this.Start] == b
}