package dusl

import (
  "io"
  "bytes"
)

type simpleDumper interface {
  Dump(out io.Writer)
}

type prettyDumper interface {
  Dump(out io.Writer, pretty bool)
}

type prefixDumper interface {
  Dump(out io.Writer, prfx string)
}

type prefixPrettyDumper interface {
  Dump(out io.Writer, prfx string, pretty bool)
}

func dumpToString(thing interface{}) string {
  buf := new(bytes.Buffer)
  if dumper, ok := thing.(simpleDumper); ok {
    dumper.Dump(buf)
  } else if dumper, ok := thing.(prettyDumper); ok {
    dumper.Dump(buf, true)
  } else if dumper, ok := thing.(prefixDumper); ok {
    dumper.Dump(buf, "")
  } else if dumper, ok := thing.(prefixPrettyDumper); ok {
    dumper.Dump(buf, "", true)
  }
  return buf.String()
}

