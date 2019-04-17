package bootstrap

import (
	"fmt"

	"gitlab.com/kode4food/ale/api"
	"gitlab.com/kode4food/ale/bootstrap/builtin"
	"gitlab.com/kode4food/ale/compiler/arity"
	"gitlab.com/kode4food/ale/compiler/encoder"
	"gitlab.com/kode4food/ale/compiler/special"
	"gitlab.com/kode4food/ale/macro"
)

// Error messages
const (
	BuiltInNotFound = "built-in not found: %s"
	SpecialNotFound = "special form not found: %s"
	MacroNotFound   = "macro not found: %s"
)

const (
	defBuiltInName = "def-builtin"
	defSpecialName = "def-special"
	defMacroName   = "def-macro"

	orMore = -1
)

func (b *bootstrap) builtIns() {
	b.specialForms()
	b.initialFunctions()
	b.availableFunctions()
}

func (b *bootstrap) specialForms() {
	b.special("do", special.Do)
	b.special("if", special.If)
	b.special("let", special.Let)
	b.special("fn", special.Fn)
	b.special("eval", special.Eval)
	b.special("declare", special.Declare)
	b.special("def", special.Bind)
	b.special("quote", special.Quote)
	b.special("defmacro", special.DefMacro)
	b.special("macroexpand-1", special.MacroExpand1)
	b.special("macroexpand", special.MacroExpand)
}

func (b *bootstrap) initialFunctions() {
	manager := b.manager

	defBuiltIn := func(args ...api.Value) api.Value {
		arity.AssertFixed(1, len(args))
		ns := manager.GetRoot()
		n := args[0].(api.LocalSymbol).Name()
		if nf, ok := b.funcMap[n]; ok {
			ns.Bind(n, nf)
			return args[0]
		}
		panic(fmt.Errorf(BuiltInNotFound, n))
	}

	defSpecial := func(args ...api.Value) api.Value {
		arity.AssertFixed(1, len(args))
		ns := manager.GetRoot()
		n := args[0].(api.LocalSymbol).Name()
		if sf, ok := b.specialMap[n]; ok {
			ns.Bind(n, sf)
			return args[0]
		}
		panic(fmt.Errorf(SpecialNotFound, n))
	}

	defMacro := func(args ...api.Value) api.Value {
		ns := manager.GetRoot()
		n := args[0].(api.LocalSymbol).Name()
		if sf, ok := b.macroMap[n]; ok {
			ns.Bind(n, sf)
			return args[0]
		}
		panic(fmt.Errorf(MacroNotFound, n))
	}

	ns := b.manager.GetRoot()
	ns.Bind(defBuiltInName, api.NormalFunction(defBuiltIn))
	ns.Bind(defSpecialName, api.NormalFunction(defSpecial))
	ns.Bind(defMacroName, api.NormalFunction(defMacro))
}

func (b *bootstrap) availableFunctions() {
	b.applicative("read", builtin.Read, 1)
	b.applicative("is-eq", builtin.IsIdentical, 1, orMore)
	b.applicative("is-nil", builtin.IsNil, 1)
	b.applicative("is-atom", builtin.IsAtom, 1)
	b.applicative("is-keyword", builtin.IsKeyword, 1)

	b.macro("syntax-quote", macro.SyntaxQuote)
	b.applicative("is-macro", builtin.IsMacro, 1)
	b.applicative("sym", builtin.Sym, 1)
	b.applicative("gensym", builtin.GenSym, 0, 1)
	b.applicative("is-symbol", builtin.IsSymbol, 1)
	b.applicative("is-local", builtin.IsLocal, 1)
	b.applicative("is-qualified", builtin.IsQualified, 1)

	b.applicative("str", builtin.Str)
	b.applicative("str!", builtin.ReaderStr)
	b.applicative("is-str", builtin.IsStr, 1)

	b.applicative("seq", builtin.Seq, 1)
	b.applicative("first", builtin.First, 1)
	b.applicative("rest", builtin.Rest, 1)
	b.applicative("last", builtin.Last, 1)
	b.applicative("cons", builtin.Cons, 2)
	b.applicative("conj", builtin.Conj, 1, orMore)
	b.applicative("len", builtin.Len, 1)
	b.applicative("nth", builtin.Nth, 2)
	b.applicative("get", builtin.Get, 2)
	b.applicative("assoc", builtin.Assoc)
	b.applicative("list", builtin.List)
	b.applicative("vector", builtin.Vector)

	b.applicative("is-seq", builtin.IsSeq, 1)
	b.applicative("is-len", builtin.IsLen, 1)
	b.applicative("is-indexed", builtin.IsIndexed, 1)
	b.applicative("is-assoc", builtin.IsAssoc, 1)
	b.applicative("is-mapped", builtin.IsMapped, 1)
	b.applicative("is-list", builtin.IsList, 1)
	b.applicative("is-vector", builtin.IsVector, 1)

	b.applicative("+", builtin.Add)
	b.applicative("-", builtin.Sub, 1, orMore)
	b.applicative("*", builtin.Mul)
	b.applicative("/", builtin.Div, 1, orMore)
	b.applicative("mod", builtin.Mod, 1, orMore)

	b.applicative("=", builtin.Eq, 1, orMore)
	b.applicative("!=", builtin.Neq, 1, orMore)
	b.applicative(">", builtin.Gt, 1, orMore)
	b.applicative(">=", builtin.Gte, 1, orMore)
	b.applicative("<", builtin.Lt, 1, orMore)
	b.applicative("<=", builtin.Lte, 1, orMore)

	b.applicative("is-pos-inf", builtin.IsPosInf, 1)
	b.applicative("is-neg-inf", builtin.IsNegInf, 1)
	b.applicative("is-nan", builtin.IsNaN, 1)

	b.applicative("partial", builtin.Partial, 1, orMore)
	b.applicative("apply", builtin.Apply, 2, orMore)
	b.applicative("is-apply", builtin.IsApply, 1)
	b.applicative("is-special", builtin.IsSpecial, 1)

	b.applicative("go*", builtin.Go, 1)
	b.applicative("chan", builtin.Chan, 0)
	b.applicative("promise", builtin.Promise, 0, 1)
	b.applicative("is-promise", builtin.IsPromise, 1)

	b.applicative("lazy-seq*", builtin.LazySequence, 1)
	b.applicative("concat", builtin.Concat)
	b.applicative("filter", builtin.Filter, 2)
	b.applicative("map", builtin.Map, 2, orMore)
	b.applicative("take", builtin.Take, 2)
	b.applicative("drop", builtin.Drop, 2)
	b.applicative("reduce", builtin.Reduce, 2, 3)
	b.applicative("for-each*", builtin.ForEach, 2)

	b.applicative("raise", builtin.Raise, 1)
	b.applicative("recover", builtin.Recover, 2)
	b.applicative("defer", builtin.Defer, 2)

	b.applicative("current-time", builtin.CurrentTime, 0)
}

func (b *bootstrap) applicative(name api.Name, call api.Call, arity ...int) {
	fn := api.ApplicativeFunction(call)
	b.builtIn(name, fn, arity...)
}

func (b *bootstrap) macro(name api.Name, call macro.Call) {
	b.macroMap[name] = call
}

func (b *bootstrap) special(name api.Name, call encoder.Call) {
	b.specialMap[name] = call
}

func (b *bootstrap) builtIn(name api.Name, fn *api.Function, a ...int) {
	fn.ArityChecker = arity.MakeChecker(a...)
	b.funcMap[name] = fn
}