package dusl

import (
  "unicode/utf8"
  "strings"
  "fmt"
)

type Token struct {
  Cat string
  Lit string
  Err string
  Ambit *Ambit
}

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

func (token *Token) String() string {
  if strings.TrimSpace(token.Lit) == "" { 
    return token.Cat // <-- "WS"
  } else {
    return token.Cat + ":" + token.Lit
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

func (this *tokenizer) TokenizeUndent(src *Source) *Syntax {
  return Undent(src).mapUnparsedAmbits(func(a *Ambit)string { return fmt.Sprintf("%v", this.Tokenize(a)) })
}