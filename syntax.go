package dusl

import (
  "io"
  "fmt"
  "bytes"
)

type Syntax struct {
  Cat string
  Lit string
  Err string
  Ambit *Ambit
  OpAmbit *Ambit
  Left *Syntax
  Right *Syntax
}

func (this *Syntax) mapUnparsedAmbits(f func(ambit *Ambit)string) *Syntax {
  if this == nil {
    return nil
  }
  if this.Cat == "UN" {
    return &Syntax{ Cat: "UN", Lit: f(this.Ambit), Err: this.Err, Ambit: this.Ambit, OpAmbit: this.OpAmbit }
  }
  return &Syntax{ Cat: this.Cat, Lit: this.Lit, Err: this.Err, Ambit: this.Ambit, OpAmbit: this.OpAmbit,
                       Left: this.Left.mapUnparsedAmbits(f),
                       Right: this.Right.mapUnparsedAmbits(f) }
}

func (this *Syntax) ToString() string {
  buf := new(bytes.Buffer)
  this.Dump(buf, "", true)
  return buf.String()
}

func (this *Syntax) Dump(out io.Writer, prfx string, pretty bool) {
  if pretty {
    this.dumpPretty(out, prfx)
  } else {
    this.dumpRaw(out, prfx)
  }
}

func (this *Syntax) dumpPretty(out io.Writer, prfx string) {
  if this == nil {
    return
  }
  if this.Cat == "SQ" {
    this.Left.dumpPretty(out, prfx)
    if this.Right != nil && this.Right.Cat != "" {
      this.Right.dumpPretty(out, prfx)
    }
  } else if this.Cat == "SN" {
    this.Left.dumpPretty(out, prfx)
    if this.Right != nil && this.Right.Cat != "" {
      this.Right.dumpPretty(out, prfx + "| ")
    }
  } else if this.Cat == "UN" {
    if this.Lit == "" {
      fmt.Fprintf(out, "%s%s", prfx, this.Ambit.ToString())
    } else {
      fmt.Fprintf(out, "%s%s\n", prfx, this.Lit)
    }
  } else {
    fmt.Fprintf(out, "%s%s:%s\n", prfx, this.Cat, this.Lit)
    this.Left.dumpPretty(out, prfx + "  ")
    this.Right.dumpPretty(out, prfx + "  ")
  }
}

func (this *Syntax) dumpRaw(out io.Writer, prfx string) {
  if this == nil {
    fmt.Fprintf(out, "%snil\n", prfx)
    return
  }
  fmt.Fprintf(out, "%s%s:%s:%s:%s\n", prfx, this.Cat, this.Lit, this.Err, this.Ambit.String())
  if this.Left != nil { this.Left.dumpRaw(out, prfx + "  ") }
  if this.Right != nil { this.Right.dumpRaw(out, prfx + "  ") }
}

func (this *Syntax) ErrorN(n int) error {
  return errorN(this.Errors(), n)
}

func (this *Syntax) Errors() []error {
  return this.gatherErrors(nil)
}

func (this *Syntax) gatherErrors(errs []error) []error {
  if this == nil {
    return errs
  }
  if this.Cat == "ERR" {
    return append(errs, fmt.Errorf("%s: %s", this.Ambit.Location(), this.Err))
  }
  errs = this.Left.gatherErrors(errs)
  errs = this.Right.gatherErrors(errs)
  return errs
}

func (this *Syntax) IsEmpty() bool {
  return this.Cat == ""
}

func (this *Syntax) IsZeroaryOp(lit string) bool {
  if this.Cat != "OP" {
    return false
  }
  if this.Left.Cat != "" || this.Right.Cat != "" {
    return false
  }
  return this.Lit == lit
}

func (this *Syntax) IsPrefixOp(lit string) bool {
  if this.Cat != "OP" {
    return false
  }
  if this.Left.Cat != "" || this.Right.Cat == "" {
    return false
  }
  return this.Lit == lit
}

func (this *Syntax) IsPostfixOp(lit string) bool {
  if this.Cat != "OP" {
    return false
  }
  if this.Left.Cat == "" || this.Right.Cat != "" {
    return false
  }
  return this.Lit == lit
}

func (this *Syntax) IsInfixOp(lit string) bool {
  if this.Cat != "OP" {
    return false
  }
  if this.Left.Cat == "" || this.Right.Cat == "" {
    return false
  }
  return this.Lit == lit
}

func (this *Syntax) First(cat, lit string) *Syntax {
  if this == nil {
    return nil
  }
  if (cat == "" || this.Cat == cat) && (lit == "" || this.Lit == lit) {
    return this
  }
  node := this.Left.First(cat, lit)
  if node != nil {
    return node
  }
  return this.Right.First(cat, lit)
}

func (this *Syntax) FirstN(cat, lit string, n int) []*Syntax {
  list := make([]*Syntax, 0, max(n, 1))
  return this.listFirstN(list, cat, lit, n)
}

func (this *Syntax) listFirstN(list []*Syntax, cat, lit string, n int) []*Syntax {
  if this == nil || (n >= 0 && len(list) >= n) {
    return list
  }
  if (cat == "" || this.Cat == cat) && (lit == "" || this.Lit == lit) {
    list = append(list, this)
  }
  list = this.Left.listFirstN(list, cat, lit, n)
  return this.Right.listFirstN(list, cat, lit, n)
}
