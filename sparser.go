package dusl

import (
  "fmt"
)

// Sparser stands for Superpermissive-Parser.
// The Sparse method converts an ambit into a syntax tree,
// the SparseUndent method converts an entire source into a syntax tree.
type Sparser interface {
  Sparse(ambit *Ambit) *Syntax
  SparseUndent(src *Source) *Syntax
}

type sparser struct {
  precedenceLevels
  spanner spannerI
}

func newSparser(spanner spannerI, precedence *precedenceLevels) Sparser {
  return &sparser{ spanner: spanner, precedenceLevels: *precedence }
}

// SparseUndent returns the syntax tree constructed for the given source.
func (this *sparser) SparseUndent(source *Source) *Syntax {
  root := Undent(source)
  this.sparseSQ(root)
  return root
}

func (this *sparser) sparseSQ(node *Syntax) {
  if node.Cat == "SQ" {
    this.sparseSQ(node.Left)
    this.sparseSQ(node.Right)
    return
  }
  if node.Cat == "SN" {
    node.Left = this.Sparse(node.Left.Ambit)
    this.sparseSQ(node.Right)
    return
  }
  // <empty
}

// Sparse returns the syntax tree constructed for the given ambit.
func (this *sparser) Sparse(ambit *Ambit) *Syntax {
  return this.sparse(ambit, this.spanner.span(ambit), 1)
}

func (this *sparser) sparse(ambit *Ambit, spans []*spanT, minPrecedence int) *Syntax {
  ambit, spans = trimSpans(ambit, spans)
  if len(spans) == 0 {
    return &Syntax{ Ambit: ambit }
  }
  if len(spans) == 1 {
    span := spans[0]
    lit := span.Lit
    if span.Children == nil {
      if span.Cat == "OP" {
        if this.precedenceEFE[lit] < minPrecedence {
          return &Syntax{ Cat: "ERR", Err: fmt.Sprintf("unexpected: %s", lit), Ambit: ambit }
        }
        return &Syntax{ Cat: "OP", Lit: lit, Ambit: ambit, OpAmbit: span.Ambit,
                        Left: &Syntax{ Ambit: span.Ambit.CollapseLeft() },
                        Right: &Syntax{ Ambit: span.Ambit.CollapseRight() } }
      }
      return &Syntax{ Cat: span.Cat, Lit: lit, Err: span.Err, Ambit: ambit, OpAmbit: span.Ambit }
    }
    // span.Cat == "BB"
    precedence, recognized := this.precedenceB[lit]
    if !recognized {
      return &Syntax{ Cat: "ERR", Err: fmt.Sprintf("unexpected: %s", lit), Ambit: ambit }
    }
    return &Syntax{ Cat: span.Cat, Lit: lit, Ambit: span.Ambit,
                    Left: this.sparse(span.SubAmbit, span.Children, precedence),
                    Right: &Syntax{ Ambit: span.Ambit.CollapseRight() } }
  }
  l, splitPrecedence, splitLoc, splitPrecLeft, splitPrecRight := len(spans)-1, maxPrecedence+1, -1, -1, -1
  if span, ws := spans[0], spans[1]; span.Cat != "OP" && ws.Cat == "WS" {
    juxtaposition := false
    if len(spans) == 2 {
      juxtaposition = true
    } else {
      cat := spans[2].Cat
      if cat != "ERR" {
        if cat != "OP" {
          juxtaposition = true
        } else {
          lit := spans[2].Lit
          if this.exclusivelyZeroAry[lit] {
            juxtaposition = true
          }
        }
      }
    }
    if juxtaposition {
      return &Syntax{ Cat: "JUXT", Lit: " ", Ambit: ambit, OpAmbit: ws.Ambit,
                      Left: this.sparse(span.Ambit, spans[:1], minPrecedence),
                      Right: this.sparse(ambit.SubtractLeft(ws.Ambit), spans[1:], minPrecedence) }
    }
  }
  if span := spans[0]; span.Cat == "OP" {
    lit := span.Lit
    prec := this.precedenceEFA[lit]
    if prec == minPrecedence {
      return &Syntax{ Cat: span.Cat, Lit: lit, Ambit: ambit, OpAmbit: span.Ambit,
                      Left: &Syntax{ Ambit: span.Ambit.CollapseLeft() },
                      Right: this.sparse(ambit.SubtractLeft(span.Ambit), spans[1:], prec) }
    }
    if prec > minPrecedence && prec < splitPrecedence {
      splitLoc, splitPrecedence, splitPrecRight = 0, prec, prec
    } 
  }
  if span := spans[l]; span.Cat == "OP" {
    lit := span.Lit
    prec := this.precedenceAFE[lit]
    if prec == minPrecedence {
      return &Syntax{ Cat: span.Cat, Lit: lit, Ambit: ambit, OpAmbit: span.Ambit,
                      Left: this.sparse(ambit.SubtractRight(span.Ambit), spans[:l], prec),
                      Right: &Syntax{ Ambit: span.Ambit.CollapseRight() } }
    }
    if prec >= minPrecedence && prec < splitPrecedence {
      splitLoc, splitPrecedence, splitPrecLeft = l, prec, prec
    }
  }
  for indexLR := 1; indexLR < l; indexLR++ {
    if span := spans[indexLR]; span.Cat == "OP" {
      lit := span.Lit
      prec := this.precedenceAFB[lit]
      if prec >= minPrecedence && prec < splitPrecedence {
        if this.checkInfixCandidate(spans, indexLR, prec, prec+1) {
          if prec == minPrecedence {
            return &Syntax{ Cat: span.Cat, Lit: lit, Ambit: ambit, OpAmbit: span.Ambit, 
                            Left: this.sparse(ambit.SubtractRight(span.Ambit), spans[:indexLR], prec+1),
                            Right: this.sparse(ambit.SubtractLeft(span.Ambit), spans[indexLR+1:], prec) }
          }
          splitLoc, splitPrecedence, splitPrecLeft, splitPrecRight = indexLR, prec, prec+1, prec
        }
      }
    }
    indexRL := l - indexLR
    if span := spans[indexRL]; span.Cat == "OP" {
      lit := span.Lit
      prec := this.precedenceBFA[lit]
      if prec >= minPrecedence && prec < splitPrecedence {
        if this.checkInfixCandidate(spans, indexRL, prec+1, prec) {
          if prec == minPrecedence {
            return &Syntax{ Cat: span.Cat, Lit: lit, Ambit: ambit, OpAmbit: span.Ambit, 
                            Left: this.sparse(ambit.SubtractRight(span.Ambit), spans[:indexRL], prec),
                            Right: this.sparse(ambit.SubtractLeft(span.Ambit), spans[indexRL+1:], prec+1) }
          }
          splitLoc, splitPrecedence, splitPrecLeft, splitPrecRight = indexRL, prec, prec, prec+1
        }
      }
    }
  }
  if splitLoc >= 0 {
    splitSpan := spans[splitLoc]
    if splitLoc == 0 {
      return &Syntax{ Cat: splitSpan.Cat, Lit: splitSpan.Lit, Ambit: ambit, OpAmbit: splitSpan.Ambit,
                      Left: &Syntax{ Ambit: splitSpan.Ambit.CollapseLeft() },
                      Right: this.sparse(ambit.SubtractLeft(splitSpan.Ambit), spans[1:], splitPrecRight) }
    } 
    if splitLoc == l {
      return &Syntax{ Cat: splitSpan.Cat, Lit: splitSpan.Lit, Ambit: ambit, OpAmbit: splitSpan.Ambit,
                      Left: this.sparse(ambit.SubtractRight(splitSpan.Ambit), spans[:l], splitPrecLeft),
                      Right: &Syntax{ Ambit: splitSpan.Ambit.CollapseRight() } }
    }
    return &Syntax{ Cat: splitSpan.Cat, Lit: splitSpan.Lit, Ambit: ambit, OpAmbit: splitSpan.Ambit, 
                    Left: this.sparse(ambit.SubtractRight(splitSpan.Ambit), spans[:splitLoc], splitPrecLeft),
                    Right: this.sparse(ambit.SubtractLeft(splitSpan.Ambit), spans[splitLoc+1:], splitPrecRight) }
  }
  firstSpan, secondSpan := spans[0], spans[1]
  if secondSpan.Cat == "WS" {
    return &Syntax{ Cat: "JUXT", Lit: " ", Ambit: ambit, OpAmbit: secondSpan.Ambit,
                    Left: this.sparse(ambit.SubtractRight(secondSpan.Ambit), spans[:1], minPrecedence),
                    Right: this.sparse(ambit.SubtractLeft(secondSpan.Ambit), spans[2:], minPrecedence) }
  }
  return &Syntax{ Cat: "GLUE", Lit: "", Ambit: ambit, OpAmbit: secondSpan.Ambit.CollapseLeft(),
                  Left: this.sparse(ambit.SubtractRight(secondSpan.Ambit), spans[:1], minPrecedence),
                  Right: this.sparse(ambit.SubtractLeft(firstSpan.Ambit), spans[1:], minPrecedence) }
}

func (this *sparser) checkInfixCandidate(spans []*spanT, index int, minPrecLeft int, minPrecRight int) bool {  
  indexRL := index-1
  for indexRL >= 0 {
    span := spans[indexRL]
    if span.Cat != "WS" {
      if span.Cat != "OP" {
        break
      }
      lit := span.Lit
      prec := this.precedenceEFE[lit]
      if prec >= minPrecLeft {
        break
      }
      prec = this.precedenceAFE[lit]
      if prec < minPrecLeft {
        return false
      }
      minPrecLeft = prec
    }
    indexRL--
  }
  if indexRL < 0 {
    return false
  }
  l := len(spans)-1
  indexLR := index+1 
  for indexLR <= l {
    span := spans[indexLR]
    if span.Cat != "WS" {
      if span.Cat != "OP" {
        break
      }
      lit := span.Lit
      prec := this.precedenceEFE[lit]
      if prec >= minPrecRight {
        break
      }
      prec = this.precedenceEFA[lit]
      if prec < minPrecRight {
        return false
      }
      minPrecRight = prec
    }
    indexLR++
  }
  if indexLR > l {
    return false
  }
  return true
}



func trimSpans(ambit *Ambit, spans []*spanT) (*Ambit, []*spanT) {
  return trimSpansLeft(trimSpansRight(ambit, spans))
}

func trimSpansLeft(ambit *Ambit, spans []*spanT) (*Ambit, []*spanT) {
  for index, span := range spans {
    if span.Cat != "WS" {
      return ambit, spans[index:]
    }
    ambit = ambit.SubtractLeft(span.Ambit)
  }
  return ambit, nil
}

func trimSpansRight(ambit *Ambit, spans []*spanT) (*Ambit, []*spanT) {
  for index := len(spans)-1; index >= 0; index-- {
    span := spans[index]
    if span.Cat != "WS" {
      return ambit, spans[:index+1]
    }
    ambit = ambit.SubtractRight(span.Ambit)
  }
  return ambit, nil
}
