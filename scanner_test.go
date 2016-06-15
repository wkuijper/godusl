package dusl

import (
  "testing"
)

func TestScanner(t *testing.T) {
  ooc := PrefixScanner("OP += ===", "OB <%", "CB %>")
  osc := ooc.Scan()
  var cat string
  var cont bool
  cat, cont = osc.Consume('%')
  if cat != "" || !cont {
    t.Log("%", cat, cont)
    t.Fail()
  }
  cat, cont = osc.Consume('>')
  if cat != "CB" || cont {
    t.Log("%>", cat, cont)
    t.Fail()
  }
  cat, cont = osc.Consume('?')
  if cat != "" || cont {
    t.Log("%>?", cat, cont)
    t.Fail()
  }
  sss := SimpleStringScanner.Scan()
  cat, cont = sss.Consume('"')
  if cat != "" || !cont {
    t.Log(`"`, cat, cont)
    t.Fail()
  }
  cat, cont = sss.Consume('a')
  if cat != "" || !cont {
    t.Log(`"a`, cat, cont)
    t.Fail()
  }
  cat, cont = sss.Consume('b')
  if cat != "" || !cont {
    t.Log(`"ab`, cat, cont)
    t.Fail()
  }
  cat, cont = sss.Consume('c')
  if cat != "" || !cont {
    t.Log(`"abc`, cat, cont)
    t.Fail()
  }
  cat, cont = sss.Consume('"')
  if cat != "STR" || cont {
    t.Log(`"abc"`, cat, cont)
    t.Fail()
  }
  cat, cont = sss.Consume('?')
  if cat != "" || cont {
    t.Log(`"abc"]?`, cat, cont)
    t.Fail()
  }
}