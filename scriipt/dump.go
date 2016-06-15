package scriipt

import (
  "io"
  "fmt"
)

func (this *idExpr) Dump(out io.Writer, prfx string, pretty bool) {
  if pretty {
    fmt.Fprintf(out, "%s", this.idName)
  } else {
    fmt.Fprintf(out, "%sIdExpr:%s:%s\n", prfx, this.idName, this.ambit)
  }
}

func (this *numExpr) Dump(out io.Writer, prfx string, pretty bool) {
  if pretty {
    fmt.Fprintf(out, "%d", this.num)
  } else {
    fmt.Fprintf(out, "%sNumExpr:%d:%s\n", prfx, this.num, this.ambit)
  }
}

func (this *strExpr) Dump(out io.Writer, prfx string, pretty bool) {
  if pretty {
    fmt.Fprintf(out, "%s", this.EscapedStr())
  } else {
    fmt.Fprintf(out, "%sStrExpr:%s:%s\n", prfx, this.EscapedStr(), this.ambit)
  }
}

func (this *callExpr) Dump(out io.Writer, prfx string, pretty bool) {
  if pretty {
    fmt.Fprintf(out, "%s(", this.funcName)
    for index, arg := range this.args {
      if index > 0 {
        fmt.Fprintf(out, ", ")
      }
      arg.Dump(out, "", true)
    }
    fmt.Fprintf(out, ")")
  } else {
    fmt.Fprintf(out, "%sCallExpr:%s:%s\n", prfx, this.funcName, this.ambit)
    for _, arg := range this.args {
      arg.Dump(out, prfx + "  ", false)
    }
  }
}

func (this *prefixExpr) Dump(out io.Writer, prfx string, pretty bool) {
  if pretty {
    fmt.Fprintf(out, "%s", this.lit)
    this.sub.Dump(out, "", true)
  } else {
    fmt.Fprintf(out, "%sPrefixExpr:%s:%s\n", prfx, this.lit, this.ambit)
    this.sub.Dump(out, prfx + "  ", false)
  }
}

func (this *postfixExpr) Dump(out io.Writer, prfx string, pretty bool) {
  if pretty {
    this.sub.Dump(out, "", true)
    fmt.Fprintf(out, "%s", this.lit)
  } else {
    fmt.Fprintf(out, "%sPostfixExpr:%s:%s\n", prfx, this.lit, this.ambit)
    this.sub.Dump(out, prfx + "  ", false)
  }
}

func (this *infixExpr) Dump(out io.Writer, prfx string, pretty bool) {
  if pretty {
    out.Write([]byte{ '(' })
    this.left.Dump(out, "", true)
    fmt.Fprintf(out, " %s ", this.lit)
    this.right.Dump(out, "", true)
    out.Write([]byte{ ')' })
  } else {
    fmt.Fprintf(out, "%sInfixExpr:%s:%s\n", prfx, this.lit, this.ambit)
    this.left.Dump(out, prfx + "  ", false)
    this.right.Dump(out, prfx + "  ", false)
  }
}

func (this *trueExpr) Dump(out io.Writer, prfx string, pretty bool) {
  if pretty {
    fmt.Fprintf(out, "true")
  } else {
    fmt.Fprintf(out, "%sTrueExpr:%s\n", prfx, this.ambit)
  }
}

func (this *falseExpr) Dump(out io.Writer, prfx string, pretty bool) {
  if pretty {
    fmt.Fprintf(out, "false")
  } else {
    fmt.Fprintf(out, "%sFalseExpr:%s\n", prfx, this.ambit)
  }
}

func (this *noopSimple) Dump(out io.Writer, prfx string, pretty bool) {
  if pretty {
    fmt.Fprintf(out, "noop")
  } else {
    fmt.Fprintf(out, "%sNoopExpr:%s\n", prfx, this.ambit)
  }
}

func (this *assignSimple) Dump(out io.Writer, prfx string, pretty bool) {
  if pretty {
    fmt.Fprintf(out, "%s = ", this.varName)
    this.expr.Dump(out, "", true)
  } else {
    fmt.Fprintf(out, "%sAssignSimple:%s:%s\n", prfx, this.varName, this.ambit)
    this.expr.Dump(out, prfx + "  ", false)
  }
}

func (this *incrSimple) Dump(out io.Writer, prfx string, pretty bool) {
  if pretty {
    fmt.Fprintf(out, "%s++", this.varName)
  } else {
    fmt.Fprintf(out, "%sIncrSimple:%s:%s\n", prfx, this.varName, this.ambit)
  }
}

func (this *decrSimple) Dump(out io.Writer, prfx string, pretty bool) {
  if pretty {
    fmt.Fprintf(out, "%s--", this.varName)
  } else {
    fmt.Fprintf(out, "%sDecrSimple:%s:%s\n", prfx, this.varName, this.ambit)
  }
}

func (this *exprSimple) Dump(out io.Writer, prfx string, pretty bool) {
  if pretty {
    this.expr.Dump(out, "", true)
  } else {
    fmt.Fprintf(out, "%sExprSimple:%s\n", prfx, this.ambit)
    this.expr.Dump(out, prfx + "  ", false)
  }
}

func (this *seqStmt) Dump(out io.Writer, prfx string, pretty bool) {
  if pretty {
    this.head.Dump(out, prfx, true)
    this.tail.Dump(out, prfx, true)
  } else {
    fmt.Fprintf(out, "%sSeqStmt:%s\n", prfx, this.ambit)
    this.head.Dump(out, prfx + "  ", false)
    this.tail.Dump(out, prfx + "  ", false)
  }
}

func (this *emptyStmt) Dump(out io.Writer, prfx string, pretty bool) {
  if pretty {
    // NOOP
  } else {
    fmt.Fprintf(out, "%sEmptyStmt:%s\n", prfx, this.ambit)
  }
}

func (this *ifStmt) Dump(out io.Writer, prfx string, pretty bool) {
  if pretty {
    fmt.Fprintf(out, "%sif ", prfx)
    this.cond.Dump(out, "", true)
    out.Write([]byte{ '\n' })
    this.then.Dump(out, prfx + "  ", true)
    if _, empty := this.el5e.(*emptyStmt); !empty {
      fmt.Fprintf(out, "%selse\n", prfx)
      this.el5e.Dump(out, prfx + "  ", true)
    }
  } else {
    fmt.Fprintf(out, "%sIfStmt:%s\n", prfx, this.ambit)
    this.cond.Dump(out, prfx + "  ", false)
    this.then.Dump(out, prfx + "  ", false)
    this.el5e.Dump(out, prfx + "  ", false)
  }
}

func (this *forStmt) Dump(out io.Writer, prfx string, pretty bool) {
  if pretty {
    fmt.Fprintf(out, "%sfor ", prfx)
    this.initial.Dump(out, "", true)
    fmt.Fprintf(out, "; ")
    this.cond.Dump(out, "", true)
    fmt.Fprintf(out, "; ")
    this.update.Dump(out, "", true)
    if _, empty := this.body.(*emptyStmt); !empty {
      out.Write([]byte{ '\n' })
      this.body.Dump(out, prfx + "  ", true)
    }
  } else {
    fmt.Fprintf(out, "%sForStmt:%s\n", prfx, this.ambit)
    this.initial.Dump(out, prfx + "  ", false)
    this.cond.Dump(out, prfx + "  ", false)
    this.update.Dump(out, prfx + "  ", false)
    this.body.Dump(out, prfx + "  ", false)
  }
}

func (this *simpleStmt) Dump(out io.Writer, prfx string, pretty bool) {
  if pretty {
    fmt.Fprintf(out, "%s", prfx)
    this.simple.Dump(out, "", true)
    out.Write([]byte{ '\n' })
  } else {
    fmt.Fprintf(out, "%sSimpleStmt:%s\n", prfx, this.ambit)
    this.simple.Dump(out, prfx + "  ", false)
  }
}
