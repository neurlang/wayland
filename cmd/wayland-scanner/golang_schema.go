package main

import "strings"
import "fmt"

type GoFile struct {
	Structs   []*GoStruct
	Constants []*GoConstant
	Name      string
	PkgName   string
	FileName  string
}

func (f *GoFile) AddStruct(s *GoStruct) {
	f.Structs = append(f.Structs, s)
}

func (f *GoFile) AddConstant(s *GoConstant) {
	f.Constants = append(f.Constants, s)
}

func (f *GoFile) IsNeedMutex() bool {
	for _, strct := range f.Structs {
		if strct.IsNeedMutex() {
			return true
		}
	}
	return false
}

func (f *GoFile) Serialize() string {
	var b strings.Builder
	b.WriteString("// This file is autogenerated from: " + f.FileName + "\n// Do not edit\n\n")
	b.WriteString("// Package " + f.PkgName + " implements the " + f.Name + " protocol\n")
	b.WriteString(`package ` + f.PkgName + "\n\nimport (\n")
	if f.IsNeedMutex() {
		b.WriteString("\t\"sync\"\n")
	}
	b.WriteString("\n)\n")
	for _, k := range f.Constants {
		b.WriteString(k.Serialize())
	}
	for _, s := range f.Structs {
		b.WriteString(s.Serialize())
	}
	return b.String()
}

type GoDispatchMethod struct {
	Events []*GoEvent
}

func (m *GoDispatchMethod) Serialize(s *GoStruct) string {
	var b strings.Builder
	b.WriteString(`// Dispatch dispatches event for object ` + s.Name + `
func (p *` + s.Name + `) Dispatch(event *Event) {
	switch event.Opcode {
`)
	for num, e := range s.Events {
		b.WriteString(e.SerializeCase(num, s))
	}
	b.WriteString("\n\t}\n}\n")
	return b.String()
}

type GoStruct struct {
	Name     string
	Comment  string
	Cons     GoConstructor
	Mets     []*GoMethod
	Events   []*GoEvent
	Handlers []*GoHandler
}

func (f *GoStruct) AddMethod(s *GoMethod) {
	f.Mets = append(f.Mets, s)
}

func (f *GoStruct) AddEvent(s *GoEvent) {
	f.Events = append(f.Events, s)
}

func (f *GoStruct) AddHandler(s *GoHandler) {
	f.Handlers = append(f.Handlers, s)
}

func (s *GoStruct) IsNeedMutex() bool {
	return len(s.Handlers) > 0
}

func (s *GoStruct) Serialize() string {
	var b strings.Builder
	b.WriteString(`// ` + s.Name + ` ` + s.Comment + `
type ` + s.Name + ` struct {
	BaseProxy
`)
	if s.IsNeedMutex() {
		b.WriteString("\tmu sync.RWMutex\n")
	}
	for _, e := range s.Handlers {
		b.WriteString(e.SerializeAttr(s))
	}
	b.WriteString("}\n" + s.Cons.Serialize(s))
	for num, m := range s.Mets {
		b.WriteString(m.Serialize(num, s))
	}
	b.WriteString((&GoDispatchMethod{
		Events: s.Events,
	}).Serialize(s))
	for _, e := range s.Events {
		b.WriteString(e.Serialize())
	}
	for _, e := range s.Handlers {
		b.WriteString(e.Serialize(s))
	}
	return b.String()
}

type GoConstructor struct{}

func (c *GoConstructor) Serialize(s *GoStruct) string {
	return `// New` + s.Name + ` is a constructor for the ` + s.Name + ` object
func New` + s.Name + `(ctx *Context) *` + s.Name + ` {
	ret := new(` + s.Name + `)
	ctx.Register(ret)
	return ret
}
`
}

type GoMethod struct {
	Name    string
	Comment string
	Args    []*GoArg
}

func (f *GoMethod) AddEventArg(s *GoArg, t *GoType) {
	s.Type = t
	f.Args = append(f.Args, s)
	if t.Type == "fd" {
		// Add the error arg for file descriptor here
		var s2 = GoArg{
			Name:    s.Name + "Error",
			Comment: s.Comment + " (error)",
			Type: &GoType{
				Type: "error",
			},
		}
		f.Args = append(f.Args, &s2)
	}
}

func (m *GoMethod) Serialize(num int, s *GoStruct) string {
	var b strings.Builder
	b.WriteString(`// ` + m.Name + " " + m.Comment + "\n" + `func (p *` + s.Name + `) ` + m.Name + `(`)
	var first = true
	for _, arg := range m.Args {
		if arg.IsInput() {
			if first {
				first = false
			} else {
				b.WriteString(", ")
			}
			b.WriteString(arg.Name + " " + arg.Type.Serialize())
		}
	}
	b.WriteString(`) (`)
	for _, arg := range m.Args {
		if arg.IsOutput() {
			b.WriteString(arg.Type.Serialize() + ", ")
		}
	}
	b.WriteString(`error) {
	`)
	for _, arg := range m.Args {
		if arg.IsOutput() {
			b.WriteString("ret" + arg.Name + " := " + arg.Type.SerializeNew())
		}
	}
	b.WriteString(`
	return `)

	for _, arg := range m.Args {
		if arg.IsOutput() {
			b.WriteString("ret" + arg.Name + ", ")
		}
	}

	b.WriteString(`p.Context().SendRequest(p, ` + fmt.Sprintf("%d", num))
	for _, arg := range m.Args {
		if arg.IsOutput() {
			b.WriteString(", ret" + arg.Name)
		}
		if arg.IsInput() {
			b.WriteString(", " + arg.Name)
		}
	}
	b.WriteString(")\n}\n")
	return b.String()
}

type GoEvent struct {
	Name    string
	Comment string
	Args    []*GoArg
	IsCb    bool
	IsPt    bool
	IsBf    bool
}

func (f *GoEvent) AddEventArg(s *GoArg, t *GoType) {
	s.Type = t
	f.Args = append(f.Args, s)
	if t.Type == "fd" {
		// Add the error arg for file descriptor here
		var s2 = GoArg{
			Name:    s.Name + "Error",
			Comment: s.Comment + " (error)",
			Type: &GoType{
				Type: "error",
			},
		}
		f.Args = append(f.Args, &s2)
	}
}

func (m *GoEvent) Serialize() string {
	var b strings.Builder
	b.WriteString(`// ` + m.Name + `Event is the ` + m.Comment + `
type ` + m.Name + `Event struct {
`)
	for _, a := range m.Args {
		b.WriteString(a.Serialize())
	}
	if m.IsCb {
		b.WriteString("\tC *Callback")
	}
	if m.IsPt {
		b.WriteString("\tP *Pointer")
	}
	if m.IsBf {
		b.WriteString("\tB *Buffer")
	}
	b.WriteString("\n}\n")
	return b.String()
}

func (m *GoEvent) SerializeCase(num int, s *GoStruct) string {
	var b strings.Builder
	b.WriteString(`	case ` + fmt.Sprintf("%d", num) + `:
		if len(p.private` + m.Name + `s) > 0 {
			ev := ` + m.Name + `Event{}
`)
	if m.IsCb {
		b.WriteString("\t\t\tev.C = p\n")
	}
	if m.IsPt {
		b.WriteString("\t\t\tev.P = p\n")
	}
	if m.IsBf {
		b.WriteString("\t\t\tev.B = p\n")
	}
	for _, a := range m.Args {
		b.WriteString(a.SerializeCall())
	}
	b.WriteString(`			p.mu.RLock()
			for _, h := range p.private` + m.Name + `s {
				h.Handle` + m.Name + `(ev)
			}
			p.mu.RUnlock()
		}
`)
	return b.String()
}

type GoConstant struct {
	Name    string
	Value   string
	Type    string
	Comment string
}

func (c *GoConstant) Serialize() string {
	var means string
	if c.Comment != "" {
		means = " means "
	}

	return `// ` + c.Name + means + c.Comment + `
const ` + c.Name + ` = ` + c.Value + `

`
}

type GoArg struct {
	Name    string
	Comment string
	Type    *GoType
}

func (a *GoArg) IsCanError() bool {
	return a.Type.IsCanError()
}

func (a *GoArg) IsOutput() bool {
	return a.Type.IsOutput()
}
func (a *GoArg) IsInput() bool {
	return a.Type.IsInput()
}

func (a *GoArg) Serialize() string {
	return "\t// " + a.Name + " is the " + a.Comment + "\n\t" + a.Name + " " + a.Type.Serialize() + "\n"
}
func (a *GoArg) SerializeCall() string {
	call := a.Type.SerializeCall()
	if call == "" {
		return ""
	}
	if a.IsCanError() {
		return "\t\t\tev." + a.Name + ", ev." + a.Name + "Error = event." + call + "\n"
	}
	return "\t\t\tev." + a.Name + " = event." + call + "\n"
}

type GoType struct {
	Type      string
	Interface string
}

func (a *GoType) SerializeNew() string {
	return "New" + a.Interface + "(p.Context())"
}
func (a *GoType) Serialize() string {
	switch a.Type {
	case "string":
		return "string"
	case "object":
		if a.Interface == "" {
			return "Proxy"
		}
		return "*" + a.Interface // pointer
	case "uint":
		return "uint32"
	case "int":
		return "int32"
	case "fd":
		return "uintptr"
	case "fixed":
		return "float32"
	case "array":
		return "[]int32"
	case "new_id":
		if a.Interface == "" {
			return "Proxy"
		}
		return "*" + a.Interface // pointer
	case "error":
		return "error"
	default:
		println("Unknown type: " + a.Type)
		return "struct {} // unknown type"
	}
}
func (a *GoType) SerializeCall() string {
	switch a.Type {
	case "string":
		return "String()"
	case "object":
		if a.Interface == "" {
			return "Proxy(p.Context())"
		}
		return "Proxy(p.Context()).(*" + a.Interface + ")"
	case "uint":
		return "Uint32()"
	case "int":
		return "Int32()"
	case "fd":
		return "FD()"
	case "fixed":
		return "Float32()"
	case "array":
		return "Array()"
	case "new_id":
		return "NewId(new(" + a.Interface + "), p.Context()).(*" + a.Interface + ")"
	case "error":
		return ""
	default:
		return ""
	}
}

func (a *GoType) IsCanError() bool {
	return a.Type == "fd"
}
func (a *GoType) isNewId() bool {
	return a.Type == "new_id" && a.Interface != ""
}
func (a *GoType) IsOutput() bool {
	return a.isNewId() && a.Type != "error"
}
func (a *GoType) IsInput() bool {
	return !a.isNewId() && a.Type != "error"
}

type GoHandler struct {
	Name string
}

func (h *GoHandler) Serialize(s *GoStruct) string {
	return `// ` + s.Name + h.Name + `Handler is the handler interface for ` + s.Name + h.Name + `Event
type ` + s.Name + h.Name + `Handler interface {
	Handle` + s.Name + h.Name + `(` + s.Name + h.Name + `Event)
}

// Add` + h.Name + `Handler removes the ` + h.Name + ` handler
func (p *` + s.Name + `) Add` + h.Name + `Handler(h ` + s.Name + h.Name + `Handler) {
	if h != nil {
		p.mu.Lock()
		p.private` + s.Name + h.Name + `s = append(p.private` + s.Name + h.Name + `s, h)
		p.mu.Unlock()
	}
}

// Remove` + h.Name + `Handler adds the ` + h.Name + ` handler
func (p *` + s.Name + `) Remove` + h.Name + `Handler(h ` + s.Name + h.Name + `Handler) {
	p.mu.Lock()
	defer p.mu.Unlock()

	for i, e := range p.private` + s.Name + h.Name + `s {
		if e == h {
			p.private` + s.Name + h.Name + `s = append(p.private` + s.Name + h.Name + `s[:i], p.private` + s.Name + h.Name + `s[i+1:]...)
			break
		}
	}
}
`
}
func (h *GoHandler) SerializeAttr(s *GoStruct) string {
	return "\t" + `private` + s.Name + h.Name + `s []` + s.Name + h.Name + `Handler` + "\n"
}
