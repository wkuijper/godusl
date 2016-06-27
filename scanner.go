package dusl

import (
  "strings"
)

type Scanner interface {
  Scan() Scan
}

type Scan interface {
  Consume(r rune) (string, bool)
  Reset()
}

type emptyScanner struct {}

func (this *emptyScanner) Scan() Scan {
  return this
}

func (this *emptyScanner) Consume(r rune) (string, bool) {
  return "", false
}

func (this *emptyScanner) Reset() {
  // NOOP
}

type seqScanner struct {
  master Scanner
  slave Scanner
}

type seqScan struct {
  master Scan
  slave Scan
}

func filterNilScanners(scanners []Scanner) []Scanner {
  if len(scanners) == 0 {
    return scanners
  }
  l := len(scanners)-1
  if scanners[l] == nil {
    return filterNilScanners(scanners[:l])
  }
  return append(filterNilScanners(scanners[:l]), scanners[l])
}

func sequenceScanners(scanners ...Scanner) Scanner {
  scanners = filterNilScanners(scanners)
  if len(scanners) == 0 {
    return nil
  }
  if len(scanners) == 1 {
    return scanners[0]
  }
  return &seqScanner{ master: scanners[0], slave: sequenceScanners(scanners[1:]...) }
}

func (this *seqScanner) Scan() Scan {
  return &seqScan{ master: this.master.Scan(), slave: this.slave.Scan() }
}

func (this *seqScan) Consume(r rune) (string, bool) {
  masterCat, masterCont := this.master.Consume(r)
  slaveCat, slaveCont := this.slave.Consume(r)
  if masterCat == "" {
    return slaveCat, masterCont || slaveCont
  }
  return masterCat, masterCont || slaveCont
}

func (this *seqScan) Reset() {
  this.master.Reset()
  this.slave.Reset()
}

func composeScanners(scanners ...Scanner) Scanner {
  scanners = filterNilScanners(scanners)
  if len(scanners) == 0 {
    return nil
  }
  if len(scanners) == 1 {
    return scanners[0]
  }
  return &compScanner{ scannerA: scanners[0], scannerB: composeScanners(scanners[1:]...) }
}

type compScanner struct {
  scannerA Scanner
  scannerB Scanner
}

type compScan struct {
  scan1 Scan
  scan2 Scan
}

func (this *compScanner) Scan() Scan {
  return &compScan{ this.scannerA.Scan(), this.scannerB.Scan() }
}

func (this *compScan) Consume(r rune) (string, bool) {
  cat1, cont1 := this.scan1.Consume(r)
  cat2, cont2 := this.scan2.Consume(r)
  if cat2 == "" || cat1 == cat2 {
    return cat1, cont1 || cont2
  }
  if cat1 == "" {
    return cat2, cont1 || cont2
  }
  return "", cont1 || cont2
}

func (this *compScan) Reset() {
  this.scan1.Reset()
  this.scan2.Reset()
}

// The SimpleStringScanner scans for double quoted string literals, 
// It recognizes escape sequences backslash-doublequote,
// backslash-backslash, backslash-n(ewline), backslash-r(eturn) and backslash-t(ab).
// It also recognizes raw string literals that start with a backtick and end at the
// end-of-line, in that case nothing needs to be escaped (not even the backtick itself).
// This scanner only reports the lexical category: "STR".
var SimpleStringScanner Scanner

type simpleStringScanner struct {}

func (this *simpleStringScanner) Scan() Scan {
  return &simpleStringScan{}
}

type simpleStringScan struct {
  state int  
}

func (this *simpleStringScan) Consume(r rune) (string, bool) {
  const (
    INIT = iota // convention requires: INIT == 0
    RAW
    RAW_R
    INSIDE
    ESCAPE
    NOMORE
  )
  switch (this.state) {
  case INIT:
    if r == '"' {
      this.state = INSIDE
      return "", true
    } else if r == '`' {
      this.state = RAW
      return "", true
    } else {
      this.state = NOMORE
      return "", false
    }
  case INSIDE:
    if r == '\\' {
      this.state = ESCAPE
      return "", true
    } else if r == '"' {
      this.state = NOMORE
      return "STR", false
    } else {
      return "", true
    }
  case ESCAPE:
    switch r {
    case 'n', 'r', 't', '"':
      this.state = INSIDE
      return "", true
    default:
      return "", false
    }
  case RAW:
    if r == '\r' {
      this.state = RAW_R
      return "STR", true
    } else if r == '\n' {
      this.state = NOMORE
      return "STR", false
    } else {
      return "", true
    }
  case RAW_R:
    if r == '\n' {
      this.state = NOMORE
      return "STR", false
    } else {
      return "", false
    }
  }
  return "", false 
}

func (this *simpleStringScan) Reset() {
  this.state = 0
}

// SimpleDecimalNumScanner scans numbers in decimal notation. 
// No fractions, exponents, base 2, 8 or 16 notations are supported by this simple scanner.
// This scanner reports only the lexical category "NUM"
var SimpleDecimalNumScanner Scanner

type simpleDecimalNumScanner struct {}

func (this *simpleDecimalNumScanner) Scan() Scan {
  return &simpleDecimalNumScan{}
}

type simpleDecimalNumScan struct {
  state int
}

func (this *simpleDecimalNumScan) Consume(r rune) (string, bool) {
  const (
    INIT = iota // convention requires: INIT == 0
    REST
    NOMORE
  )
  switch this.state {
  case INIT:
    if r >= '1' && r <= '9' {
      this.state = REST
      return "NUM", true
    } else if r == '0' {
      this.state = NOMORE
      return "NUM", false
    } else {
      this.state = NOMORE
      return "", false
    }
  case REST:
    if r >= '0' && r <= '9' {
      return "NUM", true
    } else {
      this.state = NOMORE
      return "", false
    }
  }
  return "", false
}

func (this *simpleDecimalNumScan) Reset() {
  this.state = 0
}

// SimpleIdentifierScanner scans simple ASCII based identifiers
// starting with a '_' or roman letter and containing only '_' or roman alphanumerics.
// This scanner reports only the lexical category "ID".
var SimpleIdentifierScanner Scanner

type simpleIdentifierScanner struct {}

func (this *simpleIdentifierScanner) Scan() Scan {
  return &simpleIdentifierScan{}
}

type simpleIdentifierScan struct {
  state int
}  

func (this *simpleIdentifierScan) Consume(r rune) (string, bool) {
  const (
    INIT = iota // convention requires: INIT == 0
    REST
    NOMORE
  )
  switch (this.state) {
  case INIT:
    if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') {
      this.state = REST
      return "ID", true
    } else {
      this.state = NOMORE
      return "", false
    }
  case REST:
    if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') || r == '_' {
      return "ID", true
    } else {
      this.state = NOMORE
      return "", false
    }
  }
  return "", false
}

func (this *simpleIdentifierScan) Reset() {
  this.state = 0
}

// SimpleBaseScanner recognizes ASCII whitespace,
// and comments starting with a hash-sign, ending at end-of-line.
// The simplebase scanner returns lexical category "WS" (for whitespace) 
var SimpleBaseScanner Scanner

type simpleBaseScanner struct {}

func (this *simpleBaseScanner) Scan() Scan {
  return &simpleBaseScan{}
}

type simpleBaseScan struct {
  state int
}  

func (this *simpleBaseScan) Consume(r rune) (string, bool) {
  const (
    INIT = iota // convention requires: INIT == 0
    COMMENT
    COMMENT_R
    NOMORE
  )
  switch (this.state) {
  case INIT:
    if r == ' ' || r == '\t' || r == '\r' || r == '\n' {
      return "WS", true
    } else if r == '#' {
      this.state = COMMENT
      return "WS", true
    } else {
      this.state = NOMORE
      return "", false
    }
  case COMMENT:
    if r == '\r' {
      this.state = COMMENT_R
      return "WS", true
    } else if r == '\n' {
      this.state = INIT
      return "WS", true
    } else {
      return "WS", true
    }
  case COMMENT_R:
    if r == '\n' {
      this.state = INIT
      return "WS", true
    } else {
      this.state = NOMORE
      return "", false
    }
  }
  return "", false
}

func (this *simpleBaseScan) Reset() {
  this.state = 0
}

// A PrefixScanner works as a prefixtree: for a set of fixed length tokens 
// it will return the associated lexical category.
// The scanner is defined by giving a set of descriptions, each description
// consists of a category name followed by a space followed by a list of
// token literals that are to be recognized (also separated by single spaces).
// For example: PrefixScanner("ID a b c", "NUM 1 2 3", "OP + - *")
func PrefixScanner(descs ...string) Scanner {
  prfxTree := &prfxTree{}
  for _, desc := range descs {
    parts := strings.Split(desc, " ")
    if len(parts) < 1 {
      continue
    }
    cat := parts[0]
    for _, part := range parts[1:] {
      prfxTree.add(cat, part)
    }
  }
  return prfxTree
}

// The DefaultScanner is composed of the SimpleBaseScanner, 
// the SimpleStringScanner, the SimpleIdentifierScanner and the 
// SimpleDecimalNumScanner.
// As such it reports the lexical categories: "WS", "STR", "ID", and "NUM".
var DefaultScanner Scanner

func init() {
  SimpleStringScanner = &simpleStringScanner{} 
  SimpleDecimalNumScanner = &simpleDecimalNumScanner{}
  SimpleIdentifierScanner = &simpleIdentifierScanner{}
  SimpleBaseScanner = &simpleBaseScanner{}
  DefaultScanner = composeScanners(SimpleBaseScanner, 
                                   SimpleStringScanner, 
                                   SimpleIdentifierScanner,
                                   SimpleDecimalNumScanner)
}

type metaSymbolScannerT struct {}

var metaSymbolScanner = &metaSymbolScannerT{}

func (this *metaSymbolScannerT) Scan() Scan {
  return &metaSymbolScan{}
}

type metaSymbolScan struct {
  state int
}

func (this *metaSymbolScan) Consume(r rune) (string, bool) {
  const (
    INIT = iota // convention requires: INIT == 0
    REST
    NOMORE
  )
  switch (this.state) {
  case INIT:
    if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') {
      this.state = REST
      return "$", true
    } else {
      this.state = NOMORE
      return "", false
    }
  case REST:
    if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') || r == '_' {
      return "$", true
    } else {
      this.state = NOMORE
      return "", false
    }
  }
  return "", false
}

func (this *metaSymbolScan) Reset() {
  this.state = 0
}
