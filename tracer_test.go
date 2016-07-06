package dusl

import (
	"bytes"
	"testing"
)

func TestTracer(t *testing.T) {
	lang, err := NewSpec().
		Lexical(DefaultScanner).
		Category("ID", "identifier").
		Category("NUM", "number").
		Category("STR", "string").
		OperatorEFA("+", "-").
		OperatorBFA("+", "-").
		Brackets("( )").
		SequenceLabel("XSQ", "expression sequence").
		SentenceLabel("XSN", "expression sentence").
		Label("X", "expression").
		Literal("const").
		Grammar(`
      XSQ is>
        XSN
        XSQ
      or>
        <empty
        
      XSN is> X
     
      X is> (X) or> NUM or> ID or> const X
      X is> +X or> -X or> X + X or> X - X`)

	if err != nil {
		t.Log(err)
		t.Fail()
		return
	}

	tracer := lang.Tracer()
        
	text := []byte(`1 + +(5 --4) + x`)
	source := &Source{Path: "tst", Text: text}
        
	trace := tracer.TraceUndent(source, "XSQ")

	buf := new(bytes.Buffer)
	trace.Dump(buf, "> ", false)
	res := buf.String()
	tgt := `> XSQ:0:SQ::tst[0:16]
>   XSN:0:SN::tst[0:16]
>     X:6:OP:+:tst[0:16]
>       X:6:OP:+:tst[0:12]
>         X:1:NUM:1:tst[0:1]
>         X:4:OP:+:tst[4:12]
>           X:0:BB:( ):tst[5:12]
>             X:7:OP:-:tst[6:11]
>               X:1:NUM:5:tst[6:7]
>               X:5:OP:-:tst[9:11]
>                 X:1:NUM:4:tst[10:11]
>       X:2:ID:x:tst[15:16]
>   XSQ:1:::tst[16:16]
`
	if res != tgt {
		t.Log(res)
		t.Fail()
	}

	text = []byte(`
  1 + +(5 --4) + x
  1 + +(5 --4) + "5"
  1 + +(5 --4) + x
    hello`)

	source = &Source{Path: "tst", Text: text}

	trace = tracer.TraceUndent(source, "XSQ")

	buf = new(bytes.Buffer)
	trace.Dump(buf, "> ", false)
	res = buf.String()
	tgt = `> XSQ:0:SQ::tst[3:69]
>   XSN:0:SN::tst[3:20]
>     X:6:OP:+:tst[3:19]
>       X:6:OP:+:tst[3:15]
>         X:1:NUM:1:tst[3:4]
>         X:4:OP:+:tst[7:15]
>           X:0:BB:( ):tst[8:15]
>             X:7:OP:-:tst[9:14]
>               X:1:NUM:5:tst[9:10]
>               X:5:OP:-:tst[12:14]
>                 X:1:NUM:4:tst[13:14]
>       X:2:ID:x:tst[18:19]
>   XSQ:0:SQ::tst[22:69]
>     XSN:0:SN::tst[22:41]
>       X:6:OP:+:tst[22:40]
>         X:6:OP:+:tst[22:34]
>           X:1:NUM:1:tst[22:23]
>           X:4:OP:+:tst[26:34]
>             X:0:BB:( ):tst[27:34]
>               X:7:OP:-:tst[28:33]
>                 X:1:NUM:5:tst[28:29]
>                 X:5:OP:-:tst[31:33]
>                   X:1:NUM:4:tst[32:33]
>         ERR:0:STR:"5":tst[37:40]
>     XSQ:0:SQ::tst[43:69]
>       ERR:0:SN::tst[43:69]
>       XSQ:1:::tst[69:69]
`
	if res != tgt {
		t.Log(res)
		t.Fail()
	}
}
