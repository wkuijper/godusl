package scriipt

import (
  "fmt"
)

func (this *idExpr) Eval(st map[string]interface{}) interface{} {
  return st[this.idName]
}

func (this *numExpr) Eval(st map[string]interface{}) interface{} {
  return this.num
}

func (this *strExpr) Eval(st map[string]interface{}) interface{} {
  return this.unescapedStr
}

func (this *prefixExpr) Eval(st map[string]interface{}) interface{} {
  subVal := this.sub.Eval(st)
  switch this.lit {
  case "+":
    return + subVal.(int)
  case "-":
    return - subVal.(int)
  case "!":
    return ! subVal.(bool)
  default:
    panic("unimplemented prefix op: " + this.lit)
  }
}

func (this *postfixExpr) Eval(st map[string]interface{}) interface{} {
  //subVal := this.sub.Eval(st)
  switch this.lit {
  default:
    panic("unimplemented postfix op: " + this.lit)
  }
}

func (this *infixExpr) Eval(st map[string]interface{}) interface{} {
  leftVal := this.left.Eval(st)
  rightVal := this.right.Eval(st)
  switch this.lit {
  case "+":
    return leftVal.(int) + rightVal.(int)
  case "-":
    return leftVal.(int) - rightVal.(int)
  case "*":
    return leftVal.(int) * rightVal.(int)
  case "/":
    return leftVal.(int) / rightVal.(int)
  case ">":
    return leftVal.(int) > rightVal.(int)
  case "<":
    return leftVal.(int) < rightVal.(int)
  case ">=":
    return leftVal.(int) >= rightVal.(int)
  case "<=":
    return leftVal.(int) <= rightVal.(int)
  case "==":
    return leftVal.(int) == rightVal.(int)
  case "!=":
    return leftVal.(int) != rightVal.(int)
  case "&&":
    return leftVal.(bool) && rightVal.(bool)
  case "||":
    return leftVal.(bool) || rightVal.(bool)
  default:
    panic("unimplemented infix op: " + this.lit)
  }
}

func (this *trueExpr) Eval(st map[string]interface{}) interface{} {
  return true
}

func (this *falseExpr) Eval(st map[string]interface{}) interface{} {
  return false
}

func (this *callExpr) Eval(st map[string]interface{}) interface{} {
  values := make([]interface{}, len(this.args))
  for i, arg := range this.args {
    values[i] = arg.Eval(st)
  }
  switch this.funcName {
  case "print":
    fmt.Println(values...)
    return nil
  case "max":
    if len(values) == 2 {
      if values[0].(int) > values[1].(int) {
        return values[0]
      } else {
        return values[1]
      }
    }
  case "min":
    if len(values) == 2 {
      if values[0].(int) < values[1].(int) {
        return values[0]
      } else {
        return values[1]
      }
    }
  case "panic":
    if len(values) == 1 {
      panic(values[0].(string))
    }
  }
  panic(fmt.Sprintf("%s: unknown function: %s/%d", this.Location(), this.funcName, len(values)))
}



