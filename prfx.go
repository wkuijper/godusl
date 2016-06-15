package dusl

import (
  "io"
  "fmt" 
  "sort"
  "unicode/utf8"
)

type prfxTree struct {
  cat string
  children map[rune]*prfxTree
}

var _ Scanner = &prfxTree{}

func (this *prfxTree) lookup(lit string) string {
  return this.lookupRec([]byte(lit))
}

func (this *prfxTree) lookupRec(tail []byte) string {
  if this == nil {
    return ""
  } 
  if len(tail) == 0 {
    return this.cat
  }
  c, n := utf8.DecodeRune(tail)
  return this.children[c].lookupRec(tail[n:])
}

func (this *prfxTree) add(cat string, lit string) {
  this.addRec(cat, lit, []byte(lit))
}

func (this *prfxTree) addRec(cat string, lit string, tail []byte) {
  if len(tail) == 0 {
    this.cat = cat
  } else {
    c, n := utf8.DecodeRune(tail)
    rest := tail[n:]
    children := this.children
    if children == nil {
      children = make(map[rune]*prfxTree)
      this.children = children
    }
    sub, present := this.children[c]
    if !present {
      sub = &prfxTree{}
      children[c] = sub
    }
    sub.addRec(cat, lit, rest)
  }
}

func (this *prfxTree) dump(out io.Writer, prefix string) {
  if this.cat != "" {
    fmt.Fprintf(out, "%s [accept] \"%s\"\n", prefix, this.cat)
  } else {
    fmt.Fprintf(out, "%s ...\n", prefix)
  }
  if this.children != nil {
    list := make([]string, 0, len(this.children))
    for c, _ := range this.children {
      list = append(list, fmt.Sprintf("%c", c))
    }
    sort.Strings(list)
    for _, s := range list {
      r, _ := utf8.DecodeRuneInString(s)
      sub := this.children[r]
      sub.dump(out, prefix + s)
    }
  }
}

type prfxTreeScan struct {
  root *prfxTree
  state *prfxTree
}

func (this *prfxTree) Scan() Scan {
  return &prfxTreeScan{ root: this, state: this }
}

func (this *prfxTree) Report(categories map[string]bool) {
  if this.cat != "" {
    categories[this.cat] = true
  }
  for _, child := range this.children {
    child.Report(categories)
  }
}

func (this *prfxTreeScan) Consume(r rune) (string, bool) {
  state := this.state
  if state == nil {
    return "", false
  }
  children := state.children
  state = children[r]
  this.state = state
  if state == nil {
    return "", false
  }
  return state.cat, state.children != nil
}

func (this *prfxTreeScan) Reset() {
  this.state = this.root
}
