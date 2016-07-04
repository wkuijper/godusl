package scriipt

import (
  "testing"
  "dusl"
  "fmt"
)

func TestLang(t *testing.T) {
  testTokenizer(Lang.Tokenizer(), t)
  testSparser(Lang.Sparser(), t)
  testTracer(Lang.Tracer(), t)
}

func testTokenizer(tokenizer dusl.Tokenizer, t *testing.T) {
  r := fmt.Sprintf("%v", tokenizer.Tokenize(dusl.AmbitFromString("a + b { -1 if [then] else() }")))
  if r != "[ID:a WS OP:+ WS ID:b WS OB:{ WS OP:- NUM:1 WS ID:if WS OB:[ ID:then CB:] WS ID:else OB:( CB:) WS CB:}]" {
    t.Log(r)
    t.Fail()
    return
  }
}

func testSparser(sparser dusl.Sparser, t *testing.T) {
  r := sparser.Sparse(dusl.AmbitFromString("a + b { -1 if [then] else() }")).DumpToString(true)
  if r != `OP:+
  ID:a
  JUXT: 
    ID:b
    BB:{ }
      OP:-
        :
        JUXT: 
          NUM:1
          JUXT: 
            ID:if
            JUXT: 
              BB:[ ]
                ID:then
                :
              GLUE:
                ID:else
                BB:( )
                  :
                  :
      :
` {
    t.Log(r)
    return
  }
}

func testTracer(tracer dusl.Tracer, t *testing.T) {
  s := `
    for x := 0; x < 20; x++
      if x > 10
        print(x)
      else if x < 0
        printf(-x, x)
      else
        print()`
  r := tracer.TraceUndent(dusl.SourceFromString(s), "S").DumpToString(true)
  if r != `S:1:
  hdr:1:;
    s:1:=
      ERR:0:expected: variable name
      x:1:0
    x:7:<
      x:0:x
      x:1:20
    s:2:++
      i:0:x
  S:2:
    x:7:>
      x:0:x
      x:1:10
    S:3:
      s:4:
        f:0:print
        args:1:x
          aargs:1:x
            x:0:x
      S:0:
    E:0:
      x:7:<
        x:0:x
        x:1:0
      S:3:
        s:4:
          f:0:printf
          args:1:,
            aargs:0:,
              aargs:1:-
                x:5:-
                  x:0:x
              x:0:x
        S:0:
      E:1:
        S:3:
          s:4:
            f:0:print
            args:0:
          S:0:
        S:0:
  S:0:
` {
    t.Log(r)
    t.Fail()
    return
  }
}
