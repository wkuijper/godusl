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
                slave: NewDefaultScanner()}

  spanner := newSpanner(newTokenizer(scanner), map[string]int{ "<% %>":1, "( )":1, "[ }":1 })
  spans := spanner.Span(ambit)
  if len(spans) != 3 {
    t.Log("len(spans) ==", len(spans))
    t.Fail()
  }
  res := fmt.Sprintf("%v", spans)
  tgt := `[(ID:a WS [ID:b WS ERR WS} WS ID:d WS OP:+= WS) ERR ERR]`
  if res != tgt {
    t.Log(res)
    t.Fail()
  }
}