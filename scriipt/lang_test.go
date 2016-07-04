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
  r := tracer.TraceUndent(dusl.SourceFromString(s), "Stmt").DumpToString(true)
  if r != `Stmt:1:
  loopHeader:1:;
    simpleStmt:1:=
      ERR:0:expected: variable name
      expr:1:0
    expr:7:<
      expr:0:x
      expr:1:20
    simpleStmt:2:++
      i:0:x
  Stmt:3:
    expr:7:>
      expr:0:x
      expr:1:10
    Stmt:4:
      simpleStmt:4:
        f:0:print
        args:1:x
          aargs:1:x
            expr:0:x
      Stmt:0:
    ElseCont:0:
      expr:7:<
        expr:0:x
        expr:1:0
      Stmt:4:
        simpleStmt:4:
          f:0:printf
          args:1:,
            aargs:0:,
              aargs:1:-
                expr:5:-
                  expr:0:x
              expr:0:x
        Stmt:0:
      ElseCont:1:
        Stmt:4:
          simpleStmt:4:
            f:0:print
            args:0:
          Stmt:0:
        Stmt:0:
  Stmt:0:
` {
    t.Log(r)
    t.Fail()
    return
  }
}
