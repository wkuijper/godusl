package scriipt

import (
  "dusl"
  "fmt"
)

func Parse(src *dusl.Source) (Stmt, error) {
  trace := Lang.Tracer().TraceUndent(src, "S")
  if err := trace.ErrorN(20); err != nil {
    return nil, err
  }
  return parseStmt(trace), nil
}

func parseExpr(exprTrace *dusl.Trace) Expr {
  if exprTrace.Lbl != "x" {
    panic(exprTrace.Lbl + " != x")
  }
  switch exprTrace.Idx {
  case 0:
    return newIdExpr(exprTrace.Syn.Ambit)
  case 1:
    return newNumExpr(exprTrace.Syn.Ambit)
  case 2:
    return newStrExpr(exprTrace.Syn.Ambit)
  case 3:
    return parseExpr(exprTrace.Subs[0])
  case 4:
    return newCallExpr(exprTrace.Syn.Ambit,
                       exprTrace.Subs[0].Syn.Lit, // <- funcName
                       parseArgs(exprTrace.Subs[1]))
  case 5:
    return newPrefixExpr(exprTrace.Syn.Ambit,
                         exprTrace.Syn.Lit, // <- prefix operator
                         parseExpr(exprTrace.Subs[0]))
  case 6:
    return newPostfixExpr(exprTrace.Syn.Ambit,
                          exprTrace.Syn.Lit, // <- postfix operator
                          parseExpr(exprTrace.Subs[0]))
  case 7:
    return newInfixExpr(exprTrace.Syn.Ambit,
                        exprTrace.Syn.Lit, // <- infix operator
                        parseExpr(exprTrace.Subs[0]),
                        parseExpr(exprTrace.Subs[1]))
  case 8:
    return newTrueExpr(exprTrace.Syn.Ambit)
  case 9:
    return newFalseExpr(exprTrace.Syn.Ambit)
  default:
    panic("missing case")
  }
}

func parseArgs(argsTrace *dusl.Trace) []Expr {
  if argsTrace.Lbl != "args" {
    panic(argsTrace.Lbl + " != args")
  }
  switch argsTrace.Idx {
  case 0:
    return []Expr{}
  case 1:
    return parseAargs(argsTrace.Subs[0])
  default:
    panic("missing case")
  }
}

func parseAargs(aargsTrace *dusl.Trace) []Expr {
  if aargsTrace.Lbl != "aargs" {
    panic(aargsTrace.Lbl + " != aargs")
  }
  switch aargsTrace.Idx {
  case 0:
    return append(parseAargs(aargsTrace.Subs[0]), parseExpr(aargsTrace.Subs[1]))
  case 1:
    return []Expr{ parseExpr(aargsTrace.Subs[0]) }
  default:
    panic("missing case")
  }
}

func parseStmt(stmtTrace *dusl.Trace) Stmt {
  if stmtTrace.Lbl != "S" {
    panic(stmtTrace.Lbl + " != S")
  }
  var head Stmt
  var contTrace *dusl.Trace
  switch stmtTrace.Idx {
  case 0: 
    return newEmptyStmt(stmtTrace.Syn.Ambit)
  case 1, 2:
    head, contTrace = parseForStmt(stmtTrace)
  case 3:
    head, contTrace = parseIfStmt(stmtTrace)
  case 4:
    head, contTrace = parseSimpleStmt(stmtTrace)
  default:
    panic("missing case")
  }
  return newSeqStmt(stmtTrace.Syn.Ambit, head, parseStmt(contTrace))
}

func parseSimpleStmt(simpleStmtTrace *dusl.Trace) (Stmt, *dusl.Trace) {
  if simpleStmtTrace.Lbl != "S" || simpleStmtTrace.Idx != 4 {
    panic(fmt.Sprintf("%s != S || %d != 3", simpleStmtTrace.Lbl, simpleStmtTrace.Idx))
  }
  return newSimpleStmt(simpleStmtTrace.Syn.Ambit,
                       parseSimple(simpleStmtTrace.Subs[0])), simpleStmtTrace.Subs[1]
}

func parseSimple(simpleTrace *dusl.Trace) Simple {
  if simpleTrace.Lbl != "s" {
    panic(simpleTrace.Lbl + " != s")
  }
  subs := simpleTrace.Subs
  switch simpleTrace.Idx {
  case 0:
    return newNoopSimple(simpleTrace.Syn.Ambit)
  case 1:
    return newAssignSimple(simpleTrace.Syn.Ambit,
                           subs[0].Syn.Lit,
                           parseExpr(subs[1]))
  case 2:
    return newIncrSimple(simpleTrace.Syn.Ambit,
                         subs[0].Syn.Lit)
  case 3:
    return newDecrSimple(simpleTrace.Syn.Ambit,
                         subs[0].Syn.Lit)
  case 4:
    return newExprSimple(simpleTrace.Syn.Ambit,
                         newCallExpr(simpleTrace.Syn.Ambit,
                                     simpleTrace.Subs[0].Syn.Lit,
                                     parseArgs(simpleTrace.Subs[1])))
  default:
    panic("missing case")
  }
}

func parseForStmt(forTrace *dusl.Trace) (Stmt, *dusl.Trace) {
  if forTrace.Lbl != "S" || (forTrace.Idx != 1 && forTrace.Idx != 2) {
    panic(fmt.Sprintf("%s != S || %d != 1", forTrace.Lbl, forTrace.Idx))
  }
  subs := forTrace.Subs
  contTrace := subs[2]
  body := parseStmt(subs[1])
  hdr := subs[0]
  switch hdr.Idx {
  case 0:
    return newForStmt(forTrace.Syn.Ambit,
                      newNoopSimple(hdr.Syn.Ambit),
                      newTrueExpr(hdr.Syn.Ambit),
                      newNoopSimple(hdr.Syn.Ambit),
                      body), contTrace
  case 1:
    return newForStmt(forTrace.Syn.Ambit,
                      parseSimple(hdr.Subs[0]),
                      parseExpr(hdr.Subs[1]),
                      parseSimple(hdr.Subs[2]),
                      body), contTrace
  case 2:
    return newForStmt(forTrace.Syn.Ambit,
                      newNoopSimple(hdr.Syn.Ambit.CollapseLeft()),
                      parseExpr(hdr.Subs[0]),
                      newNoopSimple(hdr.Syn.Ambit.CollapseRight()),                      
                      body), contTrace
  default:
    panic("missing case")
  }
}

func parseIfStmt(ifTrace *dusl.Trace) (Stmt, *dusl.Trace) {
  return parseElseCont(&dusl.Trace{ Syn: ifTrace.Syn, Lbl: "E", Idx: 0, Subs: ifTrace.Subs })
}

func parseElseCont(elseTrace *dusl.Trace) (Stmt, *dusl.Trace) {
  if elseTrace.Lbl != "E" {
    panic(elseTrace.Lbl + " != E")
  }
  subs := elseTrace.Subs
  switch elseTrace.Idx {
  case 0:
    el5e, contTrace := parseElseCont(subs[2])
    return newIfStmt(elseTrace.Syn.Ambit,
                     parseExpr(subs[0]),
                     parseStmt(subs[1]),
                     el5e), contTrace
  case 1:
    return parseStmt(subs[0]), subs[1]
  case 2, 3:
    return newEmptyStmt(elseTrace.Syn.Ambit), subs[0]
  default:
    panic("missing case")
  }
}
