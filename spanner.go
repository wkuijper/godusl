package dusl

import (
  "fmt"
  "strings"
)

type Spanner interface {
  Span(ambit *Ambit) []*Span
  SpanUndent(src *Source) *Syntax
}

type Span struct {
  Cat string
  Lit string
  Err string
  Ambit *Ambit
  SubAmbit *Ambit
  Children []*Span
}

type spanner struct {
  tokenizer Tokenizer
  precedenceB map[string]int
}

func newSpanner(tokenizer Tokenizer, precedenceB map[string]int) Spanner {
  return &spanner{ tokenizer: tokenizer, precedenceB: precedenceB }
}

func (this *Span) String() string {
  if len(this.Children) == 0 {
    if strings.TrimSpace(this.Lit) == "" { 
      return this.Cat // <-- "WS"
    } else {
      return this.Cat + ":" + this.Lit
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

func (this *spanner) Span(ambit *Ambit) []*Span {
  tokens := this.tokenizer.Tokenize(ambit)
  var spans []*Span
  for len(tokens) > 0 {
    spans, tokens = this.span(spans, tokens)
    if len(tokens) > 0 {
      // stray closing
      var token *Token
      token, tokens = tokens[0], tokens[1:]
      spans = append(spans, &Span{ Cat: "ERR", Err: fmt.Sprintf("unexpected closing bracket: '%s'", token.Ambit.ToString()), Ambit: token.Ambit })
    }
  }
  return spans
}

func (this *spanner) span(spans []*Span, tokens []*Token) ([]*Span, []*Token) {
  if spans == nil && len(tokens) > 0 {
    spans = make([]*Span, 0, min(32, len(tokens)))
  }
  for len(tokens) > 0 {
    token := tokens[0]
    if token.Cat == "CB" {
      break
    }
    tokens = tokens[1:]
    var span *Span
    if token.Cat == "OB" {
      opbr := token
      var children []*Span
      children, tokens = this.span(children, tokens)
      if len(tokens) < 1 {
        // missing closing
        if len(children) < 1 {
          span = &Span{ Cat: "ERR", Err: fmt.Sprintf("missing closing bracket: corresponding to opening bracket: '%s'", opbr.Ambit.ToString()), Ambit: opbr.Ambit }
        } else {
          span = &Span{ Cat: "ERR", Err: fmt.Sprintf("missing closing bracket: corresponding to opening bracket: '%s'", opbr.Ambit.ToString()), Ambit: opbr.Ambit.Merge(children[len(children)-1].Ambit) }
        }
      } else {
        var clbr *Token
        clbr, tokens = tokens[0], tokens[1:]
        brcat := opbr.Lit + " " + clbr.Lit
        if clbr.Cat != "CB" || this.precedenceB[brcat] < 1 {
          span = &Span{ Cat: "ERR", Err: fmt.Sprintf("non-matching brackets: '%s'", brcat), Ambit: opbr.Ambit.Merge(clbr.Ambit) }
        } else {
          span = &Span{ Cat: "BB", Lit: brcat, Ambit: opbr.Ambit.Merge(clbr.Ambit), SubAmbit: opbr.Ambit.Merge(clbr.Ambit).SubtractLeft(opbr.Ambit).SubtractRight(clbr.Ambit), Children: children }
        }
      }
    } else {
      span = &Span{ Cat: token.Cat, Lit: token.Lit, Err: token.Err, Ambit: token.Ambit }
    }
    spans = append(spans, span)
  }
  return spans, tokens
}

func (this *spanner) SpanUndent(src *Source) *Syntax {
  return Undent(src).mapUnparsedAmbits(func(a *Ambit)string { return fmt.Sprintf("%v", this.Span(a)) })
}