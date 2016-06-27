package dusl

import (
  "bytes"
  "fmt"
)

// An AmbitError reports a syntax error. It consists of an Ambit
// which encodes the location or region within the source to which the error applies
// and the descriptive error message to be shown back to the user.
// The Error() method is memoized so it is computed only once in a lazy fashion.
func AmbitError(ambit *Ambit, msg string) error {
  return &ambitError{ ambit: ambit, msg: msg }
}

type ambitError struct {
  ambit *Ambit
  msg string
  formattedMsg string
}

func (this *ambitError) Error() string {
  if this.formattedMsg == "" {
    this.formattedMsg = fmt.Sprintf("%s: %s", this.ambit.Location(), this.msg)
  }
  return this.formattedMsg
}

// A SummaryError summarizes a whole bunch of errors into one. The N field can be
// positive in which case it represent the cut-off value, or negative in which case
// all the errors will be reported. The Error() method is memoized so it is
// computed only once in a lazy fashion.
func SummaryError(errs []error, n int) error {
  if len(errs) == 0 || n == 0 {
    return nil
  }
  return &summaryError{ errs: errs, n: n }
}

type summaryError struct {
  errs []error
  n int
  formattedMsg string
}

func (this *summaryError) Error() string {
  if this.formattedMsg == "" {
    errs := this.errs
    n := this.n
    if len(errs) == 0 {
      return "" // <-- defensive
    }
    buf := new(bytes.Buffer)
    for index, err := range errs {
      if n < 0 || index < n {
        fmt.Fprintf(buf, "%s\n", err.Error())
      } else if index == n {
        fmt.Fprintf(buf, "and %d more error(s)\n", len(errs)-n)
      } else {
        break
      }
    }
    this.formattedMsg = buf.String()
  }
  return this.formattedMsg
}
