# DUSL

Domain UnSpecific Language.

Explore the included example:

    $ export GOPATH=...
    $ cd .../src/dusl/scriipt/scriipt
    $ go build
    $ ./scriipt
    usage: scriipt <verb> <filepath>
    verbs: undent, undent-raw, tokenize, tokenize-raw, span, span-raw, sparse, sparse-raw, trace, trace-raw, parse, parse-raw, run

Run the toy scripting language:

    $ cat fac.scriipt
    x = 10 # compute factorial of x

    z = 1
    for y = 1; y <= x; y++
      z = z * y
      if z > 10000000
        panic("too big!")

    print("fac", x, "==", z)
    $ ./scriipt run fac.scriipt
    fac 10 == 3628800

See the prettified parse tree:

    $ ./scriipt parse fac.scriipt
    x = 10
    z = 1
    for y = 1; (y <= x); y++
      z = (z * y)
      if (z > 10000000)
        panic("too big!")
    print("fac", x, "==", z)

See the raw parse tree:

    $ ./scriipt parse-raw fac.scriipt
    SeqStmt:fac.scriipt[0:140]
      SimpleStmt:fac.scriipt[0:140]
        AssignSimple:x:fac.scriipt[0:6]
          NumExpr:10:fac.scriipt[4:6]
      SeqStmt:fac.scriipt[33:140]
        SimpleStmt:fac.scriipt[33:140]
          AssignSimple:z:fac.scriipt[33:38]
            NumExpr:1:fac.scriipt[37:38]
        SeqStmt:fac.scriipt[39:140]
          ForStmt:fac.scriipt[39:140]
            AssignSimple:y:fac.scriipt[43:48]
              NumExpr:1:fac.scriipt[47:48]
            InfixExpr:<=:fac.scriipt[50:56]
              IdExpr:y:fac.scriipt[50:51]
              IdExpr:x:fac.scriipt[55:56]
            IncrSimple:y:fac.scriipt[58:61]
            SeqStmt:fac.scriipt[64:115]
              SimpleStmt:fac.scriipt[64:115]
                AssignSimple:z:fac.scriipt[64:73]
                  InfixExpr:*:fac.scriipt[68:73]
                    IdExpr:z:fac.scriipt[68:69]
                    IdExpr:y:fac.scriipt[72:73]
              SeqStmt:fac.scriipt[76:115]
                IfStmt:fac.scriipt[76:115]
                  InfixExpr:>:fac.scriipt[79:91]
                    IdExpr:z:fac.scriipt[79:80]
                    NumExpr:10000000:fac.scriipt[83:91]
                  SeqStmt:fac.scriipt[96:115]
                    SimpleStmt:fac.scriipt[96:115]
                      ExprSimple:fac.scriipt[96:113]
                        CallExpr:panic:fac.scriipt[96:113]
                          StrExpr:"too big!":fac.scriipt[102:112]
                    EmptyStmt:fac.scriipt[115:115]
                  EmptyStmt:fac.scriipt[115:115]
                EmptyStmt:fac.scriipt[115:115]
          SeqStmt:fac.scriipt[115:140]
            SimpleStmt:fac.scriipt[115:140]
              ExprSimple:fac.scriipt[115:139]
                CallExpr:print:fac.scriipt[115:139]    
                  StrExpr:"fac":fac.scriipt[121:126]
                  IdExpr:x:fac.scriipt[128:129]
                  StrExpr:"==":fac.scriipt[131:135]
                  IdExpr:z:fac.scriipt[137:138]
            EmptyStmt:fac.scriipt[140:140]

Look at the grammar:

    $ cat ../lang.go
    ...
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
        Stmt                    # continuation
      or>
        for loopHeader          # optional: initialization; condition; update
          Stmt                  #   body
        end                     # optional: end-marker
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


Explore the intermediate stages on your own.

Look at the code for details: start with `scriipt/*` as an example of
usage, then move on to the `dusl` package proper for the internal
workings.