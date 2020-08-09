package builtin

import (
	"github.com/kode4food/ale/compiler/encoder"
	"github.com/kode4food/ale/data"
	"github.com/kode4food/ale/internal/sequence"
)

// Apply performs a parameterized function call
func Apply(args ...data.Value) data.Value {
	fn := args[0].(data.Caller).Call()
	al := len(args)
	if al == 2 {
		return fn(sequence.ToValues(args[1].(data.Sequence))...)
	}
	last := al - 1
	ls := sequence.ToValues(args[last].(data.Sequence))
	prependedArgs := append(args[1:last], ls...)
	return fn(prependedArgs...)
}

// IsApply tests whether a value is callable
func IsApply(args ...data.Value) data.Value {
	_, ok := args[0].(data.Caller)
	return data.Bool(ok)
}

// IsSpecial tests whether not a function is a special form
func IsSpecial(args ...data.Value) data.Value {
	_, ok := args[0].(encoder.Call)
	return data.Bool(ok)
}
