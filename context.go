package main

import (
	"fmt"
	"strings"
)

type Context struct {
	waiting  []string
	preIdx   int
	depth    int
	focusIdx int
}

func (ctx *Context) In(module string) {
	ctx.waiting = append(ctx.waiting, module)
	ctx.depth += 1
}

func (ctx *Context) Focus() {
	ctx.focusIdx = len(ctx.waiting) - 1
}

func (ctx *Context) Out(all bool) {
	ctx.depth -= 1
	if ctx.focusIdx >= 0 || all {
		for i, pkg := range ctx.waiting {
			if i == ctx.focusIdx {
				pkg = fmt.Sprintf(">> %s <<", pkg)
				ctx.focusIdx = -1
			}
			fmt.Printf("%s %s\n", strings.Repeat(" | ", i+ctx.preIdx), pkg)
		}
		ctx.waiting = []string{}
		ctx.preIdx = ctx.depth
	} else if len(ctx.waiting) > 0 {
		newWaiting := make([]string, 0, cap(ctx.waiting))
		newWaiting = append(newWaiting, ctx.waiting[:len(ctx.waiting)-1]...)
		ctx.waiting = newWaiting
	}
}

func NewContext() *Context {
	return &Context{
		waiting:  []string{},
		focusIdx: -1,
	}
}
