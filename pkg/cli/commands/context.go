package commands

// ExecutionContext holds variables for script execution, especially the output variable.
type ExecutionContext struct {
	vars map[string]interface{}
}

func NewExecutionContext() *ExecutionContext {
	return &ExecutionContext{
		vars: make(map[string]interface{}),
	}
}

func (ctx *ExecutionContext) SetOutput(value interface{}) {
	ctx.vars["output"] = value
}

func (ctx *ExecutionContext) GetOutput() interface{} {
	return ctx.vars["output"]
}

func (ctx *ExecutionContext) SetVar(name string, value interface{}) {
	ctx.vars[name] = value
}

func (ctx *ExecutionContext) GetVar(name string) interface{} {
	return ctx.vars[name]
}

func (ctx *ExecutionContext) Vars() map[string]interface{} {
	return ctx.vars
}
