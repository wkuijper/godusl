package dusl

import (
	"bytes"
	"testing"
)

func TestSparser(t *testing.T) {

	lang, err := NewSpec().
		Lexical(NewDefaultScanner()).
		OperatorEFA("+", "-").
		OperatorBFA("+", "-").
		Brackets("( )").
		Grammar("")

	if err != nil {
		t.Log(err)
		t.Fail()
		return
	}

	sparser := lang.Sparser()

	text := []byte(`1 + +(5 --4) + x`)
	source := &Source{Path: "tst", Text: text}
	ambit := source.FullAmbit()

	tree := sparser.Sparse(ambit)

	buf := new(bytes.Buffer)
	tree.Dump(buf, "> ", false)
	res := buf.String()
	tgt := `> OP:+::tst[0:16]
>   OP:+::tst[0:12]
>     NUM:1::tst[0:1]
>     OP:+::tst[4:12]
>       :::tst[4:4]
>       BB:( )::tst[5:12]
>         OP:-::tst[6:11]
>           NUM:5::tst[6:7]
>           OP:-::tst[9:11]
>             :::tst[9:9]
>             NUM:4::tst[10:11]
>         :::tst[12:12]
>   ID:x::tst[15:16]
`
	if res != tgt {
		t.Log(res)
		t.Fail()
	}
}

func TestSparserUndent(t *testing.T) {

	lang, err := NewSpec().
		Lexical(NewDefaultScanner()).
		OperatorEFA("+", "-").
		OperatorBFA("+", "-").
		Brackets("( )").
		Grammar("")

	if err != nil {
		t.Log(err)
		t.Fail()
		return
	}

	sparser := lang.Sparser()

	text := []byte(`
func (1 + 2)
  a + b
    c + d
`)
	source := &Source{Path: "tst", Text: text}

	tree := sparser.SparseUndent(source)

	buf := new(bytes.Buffer)
	tree.Dump(buf, "> ", false)
	res := buf.String()
	tgt := `> SQ:::tst[1:32]
>   SN:::tst[1:32]
>     JUXT: ::tst[1:13]
>       ID:func::tst[1:5]
>       BB:( )::tst[6:13]
>         OP:+::tst[7:12]
>           NUM:1::tst[7:8]
>           NUM:2::tst[11:12]
>         :::tst[13:13]
>     SQ:::tst[16:32]
>       SN:::tst[16:32]
>         OP:+::tst[16:21]
>           ID:a::tst[16:17]
>           ID:b::tst[20:21]
>         SQ:::tst[26:32]
>           SN:::tst[26:32]
>             OP:+::tst[26:31]
>               ID:c::tst[26:27]
>               ID:d::tst[30:31]
>             :::tst[32:32]
>           :::tst[32:32]
>       :::tst[32:32]
>   :::tst[32:32]
`
	if res != tgt {
		t.Log(res)
		t.Fail()
	}
}
