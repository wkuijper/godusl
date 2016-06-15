package dusl

import (
	"bytes"
	"testing"
)

func TestUndent(t *testing.T) {
	text := `# title

    func f(x)
      # does stuff
      print "hello,"
            # more stuff
            "world!"
      read
    const x = 5
`
	source := &Source{Path: "tst", Text: []byte(text)}

	tree := Undent(source)

	buf := new(bytes.Buffer)
	tree.Dump(buf, "undent> ", false)
	res := buf.String()
	tgt := `undent> SQ:::tst[13:136]
undent>   SN:::tst[13:120]
undent>     UN:::tst[13:42]
undent>     SQ:::tst[48:120]
undent>       SN:::tst[48:109]
undent>         UN:::tst[48:109]
undent>         :::tst[109:109]
undent>       SQ:::tst[115:120]
undent>         SN:::tst[115:120]
undent>           UN:::tst[115:120]
undent>           :::tst[120:120]
undent>         :::tst[120:120]
undent>   SQ:::tst[124:136]
undent>     SN:::tst[124:136]
undent>       UN:::tst[124:136]
undent>       :::tst[136:136]
undent>     :::tst[136:136]
`
	if res != tgt {
		t.Log(res)
		t.Fail()
	}
}
