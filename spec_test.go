package dusl

import (
  "testing"
)

func TestSpec(t *testing.T) {
  lang, err := NewSpec().
    Lexical(DefaultScanner).
    Category("ID", "identifier").
    Category("NUM", "number").
    Category("STR", "string").
    OperatorEFA("+", "-").
    OperatorBFA("+", "-").
    Brackets("( )").
    SequenceLabel("XSQ", "expression sequence").
    SentenceLabel("XSE", "expression sentence").
    Label("X", "expression").
    Grammar(`

      XSQ is>
        XSE
        XSQ
      or>
        <empty

      XSE is> X
      
      X is> (X) or> NUM or> ID or> +X or> -X or> X + X or> X - X`);
      
  if err != nil {
    t.Log(err)
    t.Fail()
    return
  }

  t.Log(lang)
}