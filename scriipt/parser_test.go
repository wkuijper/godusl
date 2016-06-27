package scriipt

import (
  "testing"
  "dusl"
)

func TestParser(t *testing.T) {
  s := `
    for x = 0; x < 20; x++
      if x > 10
        print(x)
      else if x < 0
        printf(-x, x)
      else
        print()`
  tree, err := Parse(dusl.SourceFromString(s))
  if err != nil {
    t.Log(err.Error())
    t.Fail()
    return
  }
  r := dumpToString(tree)
  if r != `for x = 0; (x < 20); x++
  if (x > 10)
    print(x)
  else
    if (x < 0)
      printf(-x, x)
    else
      print()
` {
    t.Log(r)
    t.Fail()
    return
  }
}