package scriipt

import (
  "dusl"
)

var Lang dusl.Lang

func init() {
  var err error
  Lang, err = dusl.NewSpec().
    Lexical(dusl.DefaultScanner).
    Category("ID", "identifier").
    Category("NUM", "number").
    Category("STR", "string").
    OperatorEFE("true", "false").
    OperatorEFA("+", "-", "!").
    OperatorAFE("++", "--").
    OperatorBFA("*", "/").
    OperatorBFA("+", "-").
    OperatorBFA("==", "!=", "<", "<=", ">=", ">").
    OperatorBFA("&&").
    OperatorBFA("||").
    OperatorBFA(",").
    OperatorBFA("=").
    OperatorBFA(";").
    Brackets("( )", "[ ]", "{ }").

    Literal("if", "ID").
    Literal("else", "ID").
    Literal("for", "ID").
    Literal("end", "ID").
    Literal("noop", "ID").
    
    Label("x", "expression").
    Label("args", "argument list").
    Label("aargs", "non-empty argument list").
    
    Label("s", "simple statement").
    
    Label("f", "function name").
    Label("v", "variable name").
    Label("i", "integer variable name").
    
    Label("hdr", "for header").

    SequenceLabel("S", "statement").
    SequenceLabel("E", "else continuation statement").

    ShorthandOperator("prfx~", "-", "+").
    ShorthandOperator("~postfx", "--", "++").
    ShorthandOperator("~infix~", "-", "+", "*", "/", "==", "!=", "<", "<=", ">=", ">", "&&", "||").
    
    Grammar(`
    
      # placeholders
      # ------------
            
      f is> ID                  # function name
      i is> ID                  # integer variable name
      v is> ID                  # variable name
      
      # expressions
      # -----------
      
      x is> ID
            or> NUM
            or> STR
            or> (x)
            or> f(args)
            or> prfx~ x
            or> x ~postfx
            or> x ~infix~ x
            or> true
            or> false
      
      # (argument) lists
      # ----------------
      
      args is> <empty or> aargs
      aargs is> aargs, x or> x

      # simple statements
      # -----------------
      
      s is> noop                # no-operation
            or> v = x           # assignment
            or> i++             # increment
            or> i--             # decrement
            or> f(args)         # function call
            
      # block statements
      # ----------------
      
      S is>
        <empty
      or>
        for hdr                 # optional: initialization; condition; update
          S                     #   body
        S                       # continuation
      or>
        for hdr                 # optional: initialization; condition; update
          S                     #   body
        end                     # optional: end-marker
        S                       # continuation
      or>
        if x                    # condition
          S                     #   body
        E                       # else-if/continuation w/ optional end-marker
      or>
        s                       # simple statement
        S                       # continuation

      # else-if
      # -------
      
      E is>
        else if x               # condition
          S                     #   body
        E                       # else-if/continuation
      or>
        else
          S                     #   body
        S                       # continuation
      or>
        end
        S
      or>
        S

      # for header
      # ----------
      
      hdr is> <empty            # infinite loop
              or> s; x; s       # full for header
              or> x             # while header

  `)
  if err != nil {
    panic(err.Error())
  }
}
