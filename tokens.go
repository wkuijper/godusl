package dusl

import (
  "unicode/utf8"
  "strings"
  "fmt"
)

// A Token represents an ambit that has been identified as lexically relevant
// and that has been assigned a lexical category. The Lit field is assigned
// the text of the ambit as a literal string. The Err field is set with a
// descriptive error message iff the Cat field equals the special error
// category "ERR".
type Token struct {
  Cat string
  Lit string
  Err string
  Ambit *Ambit
}

// A Tokenizer is used for converting an Ambit to a list of Tokens using the Tokenize
// method or an entire Source to a tree of formatted lists of Tokens using the
// TokenizeUndent method. The latter method is mainly intended for unit tests and
// debugging.
type Tokenizer interface {
  Tokenize(ambit *Ambit) []*Token
  TokenizeUndent(src *Source) *Syntax
}

type tokenizer struct {
  scan Scan
}

func newTokenizer(scanner Scanner) Tokenizer {
 return &tokenizer{ scan: scanner.Scan() }
}

func (this *Token) String() string {
  cat := this.Cat
  lit := this.Lit
  if cat == "ERR" {
    lit = this.Err
  }
  if strings.TrimSpace(lit) == "" { 
    return cat // <-- "WS"
  } else {
    return cat + ":" + lit
  }
}

func (this *tokenizer) splitOnToken(ambit *Ambit) (string, *Ambit, *Ambit) {
  text := ambit.Source.Text
  end := ambit.End
  scan := this.scan
  scan.Reset()
  
  i := ambit.Start
  tokenEnd := i
  tokenCat := ""
  
  for i < end {
    c := text[i]
    var r rune
    if c <= 127 {
      r = rune(c)
      i++
    } else {
      var n int
      r, n = utf8.DecodeRune(text[i:])
      i += n
    }
    cat, cont := scan.Consume(r)
    if cat != "" {
      tokenCat = cat
      tokenEnd = i
    }
    if !cont {
      break
    }
  }
  tokenAmbit, restAmbit := ambit.SplitAtAbs(tokenEnd) 
  return tokenCat, tokenAmbit, restAmbit
}

func (this *tokenizer) splitOnError(ambit *Ambit) (*Ambit, *Ambit) {
  text := ambit.Source.Text
  end := ambit.End
  i := ambit.Start
  var errAmbit *Ambit
  var restAmbit = ambit
  for i < end {
    cat, _, _ := this.splitOnToken(restAmbit)
    if cat != "" {
      break
    }
    c := text[i]
    if c <= 127 {
      i++
    } else {
      var n int
      _, n = utf8.DecodeRune(text[i:])
      i += n
    }
    errAmbit, restAmbit = ambit.SplitAtAbs(i)
  }
  return errAmbit, restAmbit
}

// Tokenize returns the slice of Tokens obtained by scanning the given source ambit.
func (this *tokenizer) Tokenize(ambit *Ambit) []*Token {
  tokens := make([]*Token, 0, 32)
  for !ambit.IsEmpty() {
    tokenCat, tokenAmbit, restAmbit := this.splitOnToken(ambit)
    var token *Token
    if tokenCat == "" {
      tokenAmbit, restAmbit = this.splitOnError(ambit)
      tokenCat = "ERR"
      token = &Token{ Cat: tokenCat, Lit: tokenAmbit.ToString(), Err: fmt.Sprintf("unexpected character(s): '%s'", tokenAmbit.ToString()), Ambit: tokenAmbit }
    } else {
      token = &Token{ Cat: tokenCat, Lit: tokenAmbit.ToString(), Ambit: tokenAmbit }
    }
    tokens = append(tokens, token)
    ambit = restAmbit
  }
  return tokens
}

// Tokenize returns the tree of formatted token lists obtained by first undenting
// and then scanning the given source.
func (this *tokenizer) TokenizeUndent(src *Source) *Syntax {
  return Undent(src).mapUnparsedAmbits(func(a *Ambit)string { return fmt.Sprintf("%v", this.Tokenize(a)) })
}