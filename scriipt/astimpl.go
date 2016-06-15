package scriipt

import (
  "dusl"
  "strings"
  "strconv"
)

type ast struct {
  ambit *dusl.Ambit
}

func (this *ast) Location() string {
  return this.ambit.Location()
}

func (this *ast) init(ambit *dusl.Ambit) {
  this.ambit = ambit
}

type expr struct {
  ast
}

type simple struct {
  ast
}

type stmt struct {
  ast
}

type trueExpr struct {
  expr
}


func newTrueExpr(ambit *dusl.Ambit) *trueExpr {
  x := &trueExpr{}
  x.expr.init(ambit)
  return x
}

type falseExpr struct {
  expr
}

func newFalseExpr(ambit *dusl.Ambit) *falseExpr {
  x := &falseExpr{}
  x.expr.init(ambit)
  return x
}

type strExpr struct {
  expr
  escapedStr string
  unescapedStr string
}

func newStrExpr(ambit *dusl.Ambit) *strExpr {
  x := &strExpr{}
  x.expr.init(ambit)
  x.escapedStr = ambit.ToString()
  x.unescapedStr = strings.Replace(strings.Replace(strings.Replace(strings.Replace(x.escapedStr[1:len(x.escapedStr)-1], `\n`, "\n", -1), `\r`, "\r", -1), `\t`, "\t", -1), `\"`, "\"", -1)
  
  return x
}

func (this *strExpr) EscapedStr() string {
  return this.escapedStr
}

type numExpr struct {
  expr
  num int
}

func newNumExpr(ambit *dusl.Ambit) *numExpr {
  x := &numExpr{}
  x.expr.init(ambit)
  num, _ := strconv.Atoi(ambit.ToString())
  x.num = num
  return x
}

type idExpr struct {
  expr
  idName string
}

func newIdExpr(ambit *dusl.Ambit) *idExpr {
  x := &idExpr{}
  x.expr.init(ambit)
  x.idName = ambit.ToString()
  return x
}

type callExpr struct {
  expr
  funcName string
  args []Expr
}

func newCallExpr(ambit *dusl.Ambit, funcName string, args []Expr) *callExpr {
  x := &callExpr{}
  x.expr.init(ambit)
  x.funcName = funcName
  x.args = args
  return x
}

type prefixExpr struct {
  expr
  lit string
  sub Expr
}

func newPrefixExpr(ambit *dusl.Ambit, lit string, sub Expr) *prefixExpr {
  x := &prefixExpr{}
  x.expr.init(ambit)
  x.lit = lit
  x.sub = sub
  return x
}

type postfixExpr struct {
  expr
  lit string
  sub Expr
}

func newPostfixExpr(ambit *dusl.Ambit, lit string, sub Expr) *postfixExpr {
  x := &postfixExpr{}
  x.expr.init(ambit)
  x.lit = lit
  x.sub = sub
  return x
}

type infixExpr struct {
  expr
  lit string
  left Expr
  right Expr
}

func newInfixExpr(ambit *dusl.Ambit, lit string, left Expr, right Expr) *infixExpr {
  x := &infixExpr{}
  x.expr.init(ambit)
  x.lit = lit
  x.left = left
  x.right = right
  return x
}

type noopSimple struct {
  simple
}

func newNoopSimple(ambit *dusl.Ambit) *noopSimple {
  x := &noopSimple{}
  x.simple.init(ambit)
  return x
}

type assignSimple struct {
  simple
  varName string
  expr Expr
}

func newAssignSimple(ambit *dusl.Ambit, varName string, expr Expr) *assignSimple {
  x := &assignSimple{}
  x.simple.init(ambit)
  x.varName = varName
  x.expr = expr
  return x
}

type incrSimple struct {
  simple
  varName string
}

func newIncrSimple(ambit *dusl.Ambit, varName string) *incrSimple {
  x := &incrSimple{}
  x.simple.init(ambit)
  x.varName = varName
  return x
}

type decrSimple struct {
  simple
  varName string
}

func newDecrSimple(ambit *dusl.Ambit, varName string) *decrSimple {
  x := &decrSimple{}
  x.simple.init(ambit)
  x.varName = varName
  return x
}

type exprSimple struct {
  simple
  expr Expr
}

func newExprSimple(ambit *dusl.Ambit, expr Expr) *exprSimple {
  x := &exprSimple{}
  x.simple.init(ambit)
  x.expr = expr
  return x
}

type seqStmt struct {
  stmt
  head Stmt
  tail Stmt
}

func newSeqStmt(ambit *dusl.Ambit, head Stmt, tail Stmt) *seqStmt {
  x := &seqStmt{}
  x.stmt.init(ambit)
  x.head = head
  x.tail = tail
  return x
}

type emptyStmt struct {
  stmt
}

func newEmptyStmt(ambit *dusl.Ambit) *emptyStmt {
  x := &emptyStmt{}
  x.stmt.init(ambit)
  return x
}

type forStmt struct {
  stmt
  initial Simple
  cond Expr
  update Simple
  body Stmt
}

func newForStmt(ambit *dusl.Ambit,
                initial Simple,
                cond Expr,
                update Simple,
                body Stmt) *forStmt {
  x := &forStmt{}
  x.stmt.init(ambit)
  x.initial = initial
  x.cond = cond
  x.update = update
  x.body = body
  return x
}

type ifStmt struct {
  stmt
  cond Expr
  then Stmt
  el5e Stmt
}

func newIfStmt(ambit *dusl.Ambit,
               cond Expr,
               then Stmt,
               el5e Stmt) *ifStmt {
  x := &ifStmt{}
  x.stmt.init(ambit)
  x.cond = cond
  x.then = then
  x.el5e = el5e
  return x
}

type simpleStmt struct {
  stmt
  simple Simple
}

func newSimpleStmt(ambit *dusl.Ambit, simple Simple) *simpleStmt {
  x := &simpleStmt{}
  x.stmt.init(ambit)
  x.simple = simple
  return x
}
