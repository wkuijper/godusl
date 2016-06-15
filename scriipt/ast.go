package scriipt

import (
  "io"
)

type AST interface {
  Location() string
  Dump(out io.Writer, prfx string, pretty bool)
}

type Expr interface {
  AST
  Eval(map[string]interface{}) interface{}
}

type Simple interface {
  AST
  Run(st map[string]interface{})
}

type Stmt interface {
  AST
  Run(map[string]interface{})
}

type IdExpr interface {
  Expr
  IdName() string
}

type NumExpr interface {
  Expr
  Num() int
}

type StrExpr interface {
  Expr
  EscapedStr() string
  UnescapedStr() string
}

type CallExpr interface {
  Expr
  FuncName() string
  Args() []Expr
}

type PrefixExpr interface {
  Expr
  Lit() string
  Sub() Expr
}

type PostfixExpr interface {
  Expr
  Lit() string
  Sub() Expr
}

type InfixExpr interface {
  Expr
  Lit() string
  Left() Expr
  Right() Expr
}

type TrueExpr interface {
  Expr
}

type FalseExpr interface {
  Expr
}

type NoopSimple interface {
  Simple
}

type AssignSimple interface {
  Simple
  VarName() string
  Expr() Expr
}

type IncrSimple interface {
  Simple
  VarName() string
}

type DecrSimple interface {
  Simple
  VarName() string
}

type ExprSimple interface {
  Simple
  Expr() Expr
}

type SeqStmt interface {
  Stmt
  Head() Stmt
  Tail() Stmt
}

type EmptyStmt interface {
  Stmt
}

type ForStmt interface {
  Stmt
  Init() Simple
  Cond() Expr
  Update() Simple
  Body() Stmt
}

type IfStmt interface {
  Stmt
  Cond() Expr
  Then() Stmt
  Else() Stmt
}

type SimpleStmt interface {
  Stmt
  Simple() Simple
}
