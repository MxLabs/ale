package encoder

import (
	"gitlab.com/kode4food/ale/compiler/internal/analysis"
	"gitlab.com/kode4food/ale/data"
	"gitlab.com/kode4food/ale/namespace"
	"gitlab.com/kode4food/ale/runtime/isa"
)

type (
	// Type exposes an interface for stateful compiler encoding
	Type interface {
		data.Value

		Name() data.Name
		Parent() Type
		Child() Type
		NamedChild(data.Name) Type

		Emit(isa.Opcode, ...isa.Coder)
		Code() []isa.Word
		StackSize() int

		NewLabel() *Label

		Globals() namespace.Type
		Constants() data.Values
		AddConstant(data.Value) isa.Index

		Closure() data.Names
		ResolveClosure(data.LocalSymbol) (isa.Index, bool)

		PushArgs(data.Names, bool)
		PopArgs()
		ResolveArg(data.LocalSymbol) (isa.Index, bool, bool)

		LocalCount() int
		PushLocals()
		PopLocals()
		AddLocal(data.Name) isa.Index
		ResolveLocal(data.LocalSymbol) (isa.Index, bool)

		ResolveScope(data.LocalSymbol) (Scope, bool)
		InScope(data.LocalSymbol) bool
	}

	encoder struct {
		parent    Type
		globals   namespace.Type
		constants data.Values
		closure   data.Names
		name      data.Name
		args      argsStack
		locals    []Locals
		code      isa.Instructions
		maxLocal  int
		nextLocal int
		nextLabel int
	}
)

func newEncoder(globals namespace.Type) *encoder {
	return &encoder{
		globals:   globals,
		constants: data.Values{},
		closure:   data.Names{},
		args:      argsStack{},
		locals:    []Locals{{}},
		code:      isa.Instructions{},
	}
}

// NewEncoder instantiates a new Encoder
func NewEncoder(globals namespace.Type) Type {
	return newEncoder(globals)
}

func (e *encoder) child() *encoder {
	return &encoder{
		parent:    e,
		constants: data.Values{},
		closure:   data.Names{},
		args:      argsStack{},
		locals:    []Locals{{}},
		code:      isa.Instructions{},
	}
}

// Child creates a child Type
func (e *encoder) Child() Type {
	return e.child()
}

func (e *encoder) NamedChild(name data.Name) Type {
	res := e.child()
	res.name = name
	return res
}

// Name returns the name of this encoder (example: a function's name)
func (e *encoder) Name() data.Name {
	return e.name
}

// Parent returns the parent of this encoder
func (e *encoder) Parent() Type {
	return e.parent
}

// Emit adds instructions to the Type's eventual output
func (e *encoder) Emit(oc isa.Opcode, args ...isa.Coder) {
	words := make([]isa.Word, len(args))
	for i, a := range args {
		words[i] = a.Word()
	}
	e.code = append(e.code, isa.New(oc, words...))
}

// Word returns the encoder's resulting VM instructions
func (e *encoder) Code() []isa.Word {
	return isa.Flatten(e.code)
}

// StackSize returns the encoder's calculated stack size
func (e *encoder) StackSize() int {
	res, _ := analysis.CalculateStackSize(e.code)
	return res
}

// Globals returns the global name/value map
func (e *encoder) Globals() namespace.Type {
	if e.globals != nil {
		return e.globals
	}
	if e.parent != nil {
		return e.parent.Globals()
	}
	return nil
}

func (e *encoder) String() string {
	return "encoder"
}
