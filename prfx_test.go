package dusl

import (
  "testing"
  "bytes"
)

func TestPrfxTree(t *testing.T) {
  pt := &prfxTree{}
  pt.add("OP", "=")
  pt.add("OP", "==")
  pt.add("OP", "+=")
  pt.add("OP", "++")
  pt.add("OP", "===")
  pt.add("OP", "const")
  pt.add("OP", "<")
  pt.add("OB", "<%")
  buf := new(bytes.Buffer)
  pt.dump(buf, "pt> ")
  tgt := `pt>  ...
pt> + ...
pt> ++ [accept] "OP"
pt> += [accept] "OP"
pt> < [accept] "OP"
pt> <% [accept] "OB"
pt> = [accept] "OP"
pt> == [accept] "OP"
pt> === [accept] "OP"
pt> c ...
pt> co ...
pt> con ...
pt> cons ...
pt> const [accept] "OP"
`
  if buf.String() != tgt {
    t.Log(buf.String())
    t.Fail()
  }
}