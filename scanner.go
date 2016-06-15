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

var SimpleStringScanner Scanner = &simpleStringScanner{}

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
    INSIDE
    ESCAPE
    NOMORE
  )
  switch (this.state) {
  case INIT:
    if r == '"' {
      this.state = INSIDE
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
  }
  return "", false 
}

func (this *simpleStringScan) Reset() {
  this.state = 0
}

var SimpleDecimalNumScanner Scanner = &simpleDecimalNumScanner{}

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

var SimpleIdentifierScanner Scanner = &simpleIdentifierScanner{}

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

var SimpleBaseScanner = &simpleBaseScanner{}

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

func PrefixScanner(descs ...string) *prfxTree {
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

func NewDefaultScanner() Scanner {
  return composeScanners(SimpleBaseScanner, 
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
