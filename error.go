package dusl

import (
  "bytes"
  "fmt"
  "errors"
)

func errorN(errs []error, n int) error {
  if len(errs) == 0 {
    return nil
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
  return errors.New(buf.String())
}