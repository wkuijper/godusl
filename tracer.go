package dusl

import (
  "io"
  "fmt"
)

// A tracer represents a top down deterministic finite tree automaton that can be
// used to convert ambits or entire sources to traces.
type Tracer interface {
  Trace(ambit *Ambit, lbl string) *Trace
  TraceUndent(source *Source, lbl string) *Trace
  Dump(out io.Writer, prfx string)
}

type tracer struct {
  sparser Sparser
  templates map[string][]*templateT
  descriptions map[string]string
}

// A node in an acceptance trace of a top down deterministic finite tree automaton.
// The Lbl field containsthe label for the current tree automaton state,
// the Idx field contains the index of the transition rule that was applied,
// the Syn field refers to the current node of the syntax tree over which the tree
// automaton was run, the Subs field contains the subtraces in order of a
// left-to-right traversal of the transition rule template.
// If there was no rule that could be applied the Lbl will be "ERR" and the Err field
// will contain a descriptive error message.
type Trace struct {
  Lbl string
  Idx int
  Syn *Syntax
  Err string
  Subs []*Trace
}

type templateT struct {
  lbl string
  subCount int
  matchCat bool
  cat string
  matchLit bool
  lit string
  litSet map[string]bool
  left *templateT
  right *templateT
}

type waitingItemT struct {
  node *Syntax
  trace *Trace
}

type waitingT struct {
  list []waitingItemT
}

func (this *Trace) DumpToString(pretty bool) string {
  return dumpToString(this, pretty)
}

func (this *Trace) Dump(out io.Writer, prfx string, pretty bool) {
  if pretty {
    this.dumpPretty(out, prfx)
  } else {
    this.dumpRaw(out, prfx)
  }
}

func (this *Trace) dumpPretty(out io.Writer, prfx string) {
  if this == nil {
    return
  }
  if this.Err != "" {
    fmt.Fprintf(out, "%s%s:%d:%s\n", prfx, this.Lbl, this.Idx, this.Err)
  } else {
    fmt.Fprintf(out, "%s%s:%d:%s\n", prfx, this.Lbl, this.Idx, this.Syn.Lit)
  }
  subs := this.Subs
  if len(subs) > 0 {
    subPrfx := prfx + "  "
    for _, sub := range subs {
      sub.dumpPretty(out, subPrfx)
    }
  }
}

func (this *Trace) dumpRaw(out io.Writer, prfx string) {
  if this == nil {
    return
  }
  fmt.Fprintf(out, "%s%s:%d:%s:%s:%s\n", prfx, this.Lbl, this.Idx, this.Syn.Cat, this.Syn.Lit, this.Syn.Ambit)
  subs := this.Subs
  if len(subs) > 0 {
    subPrfx := prfx + "  "
    for _, sub := range subs {
      sub.dumpRaw(out, subPrfx)
    }
  }
}

func (this *Trace) ErrorN(n int) error {
  return SummaryError(this.errors(), n)
}

func (this *Trace) errors() []error {
  if this == nil {
    return nil
  }
  errs := this.Syn.errors()
  return this.gatherErrors(errs)
}

func (this *Trace) gatherErrors(errs []error) []error {
  if this == nil {
    return errs
  }
  if this.Lbl == "ERR" {
    return append(errs, AmbitError(this.Syn.Ambit, this.Err))
  }
  for _, sub := range this.Subs {
    errs = sub.gatherErrors(errs)
  }
  return errs
}

func (this *templateT) subCountOrZero() int {
  if this == nil {
    return 0
  }
  return this.subCount
}
func (this *templateT) dump(out io.Writer, prfx string) {
  if this == nil {
    return
  }
  fmt.Fprintf(out, "%s%s:%s:%s:%d\n", prfx, this.lbl, this.cat, this.lit, this.subCount)
  this.left.dump(out, prfx+"  ")
  this.right.dump(out, prfx+"  ")
}

func newTracer(sparser Sparser, templates map[string][]*templateT, descriptions map[string]string) Tracer {
  return &tracer{ sparser: sparser, templates: templates, descriptions: descriptions }
}

func (this *tracer) Trace(ambit *Ambit, lbl string) *Trace {
  root := this.sparser.Sparse(ambit)
  return this.label(root, lbl)
}

func (this *tracer) TraceUndent(source *Source, lbl string) *Trace {
  root := this.sparser.SparseUndent(source)
  return this.label(root, lbl)
}

func (this *tracer) label(root *Syntax, lbl string) *Trace {
  start := &Trace{ Lbl: lbl, Syn: root }
  waiting := &waitingT{ list: []waitingItemT{ waitingItemT{ node: root, trace: start } } }
  
  for {
    l := len(waiting.list)-1
    if l < 0 {
      break
    }
    item := waiting.list[l]
    waiting.list = waiting.list[:l]

    node, trace := item.node, item.trace
    
    templates := this.templates[trace.Lbl]
    matched := false
    for idx, template := range templates {
      if template.checkMatch(node) {
        trace.Idx = idx
        trace.Subs = make([]*Trace, template.subCount)
        template.performMatch(node, waiting, 0, trace.Subs)
        matched = true
        break
      }
    }
    if !matched {
      trace.Err = fmt.Sprintf("expected: %s", this.descriptions[trace.Lbl])
      trace.Lbl = "ERR"
    }
  }

  return start
}

func (this *templateT) checkMatch(node *Syntax) bool {
  if node == nil ||
       (this.matchCat && this.cat != node.Cat) ||
       (this.matchLit &&
         ((this.litSet != nil && !this.litSet[node.Lit]) || (this.litSet == nil && this.lit != node.Lit))) {
    return false
  }
  if this.left != nil { // implies this.right != nil
    return this.left.checkMatch(node.Left) && this.right.checkMatch(node.Right)
  }
  return true
}

func (this *templateT) performMatch(node *Syntax, waiting *waitingT, subi int, subs []*Trace) int {
  lbl := this.lbl
  if lbl != "" {
    sub := &Trace{ Lbl: lbl, Syn: node }
    subs[subi] = sub
    subi++
    waiting.list = append(waiting.list, waitingItemT{ node: node, trace: sub})
    return subi
  }
  if this.left != nil { // implies this.right != nil
    subi = this.left.performMatch(node.Left, waiting, subi, subs)
    subi = this.right.performMatch(node.Right, waiting, subi, subs)
  }
  return subi
}

func (this *tracer) Dump(out io.Writer, prfx string) {
  for lbl, templates := range this.templates {
    prfx2 := fmt.Sprintf("%s%s> ", prfx, lbl)
    for index, template := range templates {
      prfx3 := fmt.Sprintf("%s%d> ", prfx2, index)
      template.dump(out, prfx3)
    }
  }
}
