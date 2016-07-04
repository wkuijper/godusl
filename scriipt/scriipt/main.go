package main

import(
  "os"
  "fmt"
  "io/ioutil"
  "dusl/scriipt"
  "dusl"
)

func main() {
  err := doIt(os.Args[1:])
  if err != nil {
    fmt.Fprintf(os.Stderr, "%s\n", err.Error())
  }
}

func doIt(args []string) error {
  if len(args) < 2 || len(args) > 2 {
    return fmt.Errorf("usage: scriipt <verb> <filepath>\nverbs: undent, undent-raw, tokenize, tokenize-raw, sparse, sparse-raw, trace, trace-raw, parse, parse-raw, run")
  }
  verb, path := args[0], args[1]
  text, err := ioutil.ReadFile(path)
  if err != nil {
    return err
  }
  src := &dusl.Source{ Path: path, Text : text }
  switch verb {
  case "undent", "undent-pretty":
    undent(src, true)
  case "tokenize", "tokenize-pretty":
    tokenize(src, true)
  case "sparse", "sparse-pretty":
    sparse(src, true)
  case "trace", "trace-pretty":
    trace(src, true)
  case "undent-raw":
    undent(src, false)
  case "tokenize-raw":
    tokenize(src, false)
  case "sparse-raw":
    sparse(src, false)
  case "trace-raw":
    trace(src, false)
  case "parse", "parse-raw", "run":
    pgm, err := scriipt.Parse(src)
    if err != nil {
      return err
    }
    if verb == "parse" {
      pgm.Dump(os.Stdout, "", true)
      break
    } else if verb == "parse-raw" {
      pgm.Dump(os.Stdout, "", false)
      break
    } else if verb == "run" { // verb == "run"
      pgm.Run(make(map[string]interface{}, 1024))
      break
    }
    fallthrough
  default:
    return fmt.Errorf("unknown verb: %s: run without arguments to get list of valid verbs", verb)
  }
  return nil
}

func undent(src *dusl.Source, pretty bool) {
  syn := dusl.Undent(src)
  syn.Dump(os.Stdout, "", pretty)
}

func tokenize(src *dusl.Source, pretty bool) {
  syn := scriipt.Lang.Tokenizer().TokenizeUndent(src)
  syn.Dump(os.Stdout, "", pretty)
}

func sparse(src *dusl.Source, pretty bool) {
  syn := scriipt.Lang.Sparser().SparseUndent(src)
  syn.Dump(os.Stdout, "", pretty)
}

func trace(src *dusl.Source, pretty bool) {
  trace := scriipt.Lang.Tracer().TraceUndent(src, "Stmt")
  trace.Dump(os.Stdout, "", pretty)
}
