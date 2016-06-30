package dusl

import (
  "io"
  "fmt"
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

func (this *Syntax) DumpToString(pretty bool) string {
  return dumpToString(this, pretty)
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
  cat := this.Cat
  if cat == "SQ" {
    this.Left.dumpPretty(out, prfx)
    if this.Right != nil && this.Right.Cat != "" {
      this.Right.dumpPretty(out, prfx)
    }
    return
  }
  if cat == "SN" {
    this.Left.dumpPretty(out, prfx)
    if this.Right != nil && this.Right.Cat != "" {
      this.Right.dumpPretty(out, prfx + "| ")
    }
    return
  }
  lit := this.Lit
  if cat == "UN" {
    if lit == "" {
      fmt.Fprintf(out, "%s%s", prfx, this.Ambit.ToString())
    } else {
      fmt.Fprintf(out, "%s%s\n", prfx, lit)
    }
    return
  }
  if cat == "ERR" {
    lit = this.Err
  }
  fmt.Fprintf(out, "%s%s:%s\n", prfx, cat, lit)
  this.Left.dumpPretty(out, prfx + "  ")
  this.Right.dumpPretty(out, prfx + "  ")
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

// Returns a SummaryError for the first n errors found in the tree or nil
// of there were no errors.
func (this *Syntax) ErrorN(n int) error {
  return SummaryError(this.errors(), n)
}

func (this *Syntax) errors() []error {
  return this.gatherErrors(nil)
}

func (this *Syntax) gatherErrors(errs []error) []error {
  if this == nil {
    return errs
  }
  if this.Cat == "ERR" {
    return append(errs, AmbitError( this.Ambit, this.Err))
  }
  errs = this.Left.gatherErrors(errs)
  errs = this.Right.gatherErrors(errs)
  return errs
}

// IsEmpty iff Cat == ""
func (this *Syntax) IsEmpty() bool {
  return this.Cat == ""
}

// IsZeroaryOp iff left and right are empty.
// Give empty string to lit to wildcard match any literal.
func (this *Syntax) IsZeroaryOp(lit string) bool {
  return this.Cat == "OP" && this.Left.IsEmpty() && this.Right.IsEmpty() && (lit == "" || this.Lit == lit)
}

// IsPrefiOp iff left is not empty and right is empty.
// Give empty string to lit to wildcard match any literal.
func (this *Syntax) IsPrefixOp(lit string) bool {
  return this.Cat == "OP" && this.Left.IsEmpty() && !this.Right.IsEmpty() && (lit == "" || this.Lit == lit)
}

// IsPostfixOp iff left is empty and right is not empty.
// Give empty string to lit to wildcard match any literal.
func (this *Syntax) IsPostfixOp(lit string) bool {
  return this.Cat == "OP" && !this.Left.IsEmpty() && this.Right.IsEmpty() && (lit == "" || this.Lit == lit)
}

// IsInfixOp iff left and right are not empty.
// Give empty string to lit to wildcard match any literal.
func (this *Syntax) IsInfixOp(lit string) bool {
  return this.Cat == "OP" && !this.Left.IsEmpty() && !this.Right.IsEmpty() && (lit == "" || this.Lit == lit)
}

// Return the first node in a pre--order, left to right traversal that match cat and lit.
// Give an empty string to either argument for wildcard match.
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

// Return the first N nodes in a pre--order, left to right traversal that match cat and lit.
// Give an empty string to either argument for wildcard match.
// Give negative N for getting all nodes that match.
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
