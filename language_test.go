package dusl

import (
  "testing"
  "bytes"
)

func TestLanguage(t *testing.T) {
  lang, err := NewSpec().
    Lexical(DefaultScanner).
    Category("NUM", "number").
    Category("ID", "identifier").
    Category("STR", "string").
    OperatorEFA("+", "-", "*", "&", "!", "^").
    OperatorBFA("*", "/", "%", "<<", ">>", "&", "&^").
    OperatorBFA("+", "-", "|", "^").
    OperatorBFA("==", "!=", "<", "<=", ">=", ">").
    OperatorBFA("&&").
    OperatorBFA("||").
    OperatorBFA("..").
    OperatorBFA(",").
    OperatorBFA(";").
    OperatorEFA("range").
    OperatorBFA("=", ":=").
    OperatorAFE("++", "--").
    OperatorEFA("if", "else", "for", "func").
    OperatorEFE("else").
    Brackets("( )", "[ ]", "{ }").
    
    Label("x", "expression").
    Label("xs", "one or more comma separated expressions").
    Label("args", "zero or more comma separated expressions").
    
    Label("s", "simple statement").

    Label("f", "for clause").

    SequenceLabel("S", "statement").
    SequenceLabel("E", "else continuation statement").

    ShorthandOperator("prfx", "-", "+", "*", "&", "!", "^").
    ShorthandOperator("infix", "prfx", "/", "%", "<<", ">>", "&^", "|",
                               "==", "!=", "<", "<=", ">=", ">", "&&", "||").
    
    Grammar(`
      x is> ID or> NUM or> STR or> (x) or> x(args) or> prfx x or> x infix x

      xs is> x, xs or> x
      args is> <empty or> xs
      
      s is> x := x or> x = x or> x++ or> x-- or> x

      S is>
        for f
          S
        S
      S is>
        if x
          S
        E
      S is>
        s
      or>
        <empty

      E is>
        else if x
          S
        E
      or>
        else
          S
        S
      or>
        S

      f is> x := range x
      f is> s; x; s
      f is> x
      f is> <empty`);

  if err != nil {
    t.Log(err)
    t.Fail()
    return
  }

  tracer := lang.Tracer()
  
  src := `
    for x != 0
      if x > 3
        print(x)
      else if x < 0
        print(x, -x)
      else
        print()`

  source := &Source{ Path: "src", LineOffset: 92, Text: []byte(src) }
  
  trace := tracer.TraceUndent(source, "S")

  buf := new(bytes.Buffer)
  trace.Dump(buf, "> ", false)
  res := buf.String()
  tgt := `> S:0:SQ::src[5:115]
>   f:2:OP:!=:src[9:15]
>     x:6:OP:!=:src[9:15]
>       x:0:ID:x:src[9:10]
>       x:1:NUM:0:src[14:15]
>   S:1:SQ::src[22:115]
>     x:6:OP:>:src[25:30]
>       x:0:ID:x:src[25:26]
>       x:1:NUM:3:src[29:30]
>     S:2:SQ::src[39:48]
>       s:4:GLUE::src[39:47]
>         x:4:GLUE::src[39:47]
>           x:0:ID:print:src[39:44]
>           args:1:ID:x:src[45:46]
>             xs:1:ID:x:src[45:46]
>               x:0:ID:x:src[45:46]
>     E:0:SQ::src[54:115]
>       x:6:OP:<:src[62:67]
>         x:0:ID:x:src[62:63]
>         x:1:NUM:0:src[66:67]
>       S:2:SQ::src[76:89]
>         s:4:GLUE::src[76:88]
>           x:4:GLUE::src[76:88]
>             x:0:ID:print:src[76:81]
>             args:1:OP:,:src[82:87]
>               xs:0:OP:,:src[82:87]
>                 x:0:ID:x:src[82:83]
>                 xs:1:OP:-:src[85:87]
>                   x:5:OP:-:src[85:87]
>                     x:0:ID:x:src[86:87]
>       E:1:SQ::src[95:115]
>         S:2:SQ::src[108:115]
>           s:4:GLUE::src[108:115]
>             x:4:GLUE::src[108:115]
>               x:0:ID:print:src[108:113]
>               args:0:::src[114:114]
>         S:3:::src[115:115]
>   S:3:::src[115:115]
`
  
  if res != tgt {
    t.Log(res)
    t.Fail()
    return
  }
}