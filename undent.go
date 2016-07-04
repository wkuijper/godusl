package dusl

import (
  "fmt"
)

// Undent transforms a source text into a syntax tree of unparsed sentences.
func Undent(src *Source) *Syntax {
  ambit := src.FullAmbit()
  return undent(ambit)
}

func undent(ambit *Ambit) *Syntax {
  for !ambit.IsEmpty() {
    indentedLineAmbit, remainderAmbit := ambit.SplitLine()
    lineIndent, lineAmbit := indentedLineAmbit.StripIndent()
    if lineAmbit.IsWhitespace() || lineAmbit.FirstByteIs('#') {
      ambit = remainderAmbit
      continue
    }
    margin := lineIndent
    if margin % 2 != 0 {
      return &Syntax{ Cat: "ERR", Err: fmt.Sprintf("first line indented with odd number of spaces"), Ambit: lineAmbit }
    }
    root, _ := undentSequence(margin, 0, ambit)
    return root
  }
  return &Syntax{ Ambit: ambit }
}

func undentSequence(margin int,
                    currIndent int,
                    ambit *Ambit) (*Syntax, *Ambit) {
  for !ambit.IsEmpty() {
    indentedLineAmbit, remainderAmbit := ambit.SplitLine()
    lineIndent, lineAmbit := indentedLineAmbit.StripIndent()
    if lineAmbit.IsWhitespace() || lineAmbit.FirstByteIs('#') {
      ambit = remainderAmbit
      continue
    }
    lineIndent -= margin
    if lineIndent >= 0 && lineIndent < currIndent {
      return &Syntax{ Ambit: ambit.CollapseLeft() }, ambit
    }
    var head *Syntax
    if lineIndent < 0 {
      head = &Syntax{ Cat: "ERR", Err: fmt.Sprintf("line indented %d spaces before source margin", -lineIndent), Ambit: lineAmbit }
    } else if lineIndent == currIndent {
      head, remainderAmbit = undentSentence(margin, currIndent, lineAmbit, remainderAmbit)
    } else if lineIndent == currIndent+1 {
      head = &Syntax{ Cat: "ERR", Err: fmt.Sprintf("line indented with odd number of spaces"), Ambit: lineAmbit }
    } else if lineIndent >= currIndent+2 && lineIndent < currIndent+5 {
      head = &Syntax{ Cat: "ERR", Err: fmt.Sprintf("line indented more than 2 and less than 5 spaces with respect to previous line: indent 2 spaces for sub-block: indent 5 spaces or more for continuing previous line"), Ambit: lineAmbit }
    } else if lineIndent >= currIndent + 5 {
      head = &Syntax{ Cat: "ERR", Err: fmt.Sprintf("line continuation not possible here: indent less than 5 spaces with respect to previous line"), Ambit: lineAmbit }
    }
    var tail *Syntax
    tail, remainderAmbit = undentSequence(margin, currIndent, remainderAmbit)
    return &Syntax{ Cat: "SQ", Ambit: head.Ambit.Merge(tail.Ambit),
                    Left: head,
                    Right: tail }, remainderAmbit
  }
  return &Syntax{ Ambit: ambit }, ambit
}

func undentSentence(margin int,
                    currIndent int,
                    firstLineAmbit *Ambit,
                    ambit *Ambit) (*Syntax, *Ambit) {
  sentenceAmbit := firstLineAmbit
  for !ambit.IsEmpty() {
    indentedLineAmbit, remainderAmbit := ambit.SplitLine()
    lineIndent, lineAmbit := indentedLineAmbit.StripIndent()
    if lineAmbit.IsWhitespace()|| lineAmbit.FirstByteIs('#') {
      sentenceAmbit = sentenceAmbit.Merge(lineAmbit)
      ambit = remainderAmbit
      continue
    }
    if lineIndent < margin + currIndent + 5 {
      break
    }  
    sentenceAmbit = sentenceAmbit.Merge(lineAmbit)
    ambit = remainderAmbit
  }
  var subSequence *Syntax
  subSequence, ambit =  undentSequence(margin, currIndent+2, ambit)
  return &Syntax{ Cat: "SN", Ambit: sentenceAmbit.Merge(subSequence.Ambit),
                  Left: &Syntax{ Cat: "UN", Ambit: sentenceAmbit },
                  Right: subSequence }, ambit
}