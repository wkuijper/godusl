package dusl

import (
  "strings"
  "bytes"
  "errors"
  "runtime"
  "fmt"
  //"os"
)

const maxPrecedence = 1000000000

type Spec interface {
  Lexical(scanner Scanner) Spec
  Category(cat string, desc string) Spec
  OperatorAFB(ops ...string) Spec
  OperatorBFA(ops ...string) Spec
  OperatorEFA(ops ...string) Spec
  OperatorAFE(ops ...string) Spec
  OperatorEFE(ops ...string) Spec
  Brackets(pairs ...string) Spec
  SequenceLabel(lbl string, desc string) Spec
  SentenceLabel(lbl string, desc string) Spec
  Label(lbl string, desc string) Spec
  Literal(lit string, cat string) Spec
  ShorthandOperator(op string, ops ...string) Spec
  Grammar(grammar string) (Lang, error)
}

type Lang interface {
  Tokenizer() Tokenizer
  Spanner() Spanner
  Sparser() Sparser
  Tracer() Tracer
}

type spec struct {
  scanner Scanner
  layers []*specLayer
  symbols []*specSymbol
}

type lang struct {
  tokenizer Tokenizer
  spanner Spanner
  sparser Sparser
  tracer Tracer
}

func (this *lang) Tokenizer() Tokenizer {
  return this.tokenizer
}

func (this *lang) Spanner() Spanner {
  return this.spanner
}

func (this *lang) Sparser() Sparser {
  return this.sparser
}

func (this *lang) Tracer() Tracer {
  return this.tracer
}

type specLayer struct {
  pattern string
  args []string
}

const (
  spec_Label = iota
  spec_SentenceLabel
  spec_SequenceLabel
  spec_Literal
  spec_Category
  spec_ShorthandOperator
)

type specSymbol struct {
  symb string
  typ int
  lbl string
  cat string
  lit string
  desc string
  ops []string
}

func (this *specSymbol) typName() string {
  switch this.typ {
  case spec_Label:
    return "label"
  case spec_SentenceLabel:
    return "sentence label"
  case spec_SequenceLabel:
    return "sequence label"
  case spec_Literal:
    return "literal"
  case spec_Category:
    return "category"
  }
  return "<<<missing typName>>>"
}

type precedenceLevels struct {
  precedenceB map[string]int
  precedenceEFE map[string]int
  precedenceEFA map[string]int
  precedenceAFE map[string]int
  precedenceAFB map[string]int
  precedenceBFA map[string]int
}

func NewSpec() Spec {
  return &spec{}
}

func (this *spec) Lexical(scanner Scanner) Spec {
  if this.scanner == nil {
    this.scanner = scanner
  } else {
    this.scanner = &seqScanner{ master: scanner, slave: this.scanner }
  }
  return this
}

func (this *spec) ShorthandOperator(op string, ops ...string) Spec {
  this.symbols = append(this.symbols, &specSymbol{ typ: spec_ShorthandOperator, symb: op, ops: ops })
  return this
}

func (this *spec) OperatorAFB(ops ...string) Spec {
  return this.layer("AFB", ops)
}

func (this *spec) OperatorBFA(ops ...string) Spec {
  return this.layer("BFA", ops)
}

func (this *spec) OperatorEFA(ops ...string) Spec {
  return this.layer("EFA", ops)
}

func (this *spec) OperatorAFE(ops ...string) Spec {
  return this.layer("AFE", ops)
}

func (this *spec) OperatorEFE(ops ...string) Spec {
  return this.layer("EFE", ops)
}

func (this *spec) Brackets(ops ...string) Spec {
  return this.layer("B", ops)
}

func (this *spec) layer(pattern string, args []string) Spec {
  this.layers = append(this.layers, &specLayer{ pattern: pattern, args: args })
  return this
}

func (this *spec) Category(cat string, desc string) Spec {
  return this.symbol(spec_Category, cat, "", cat, "", desc)
}

func (this *spec) Literal(lit string, cat string) Spec {
  return this.symbol(spec_Literal, lit, "", cat, lit, lit)
}

func (this *spec) SequenceLabel(lbl string, desc string) Spec {
  return this.symbol(spec_SequenceLabel, lbl, lbl, "SQ", "", desc)
}

func (this *spec) SentenceLabel(lbl string, desc string) Spec {
  return this.symbol(spec_SentenceLabel, lbl, lbl, "SQ", "", desc)
}

func (this *spec) Label(lbl string, desc string) Spec {
  return this.symbol(spec_Label, lbl, lbl, "", "", desc)
}

func (this *spec) symbol(typ int, symb, lbl string, cat string, lit string, desc string) Spec {
  this.symbols = append(this.symbols, &specSymbol{ typ: typ, symb: symb, lbl: lbl, cat: cat, lit: lit, desc: desc })
  return this
}

func (this *spec) Grammar(grammar string) (Lang, error) {
  
  precMap := make(map[string]map[string]int, 16)
  layers := this.layers
  l := len(layers)-1
  for i := 0; i <= l; i++ {
    layer := layers[i]
    precedence := 10 + (l-i)
    pattMap := precMap[layer.pattern]
    if pattMap == nil {
      pattMap = make(map[string]int, 8*len(layer.args))
      precMap[layer.pattern] = pattMap
    }
    for _, arg := range layer.args {
      if pattMap[arg] != 0 {
        return nil, fmt.Errorf("double declaration of %s operator/bracket: '%s'", layer.pattern, arg)
      }
      pattMap[arg] = precedence
    }
  }
  /*
  conflicts := map[string][]string {
    "BFA": []string{ "AFB", "AFE" },
    "AFB": []string{ "BFA", "EFA" },
    "EFA": []string{ "AFE" },
    "AFE": []string{ "EFA" },
  }
  for pattern, pattMap := range precMap {
    for _, conflictPattern := range conflicts[pattern] {
      conflictPattMap := precMap[conflictPattern]
      for op, _ := range pattMap {
        if _, present := conflictPattMap[op]; present {
          return nil, fmt.Errorf("operator cannot simultaneously be declared as %s and %s operator: '%s'", pattern, conflictPattern, op)
        }
      }
    }
  }*/

  prfxScanner := &prfxTree{}
  var scanner Scanner
  if this.scanner == nil {
    scanner = prfxScanner
  } else {
    scanner = &seqScanner{ master: prfxScanner, slave: this.scanner }
  }
  
  prfxMetaScanner := &prfxTree{}
  var metaScanner Scanner
  if this.scanner == nil {
    metaScanner = &seqScanner{ master: prfxMetaScanner,
                               slave: metaSymbolScanner }
  } else {
    metaScanner = &seqScanner{ master: prfxMetaScanner,
                               slave: &seqScanner{ master: metaSymbolScanner, slave: this.scanner } }
  }
  
  for _, patt := range []string{ "AFE", "EFA", "AFB", "BFA", "EFE" } {
    for op, _ := range precMap[patt] {
      prfxScanner.add("OP", op)
      prfxMetaScanner.add("OP", op)
    }
  }
  
  for brs, _ := range precMap["B"] {
    parts := strings.Split(brs, " ")
    if len(parts) < 2 {
      return nil, fmt.Errorf("expected pair of brackets separated by blank space: '%s'", brs)
    }
    if len(parts) > 2 {
      return nil, fmt.Errorf("expected pair of brackets separated by single blank space: '%s'", brs)
    }
    ob, cb := parts[0], parts[1]
    obExisting := prfxScanner.lookup(ob)
    if obExisting == "OP" {
      return nil, fmt.Errorf("declared open bracket conflicts with declared operator: '%s'", ob)
    }
    if obExisting == "CB" {
      return nil, fmt.Errorf("declared open bracket conflicts with declared close bracket: '%s'", ob)
    }
    if obExisting != "" {
      return nil, fmt.Errorf("double declaration of open bracket: '%s'", ob)
    }
    cbExisting := prfxScanner.lookup(cb)
    if cbExisting == "OP" {
      return nil, fmt.Errorf("declared close bracket conflicts with declared operator: '%s'", cb)
    }
    if cbExisting == "CB" {
      return nil, fmt.Errorf("declared close bracket conflicts with declared open bracket: '%s'", cb)
    }
    if cbExisting != "" {
      return nil, fmt.Errorf("double declaration of close bracket: '%s'", cb)
    }
    prfxScanner.add("OB", ob)
    prfxScanner.add("CB", cb)
    prfxMetaScanner.add("OB", ob)
    prfxMetaScanner.add("CB", cb)
  }

  if prfxScanner.lookup("is>") != "" {
    return nil, fmt.Errorf("conflicting declaration of meta operator: 'is>'")
  }
  if prfxScanner.lookup("or>") != "" {
    return nil, fmt.Errorf("conflicting declaration of meta operator: 'or>'")
  }
  if prfxScanner.lookup("<empty") != "" {
    return nil, fmt.Errorf("conflicting declaration of meta operator: '<empty'")
  }

  prfxMetaScanner.add("OP", "is>")
  prfxMetaScanner.add("OP", "or>")
  prfxMetaScanner.add("OP", "<empty")

  if precMap["ABF"] == nil {
    precMap["AFB"] = make(map[string]int, 2)
  }
  if precMap["AFE"] == nil {
    precMap["AFE"] = make(map[string]int, 2)
  }
  if precMap["EFE"] == nil {
    precMap["EFE"] = make(map[string]int, 2)
  }
  
  precMap["AFB"]["is>"] = 1
  precMap["AFE"]["is>"] = 2
  precMap["AFB"]["or>"] = 3
  precMap["EFE"]["or>"] = 4
  precMap["EFE"]["<empty"] = 5
  
  precedence := &precedenceLevels {
    precedenceB: precMap["B"],
    precedenceEFE: precMap["EFE"],
    precedenceEFA: precMap["EFA"],
    precedenceAFE: precMap["AFE"],
    precedenceAFB: precMap["AFB"],
    precedenceBFA: precMap["BFA"],
  }

  symbolTable := make(map[string]*specSymbol, len(this.symbols))

  reserved := map[string]string{ "ERR": "error category",
                                 "OP": "operator category",
                                 "OB": "open-bracket category",
                                 "CB": "close-bracket category",
                                 "BB": "bracket-pair category",
                                 "WS": "whitespace category" }
  for _, symbol := range this.symbols {
    symb := symbol.symb
    trimmedSymb := strings.TrimSpace(symb)
    if trimmedSymb == "" {
      return nil, fmt.Errorf("cannot declare empty string as %s: '%s'", symbol.typName(), symb)
    }
    if trimmedSymb != symb {
      return nil, fmt.Errorf("leading/trailing whitespace in %s: '%s'", symbol.typName(), symb)
    }
    for i, c := range symb {
      if (c >= 'A' && c <= 'Z') || (c >= 'a' && c <= 'z') || c == '_' {
        continue
      }
      if i > 0 && (c >= '0' && c <= '9') {
        continue
      }
      return nil, fmt.Errorf("unexpected symbol in %s: '%s***HERE***%s'",
                             symbol.typName(), symb[:i], symb[i:])
    }
    for reservedSymb, reservedSymbName := range reserved {
      if symb == reservedSymb {
        return nil, fmt.Errorf("%s conflicts with reserved %s: '%s'", symbol.typName(), reservedSymbName, symb)
      }
    }
    if prfxScanner.lookup(symb) != "" {
      return nil, fmt.Errorf("%s conflicts with operator/bracket: '%s'", symbol.typName(), symb)
    }
    existingSymbol := symbolTable[symb]
    if existingSymbol != nil {
      if existingSymbol.typName() == symbol.typName() {
        return nil, fmt.Errorf("double declaration of %s: '%s'", symbol.typName(), symb)
      } else {
        return nil, fmt.Errorf("%s conflicts with %s: '%s'", symbol.typName(), existingSymbol.typName(), symb)
      }
    }
    switch symbol.typ {
    case spec_SequenceLabel:
      fallthrough
    case spec_SentenceLabel:
      fallthrough
    case spec_Label:
      // NOOP
    case spec_Category:
      // NOOP
    case spec_ShorthandOperator:
      var pEFE, pEFA, pAFE, pBFA, pAFB int
      ops := make([]string, 0, len(symbol.ops))
      for _, op := range symbol.ops {
        existingSymbol := symbolTable[op]
        if existingSymbol != nil && existingSymbol.typ == spec_ShorthandOperator {
          ops = append(ops, existingSymbol.ops...)
        } else {
          ops = append(ops, op)
        }
      }
      for _, op := range ops {
        if prfxScanner.lookup(op) != "OP" {
          return nil, fmt.Errorf("undeclared symbol: %s: in definition of shorthand operator: %s", op, symbol.symb)
        }
        if p := precedence.precedenceEFE[op]; pEFE == 0 || p < pEFE { pEFE = p }
        if p := precedence.precedenceEFA[op]; pEFA == 0 || p < pEFA { pEFA = p }
        if p := precedence.precedenceAFE[op]; pAFE == 0 || p < pAFE { pAFE = p }
        if p := precedence.precedenceBFA[op]; pBFA == 0 || p < pBFA { pBFA = p }
        if p := precedence.precedenceAFB[op]; pAFB == 0 || p < pAFB { pAFB = p }
      }
      precedence.precedenceEFE[symbol.symb] = pEFE
      precedence.precedenceEFA[symbol.symb] = pEFA
      precedence.precedenceAFE[symbol.symb] = pAFE
      precedence.precedenceBFA[symbol.symb] = pBFA
      precedence.precedenceAFB[symbol.symb] = pAFB
      prfxMetaScanner.add("OP", symbol.symb)
    case spec_Literal:
      prfxMetaScanner.add(symbol.cat, symbol.lit)
    }
    symbolTable[symb] = symbol
  }
  
  metaTokenizer := newTokenizer(metaScanner)
  metaSpanner := newSpanner(metaTokenizer, precedence.precedenceB)
  metaSparser := newSparser(metaSpanner, precedence)
  
  _, path, lastLineOffset, _ := runtime.Caller(1)
  grammarSource := &Source{ Path: path, LineOffset: lastLineOffset-strings.Count(grammar, "\n")-1, Text: []byte(grammar) }

  grammarTree := metaSparser.SparseUndent(grammarSource)

  //grammarTree.Dump(os.Stdout, "grammar> ")
  
  templateParser := &tpT{ symbolTable: symbolTable, templates: make(map[string][]*templateT) }
  
  for _, errNode := range grammarTree.FirstN("ERR", "", 20) {
    templateParser.err(errNode, errNode.Err)
  }

  templateParser.topSequence(grammarTree)

  if len(templateParser.errs) > 0 {
    buf := new(bytes.Buffer)
    for index, err := range templateParser.errs {
      if index < 20 {
        fmt.Fprintf(buf, "%s\n", err.Error())
      } else if index == 20 {
        fmt.Fprintf(buf, "and %d more error(s)\n", len(templateParser.errs)-20)
      } else {
        break
      }
    }
    return nil, errors.New(buf.String())
  }

  descriptions := make(map[string]string, len(symbolTable))

  for symb, symbol := range symbolTable {
    descriptions[symb] = symbol.desc
  }
  
  tokenizer := newTokenizer(scanner)
  spanner := newSpanner(tokenizer, precedence.precedenceB)
  sparser := newSparser(spanner, precedence)
  tracer := newTracer(sparser, templateParser.templates, descriptions)

  return &lang{ tokenizer: tokenizer, spanner: spanner, sparser: sparser, tracer: tracer }, nil
}

type tpT struct {
  symbolTable map[string]*specSymbol
  templates map[string][]*templateT
  errs []error
}

func (this *tpT) err(node *Syntax, format string, args ...interface{}) {
  ambit := node.Ambit
  if node.OpAmbit != nil {
    ambit = node.OpAmbit
  }
  this.errs = append(this.errs, fmt.Errorf(ambit.Location() + ": " + format, args...))
}

func (this *tpT) topSequence(node *Syntax) {
  if node.Cat == "SQ" && node.Left.Cat == "SN" {
    if node.Left.Right.Cat == "SQ" {
      this.multiSentenceRule(node)
    } else {
      this.singleSentenceRule(node.Left.Left)
      this.topSequence(node.Right)
    }
    return
  }
  if node.Cat == "" {
    // <empty
    return
  }
  this.err(node, "expected rule")
}

func (this *tpT) multiSentenceRule(node *Syntax) {
  sn, lbl := this.multiSentenceRuleHead(node.Left.Left)
  this.multiSentenceRuleBody(node.Left.Right, sn, lbl)
  this.multiSentenceRuleContinuation(node.Right, sn, lbl)
}

func (this *tpT) multiSentenceRuleContinuation(node *Syntax, sn bool, lbl string) {
  if node.Cat == "SQ" && node.Left.Cat == "SN" && node.Left.Left.IsZeroaryOp("or>") {
    this.multiSentenceRuleBody(node.Left.Right, sn, lbl)
    this.multiSentenceRuleContinuation(node.Right, sn, lbl)
  } else {
    this.topSequence(node)
  }
}

func (this *tpT) multiSentenceRuleBody(node *Syntax, sn bool, lbl string) {
  template := this.multiSentenceTemplate(node)
  if template != nil && lbl != "" {
    if sn {
      if !(template.cat == "SQ" && template.left.cat == "SN" && template.right.cat == "") {
        this.err(node, "cannot match sequence with sentence label: '%s'", lbl)
        return
      }
      template = template.left
    }
    this.templates[lbl] = append(this.templates[lbl], template)
  }
}

func (this *tpT) multiSentenceRuleHead(node *Syntax) (bool, string) {
  if !node.IsPostfixOp("is>") {
    this.err(node, "expected: <label> is>")
    return false, ""
  }
  left := node.Left
  if left.Cat != "$" {
    this.err(left, "expected: <label> is>")
    return false, ""
  }
  symbol := this.symbolTable[left.Lit]
  if symbol == nil {
    this.err(left, "undeclared symbol: '%s'", left.Lit)
    return false, ""
  }
  var lbl string
  var sn bool
  switch symbol.typ {
  case spec_SentenceLabel:
    sn = true
    fallthrough
  case spec_SequenceLabel:
    lbl = left.Lit
  case spec_Label:
    this.err(left, "expected sequence label instead of ordinary label: '%s'", left.Lit)
  case spec_Literal:
    this.err(left, "expected sequence label instead of literal: '%s'", left.Lit)
  case spec_Category:
    this.err(left, "expected sequence label instead of category: '%s'", left.Lit)
  default:
    this.err(left, "expected sequence label instead of: '%s'", left.Lit)
  }
  return sn, lbl
}

func (this *tpT) singleSentenceRule(node *Syntax) {
  if !node.IsInfixOp("is>") {
    this.err(node, "expected: <label> is> ...")
    return
  }
  left := node.Left
  if left.Cat != "$" {
    this.err(left, "expected: <label> is> ...")
    return
  }
  symbol := this.symbolTable[left.Lit]
  if symbol == nil {
    this.err(left, "undeclared symbol: '%s'", left.Lit)
    return
  }
  var lbl string
  var sn bool
  switch symbol.typ {
  case spec_SequenceLabel:
    this.err(left, "missing indented body for sequence label: '%s'", left.Lit)
    return
  case spec_SentenceLabel:
    sn = true
    fallthrough
  case spec_Label:
    lbl = left.Lit
  case spec_Literal:
    this.err(left, "expected label instead of literal: '%s'", left, left.Lit)
    return
  case spec_Category:
    this.err(left, "expected label instead of category: '%s'", left.Lit)
    return
  default:
    this.err(left, "expected label instead of: '%s'", left.Lit)
    return
  }
  this.singleSentenceRuleBody(node.Right, sn, lbl)
}

func (this *tpT) singleSentenceRuleBody(node *Syntax, sn bool, lbl string) {
  if node.Cat == "OP" && node.Lit == "or>" {
    // order of invocation matters here:
    this.singleSentenceRuleBody(node.Left, sn, lbl) 
    this.singleSentenceRuleBody(node.Right, sn, lbl)
    return
  }
  template := this.intraSentenceTemplate(node)
  if template != nil {
    if sn {
      template = &templateT{ matchCat: true, cat: "SN",
                             subCount: template.subCount,
                             left: template,
                             right: &templateT{ matchCat: true, cat: ""} }
    }
    this.templates[lbl] = append(this.templates[lbl], template)
  }
}

func (this *tpT) intraSentenceTemplate(node *Syntax) *templateT {
  if node.IsEmpty() {
    this.err(node, "expected template expression")
    return nil
  }
  return this.possiblyEmptyIntraSentenceTemplate(node)
}

func (this *tpT) possiblyEmptyIntraSentenceTemplate(node *Syntax) *templateT {
  if node == nil {
    return nil
  }
  if node.Cat == "$" {
    symbol := this.symbolTable[node.Lit]
    if symbol == nil {
      this.err(node, "undeclared symbol: '%s'", node.Lit)
      return nil
    }
    switch symbol.typ {
    case spec_SequenceLabel:
      this.err(node, "nested expression cannot be labeled with sequence label: '%s'", node.Lit)
    case spec_SentenceLabel:
      this.err(node, "nested expression cannot be labeled with sentence label: '%s'", node.Lit)
    case spec_Label:
      return &templateT{ lbl: node.Lit, subCount: 1 }
    case spec_Literal:
      return &templateT{ matchCat: true, cat: symbol.cat, matchLit: true, lit: symbol.lit }
    case spec_Category:
      return &templateT{ matchCat: true, cat: symbol.cat }
    default:
      this.err(node, "unknown symbol: '%s'", node.Lit) // <-- defensive
    }
    return nil
  }
  if node.IsZeroaryOp("<empty") {
    return &templateT{ matchCat: true, cat: "" }
  }
  if symbol := this.symbolTable[node.Lit]; symbol != nil && symbol.typ == spec_ShorthandOperator {
    litSet := make(map[string]bool, len(symbol.ops))
    for _, op := range symbol.ops {
      litSet[op] = true
    }
    template := &templateT{ matchCat: true, cat: "OP", matchLit: true, litSet: litSet, lit: symbol.symb,
                            left: this.possiblyEmptyIntraSentenceTemplate(node.Left),
                            right: this.possiblyEmptyIntraSentenceTemplate(node.Right) }
    template.subCount = template.left.subCountOrZero() + template.right.subCountOrZero()
    return template         
  }
  template := &templateT{ matchCat: true, cat: node.Cat, matchLit: true, lit: node.Lit,
                          left: this.possiblyEmptyIntraSentenceTemplate(node.Left),
                          right: this.possiblyEmptyIntraSentenceTemplate(node.Right) }
  template.subCount = template.left.subCountOrZero() + template.right.subCountOrZero()
  return template
}

func (this *tpT) multiSentenceTemplate(node *Syntax) *templateT {
  if node.IsEmpty() {
    this.err(node, "expected template expression")
    return nil
  }
  return this.possiblyEmptyMultiSentenceTemplate(node)
}

func (this *tpT) possiblyEmptyMultiSentenceTemplate(node *Syntax) *templateT {
  if node == nil {
    return nil
  }
  if node.Cat == "SQ" {
    if node.Left.Cat == "SN" {
      if node.Left.Left.Cat == "$" {
        lbl := node.Left.Left.Lit
        symbol := this.symbolTable[lbl]
        if symbol != nil && symbol.typ == spec_SequenceLabel {
          return &templateT{ lbl: lbl, subCount: 1 }
        }
      } else if node.Left.Left.IsZeroaryOp("<empty") {
        return &templateT{ matchCat: true, cat: "" }
      }
    }
    template := &templateT{ matchCat: true, cat: "SQ",
                            left: this.possiblyEmptyMultiSentenceTemplate(node.Left),
                            right: this.possiblyEmptyMultiSentenceTemplate(node.Right) }
    template.subCount = template.left.subCountOrZero() + template.right.subCountOrZero()
    return template
  }
  if node.Cat == "SN" {
    if node.Left.Cat == "$" {
      lbl := node.Left.Lit
      symbol := this.symbolTable[lbl]
      if symbol != nil && symbol.typ == spec_SentenceLabel {
        return &templateT{ lbl: lbl, subCount: 1 }
      }
    }
    template := &templateT{ matchCat: true, cat: "SN",
                            left: this.possiblyEmptyIntraSentenceTemplate(node.Left),
                            right: this.possiblyEmptyMultiSentenceTemplate(node.Right) }
    template.subCount = template.left.subCountOrZero() + template.right.subCountOrZero()
    return template
  }
  return this.possiblyEmptyIntraSentenceTemplate(node)
}