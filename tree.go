package main

import (
	"strings"
)

type Tree struct {
	RootM   *Module
	Modules map[string]*Module
}

func (t *Tree) AddModule(name string) *Module {
	if m, ok := t.Modules[name]; ok {
		return m
	} else {
		t.Modules[name] = NewModule(name)
		return t.Modules[name]
	}
}

func (t *Tree) Graph() {
	ctx := NewContext()
	t.RootM.FindDep(ctx, "", true)
}

func (t *Tree) Dep(name string, tail bool) {
	ctx := NewContext()
	t.RootM.FindDep(ctx, name, tail)
}

func (t *Tree) Ref(name string) {
	ctx := NewContext()
	for n, m := range t.Modules {
		if name == "" || strings.Contains(n, name) {
			m.FindRef(ctx)
		}
	}
}

func NewTree(content string) *Tree {
	t := &Tree{
		Modules: make(map[string]*Module),
	}
	lines := strings.Split(content, "\n")
	for _, line := range lines {
		ms := strings.SplitN(line, " ", 2)
		if len(ms) != 2 {
			continue
		}
		m1 := t.AddModule(ms[0])
		m2 := t.AddModule(ms[1])
		if m1 == nil || m2 == nil {
			continue
		}
		m1.Dep(m2)
		if m1.Version == "" {
			t.RootM = m1
		} else if m2.Version == "" {
			t.RootM = m2
		}
	}

	return t
}

type Module struct {
	Name    string
	Package string
	Version string

	Deps []*Module
	Refs []*Module
}

func (m *Module) Dep(another *Module) {
	m.Deps = append(m.Deps, another)
	another.Refs = append(another.Refs, m)
}

func (m *Module) FindDep(ctx *Context, name string, tail bool) {
	ctx.In(m.Name)
	defer ctx.Out(name == "")

	if name != "" && strings.Contains(m.Name, name) {
		ctx.Focus()
		if !tail {
			return
		}
	}

	for _, d := range m.Deps {
		d.FindDep(ctx, name, tail)
	}
}

func (m *Module) FindRef(ctx *Context) {
	ctx.In(m.Name)
	defer ctx.Out(true)

	for _, d := range m.Refs {
		d.FindRef(ctx)
	}
}

func NewModule(name string) *Module {
	parts := strings.SplitN(name, "@", 2)
	m := &Module{
		Name:    name,
		Package: parts[0],
	}
	if len(parts) > 1 {
		m.Version = parts[1]
	}
	if m.Name == "" {
		return nil
	}
	return m
}
