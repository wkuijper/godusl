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
    JuxtapositionLWA("if", "else", "for").
    
    Literal("end", "noop").
    
    Label("expr", "expression").
    Label("args", "argument list").
    Label("aargs", "non-empty argument list").
    
    Label("simpleStmt", "simple statement").
    
    Label("f", "function name").
    Label("v", "variable name").
    Label("i", "integer variable name").
    
    Label("loopHeader", "loop header").

    SequenceLabel("Stmt", "statement").
    SequenceLabel("ElseCont", "else continuation statement").

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
      
      expr is> ID
           or> NUM
           or> STR
           or> (expr)
           or> f(args)
           or> prfx~ expr
           or> expr ~postfx
           or> expr ~infix~ expr
           or> true
           or> false
      
      # (argument) lists
      # ----------------
      
      args is> <empty or> aargs
      aargs is> aargs, expr or> expr

      # simple statements
      # -----------------
      
      simpleStmt is> noop       # no-operation
                 or> v = expr   # assignment
                 or> i++        # increment
                 or> i--        # decrement
                 or> f(args)    # function call
            
      # block statements
      # ----------------
      
      Stmt is>
        <empty
      or>
        for loopHeader          # optional: initialization; condition; update
          Stmt                  #   body
        end                     # optional: end-marker
        Stmt                    # continuation
      or>
        for loopHeader          # optional: initialization; condition; update
          Stmt                  #   body
        Stmt                    # continuation
      or>
        if expr                 # condition
          Stmt                  #   body
        ElseCont                # else-if/continuation w/ optional end-marker
      or>
        simpleStmt              # simple statement
        Stmt                    # continuation

      # else-if
      # -------
      
      ElseCont is>
        else if expr            # condition
          Stmt                  #   body
        ElseCont                # else-if/continuation
      or>
        else
          Stmt                  #   body
        Stmt                    # continuation
      or>
        end                     # optional end marker
        Stmt
      or>
        Stmt

      # for header
      # ----------
      
      loopHeader is> <empty                         # infinite loop
                 or> simpleStmt; expr; simpleStmt   # full for loop header
                 or> expr                           # while loop header

  `)
  if err != nil {
    panic(err.Error())
  }
}
