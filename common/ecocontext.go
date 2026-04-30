package common

import (
	"context"
	"io"
)

type EcoContext struct {
	context.Context
	depth  int
	Logger *ecoLogger
}

func NewEcoContext(w io.Writer) *EcoContext {
	var ctx = &EcoContext{
		Context: context.Background(),
		depth:   0,
	}
	ctx.Logger = newEcoLogger(w, ctx)

	return ctx
}

func (ctx *EcoContext) Dec() *EcoContext {
	ctx.depth--
	return ctx
}

func (ctx *EcoContext) Depth() int {
	return ctx.depth
}

func (ctx *EcoContext) Inc() *EcoContext {
	ctx.depth++
	return ctx
}

// func (ctx *ecoContext) indent() string {
// 	const tab = "  "
// 	return strings.Repeat(tab, ctx.level)
// }

// func (ctx *ecoContext) printf(format string, a ...any) (int, error) {
// 	format = fmt.Sprintf("%s%s", ctx.indent(), format)
// 	return fmt.Printf(format, a...)
// }

// func (ctx *ecoContext) println(a ...any) (int, error) {
// 	args := fmt.Sprintln(a...)
// 	return fmt.Printf("%s%s", ctx.indent(), args)
// }
