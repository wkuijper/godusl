package dusl

import (
  "testing"
  "fmt"
)

func TestSpanner(t *testing.T) {
  text := []byte(`(a [b {c] } d += )](`)
  source := &Source{ Path: "string", Text: text }
  ambit := source.FullAmbit()
  scanner :=
    &seqScanner{master: PrefixScanner("OP += ===", "OB <% ( [ {", "CB } ] ) %>"),
                slave: DefaultScanner}

  spanner := newSpanner(newTokenizer(scanner), map[string]int{ "<% %>":1, "( )":1, "[ }":1 })
  spans := spanner.span(ambit)
  if len(spans) != 3 {
    t.Log("len(spans) ==", len(spans))
    t.Fail()
  }
  res := fmt.Sprintf("%v", spans)
  tgt := `[(ID:a WS [ID:b WS ERR:non-matching brackets: '{ ]' WS} WS ID:d WS OP:+= WS) ERR:unexpected closing bracket: ']' ERR:missing closing bracket: corresponding to opening bracket: '(']`
  if res != tgt {
    t.Log(res)
    t.Fail()
  }
}