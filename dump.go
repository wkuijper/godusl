package dusl

import (
  "io"
  "bytes"
)

type Dumper interface {
  Dump(out io.Writer)
}

type PrettyDumper interface {
  Dump(out io.Writer, pretty bool)
}

type PrefixDumper interface {
  Dump(out io.Writer, prfx string)
}

type PrefixPrettyDumper interface {
  Dump(out io.Writer, prfx string, pretty bool)
}

func Dump2String(thing interface{}) string {
  buf := new(bytes.Buffer)
  if dumper, ok := thing.(Dumper); ok {
    dumper.Dump(buf)
  } else if dumper, ok := thing.(PrettyDumper); ok {
    dumper.Dump(buf, true)
  } else if dumper, ok := thing.(PrefixDumper); ok {
    dumper.Dump(buf, "")
  } else if dumper, ok := thing.(PrefixPrettyDumper); ok {
    dumper.Dump(buf, "", true)
  }
  return buf.String()
}

