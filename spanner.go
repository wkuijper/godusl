package dusl

import (
  "fmt"
  "strings"
)

// A Spanner is used to convert ambits into sequences of Spans using the Span method,
// or an entire source to a tree of formatted lists of Spans using the SpanUndent
// method. The latter method is mainly intended for unit tests and debugging.
type spannerI interface {
  span(ambit *Ambit) []*spanT
  spanUndent(src *Source) *Syntax
}

// A Span represents either a Token or a sequence of subspans based on explicit
// grouping with bracket tokens.  The Lit field is assigned
// the text of the ambit as a literal string or the brackets in case of an explicit grouping.
// The Err field is set with a descriptive error message iff the Cat field equals the special
// error category "ERR".
// The Children field is set only for an explicit grouping with a slice of child spans.
type spanT struct {
  Cat string
  Lit string
  Err string
  Ambit *Ambit
  SubAmbit *Ambit
  Children []*spanT
}

type spanner struct {
  tokenizer Tokenizer
  precedenceB map[string]int
}

func newSpanner(tokenizer Tokenizer, precedenceB map[string]int) spannerI {
  return &spanner{ tokenizer: tokenizer, precedenceB: precedenceB }
}

func (this *spanT) String() string {
  cat := this.Cat
  lit := this.Lit
  if cat == "ERR" {
    lit = this.Err
  }
  if len(this.Children) == 0 {
    if strings.TrimSpace(lit) == "" { 
      return cat // <-- "WS"
    } else {
      return cat + ":" + lit
    }
  }
  sublist := fmt.Sprintf("%s", this.Children)
  parts := strings.Split(this.Lit, " ")
  if len(sublist) >= 2 && len(parts) >= 2 { 
    return parts[0] + sublist[1:len(sublist)-1] + parts[1]
  }
  // defensive
  return this.Cat + ":" + sublist
}

// Span returns the slice of Spans obtained by scanning the given source ambit.
func (this *spanner) span(ambit *Ambit) []*spanT {
  tokens := this.tokenizer.Tokenize(ambit)
  var spans []*spanT
  for len(tokens) > 0 {
    spans, tokens = this.span2(spans, tokens)
    if len(tokens) > 0 {
      // stray closing
      var token *Token
      token, tokens = tokens[0], tokens[1:]
      spans = append(spans, &spanT{ Cat: "ERR", Err: fmt.Sprintf("unexpected closing bracket: '%s'", token.Ambit.ToString()), Ambit: token.Ambit })
    }
  }
  return spans
}

func (this *spanner) span2(spans []*spanT, tokens []*Token) ([]*spanT, []*Token) {
  if spans == nil && len(tokens) > 0 {
    spans = make([]*spanT, 0, min(32, len(tokens)))
  }
  for len(tokens) > 0 {
    token := tokens[0]
    if token.Cat == "CB" {
      break
    }
    tokens = tokens[1:]
    var span *spanT
    if token.Cat == "OB" {
      opbr := token
      var children []*spanT
      children, tokens = this.span2(children, tokens)
      if len(tokens) < 1 {
        // missing closing
        if len(children) < 1 {
          span = &spanT{ Cat: "ERR", Err: fmt.Sprintf("missing closing bracket: corresponding to opening bracket: '%s'", opbr.Ambit.ToString()), Ambit: opbr.Ambit }
        } else {
          span = &spanT{ Cat: "ERR", Err: fmt.Sprintf("missing closing bracket: corresponding to opening bracket: '%s'", opbr.Ambit.ToString()), Ambit: opbr.Ambit.Merge(children[len(children)-1].Ambit) }
        }
      } else {
        var clbr *Token
        clbr, tokens = tokens[0], tokens[1:]
        brcat := opbr.Lit + " " + clbr.Lit
        if clbr.Cat != "CB" || this.precedenceB[brcat] < 1 {
          span = &spanT{ Cat: "ERR", Err: fmt.Sprintf("non-matching brackets: '%s'", brcat), Ambit: opbr.Ambit.Merge(clbr.Ambit) }
        } else {
          span = &spanT{ Cat: "BB", Lit: brcat, Ambit: opbr.Ambit.Merge(clbr.Ambit), SubAmbit: opbr.Ambit.Merge(clbr.Ambit).SubtractLeft(opbr.Ambit).SubtractRight(clbr.Ambit), Children: children }
        }
      }
    } else {
      span = &spanT{ Cat: token.Cat, Lit: token.Lit, Err: token.Err, Ambit: token.Ambit }
    }
    spans = append(spans, span)
  }
  return spans, tokens
}

// spanUndent returns the tree of formatted span lists obtained by first undenting
// and then spanning the given source.
func (this *spanner) spanUndent(src *Source) *Syntax {
  return Undent(src).mapUnparsedAmbits(func(a *Ambit)string { return fmt.Sprintf("%v", this.span(a)) })
}