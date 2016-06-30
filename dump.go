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

func dumpToString(thing interface{}, pretty bool) string {
  buf := new(bytes.Buffer)
  if dumper, ok := thing.(simpleDumper); ok {
    dumper.Dump(buf)
  } else if dumper, ok := thing.(prettyDumper); ok {
    dumper.Dump(buf, pretty)
  } else if dumper, ok := thing.(prefixDumper); ok {
    dumper.Dump(buf, "")
  } else if dumper, ok := thing.(prefixPrettyDumper); ok {
    dumper.Dump(buf, "", pretty)
  }
  return buf.String()
}

