package scriipt

func (this *noopSimple) Run(st map[string]interface{}) {
  // NOOP
}

func (this *assignSimple) Run(st map[string]interface{}) {
  st[this.varName] = this.expr.Eval(st)
}

func (this *incrSimple) Run(st map[string]interface{}) {
  st[this.varName] = st[this.varName].(int) + 1
}

func (this *decrSimple) Run(st map[string]interface{}) {
  st[this.varName] = st[this.varName].(int) - 1
}

func (this *exprSimple) Run(st map[string]interface{}) {
  this.expr.Eval(st)
}

func (this *seqStmt) Run(st map[string]interface{}) {
  this.head.Run(st)
  this.tail.Run(st)
}

func (this *emptyStmt) Run(st map[string]interface{}) {
  // NOOP
}

func (this *forStmt) Run(st map[string]interface{}) {
  for this.initial.Run(st); this.cond.Eval(st).(bool); this.update.Run(st) {
    this.body.Run(st)
  }
}

func (this *ifStmt) Run(st map[string]interface{}) {
  if this.cond.Eval(st).(bool) {
    this.then.Run(st)
  } else {
    this.el5e.Run(st)
  }
}

func (this *simpleStmt) Run(st map[string]interface{}) {
  this.simple.Run(st)
}
