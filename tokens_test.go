package dusl

import (
  "testing"
  "fmt"
)

func TestTokens(t *testing.T) {
  text := []byte(`a += 1 "hello(), \n"()[world], <% ***%> "===`)
  source := &Source{ Path: "string", Text: text }
  ambit := source.FullAmbit()
  scanner :=
    &seqScanner{PrefixScanner("OP +=", "OP ===", "OB ( [ { <%", "CB } ] ) %>"),
                NewDefaultScanner()}

  tokenizer := newTokenizer(scanner)
  tokens := tokenizer.Tokenize(ambit)
  res := fmt.Sprintf("%s", tokens)
  tgt := `[ID:a WS OP:+= WS NUM:1 WS STR:"hello(), \n" OB:( CB:) OB:[ ID:world CB:] ERR:unexpected character(s): ',' WS OB:<% WS ERR:unexpected character(s): '***' CB:%> WS ERR:unexpected character(s): '"' OP:===]`
  if res != tgt {
    t.Log(res)
    t.Fail()
  }

}