package data

import (
	"bytes"
	"errors"
	"fmt"
	"sort"
)

type (
	// Object maps a set of Values, known as keys, to another set of Values
	Object interface {
		object() // marker
		Sequence
		Mapped
		Counted
	}

	object struct {
		pair     Pair
		children [32]*object
	}
)

// Standard Keys
const (
	NameKey     = Keyword("name")
	TypeKey     = Keyword("type")
	CountKey    = Keyword("count")
	InstanceKey = Keyword("instance")
)

// Error messages
const (
	ErrMapNotPaired = "map does not contain an even number of elements"
)

// EmptyObject represents an empty Object
var EmptyObject = &object{}

// NewObject instantiates a new Object instance. Based on Phil Bagwell's
// Hashed Array Mapped Trie data structure. More information can be
// found at http://lampwww.epfl.ch/papers/idealhashtrees.pdf
func NewObject(pairs ...Pair) Object {
	res := &object{}
	for _, p := range pairs {
		res = res.Put(p).(*object)
	}
	return res
}

// ValuesToObject interprets a set of Values as an Object
func ValuesToObject(v ...Value) (Object, error) {
	if len(v)%2 != 0 {
		return nil, errors.New(ErrMapNotPaired)
	}
	res := &object{}
	for i := len(v) - 2; i >= 0; i -= 2 {
		res = res.Put(NewCons(v[i], v[i+1])).(*object)
	}
	return res, nil
}

func (*object) object() {}

func (o *object) Get(k Value) (Value, bool) {
	if o.pair == nil {
		return Nil, false
	}
	h := HashCode(k)
	return o.get(k, h)
}

func (o *object) get(k Value, hash uint64) (Value, bool) {
	if o.pair.Car().Equal(k) {
		return o.pair.Cdr(), true
	}
	bucket := o.children[hash&0x1f]
	if bucket != nil {
		return bucket.get(k, hash>>5)
	}
	return Nil, false
}

func (o *object) Put(p Pair) Sequence {
	if o.pair == nil {
		return &object{
			pair: p,
		}
	}
	h := HashCode(p.Car())
	return o.put(p, h)
}

func (o *object) put(p Pair, hash uint64) *object {
	if o.pair.Car().Equal(p.Car()) {
		return &object{
			pair:     p,
			children: o.children,
		}
	}

	idx := hash & 0x1f
	bucket := o.children[idx]
	if bucket == nil {
		bucket = &object{pair: p}
	} else {
		bucket = bucket.put(p, hash>>5)
	}

	// return a copy with the new bucket
	res := *o
	res.children[idx] = bucket
	return &res
}

func (o *object) Remove(k Value) (Value, Sequence, bool) {
	if o.pair == nil {
		return Nil, o, false
	}
	h := HashCode(k)
	if v, r, ok := o.remove(k, h); ok {
		if r != nil {
			return v, r, true
		}
		return v, EmptyObject, true
	}
	return Nil, o, false
}

func (o *object) remove(k Value, hash uint64) (Value, *object, bool) {
	if o.pair.Car().Equal(k) {
		return o.pair.Cdr(), o.promote(), true
	}
	idx := hash & 0x1f
	if bucket := o.children[idx]; bucket != nil {
		if v, r, ok := bucket.remove(k, hash>>5); ok {
			res := *o
			res.children[idx] = r
			return v, &res, true
		}
	}
	return nil, nil, false
}

func (o *object) promote() *object {
	for i, c := range o.children {
		if c != nil {
			res := *o
			res.pair = c.pair
			res.children[i] = c.promote()
			return &res
		}
	}
	return nil
}

func (o *object) First() Value {
	f, _, _ := o.Split()
	return f
}

func (o *object) Rest() Sequence {
	_, r, _ := o.Split()
	return r
}

func (o *object) Split() (Value, Sequence, bool) {
	if f := o.pair; f != nil {
		if r := o.promote(); r != nil {
			return f, r, true
		}
		return f, EmptyObject, true
	}
	return Nil, EmptyObject, false
}

func (o *object) Count() int {
	res := 0
	for _, r, ok := o.Split(); ok; _, r, ok = r.Split() {
		res++
	}
	return res
}

func (o *object) IsEmpty() bool {
	return o.pair == nil
}

// Call turns Object into a Function
func (o *object) Call(args ...Value) Value {
	res, ok := o.Get(args[0])
	if !ok && len(args) > 1 {
		return args[1]
	}
	return res
}

// Convention returns the Function's calling convention
func (o *object) Convention() Convention {
	return ApplicativeCall
}

// CheckArity performs a compile-time arity check for the Function
func (o *object) CheckArity(argCount int) error {
	return checkRangedArity(1, 2, argCount)
}

func (o *object) Equal(v Value) bool {
	if v, ok := v.(*object); ok {
		lp := sortedPairs(o.pairs())
		rp := sortedPairs(v.pairs())
		if len(lp) != len(rp) {
			return false
		}
		for i, l := range lp {
			if !l.Equal(rp[i]) {
				return false
			}
		}
		return true
	}
	return false
}

func (o *object) HashCode() uint64 {
	var h uint64
	for f, r, ok := o.Split(); ok; f, r, ok = r.Split() {
		p := f.(Pair)
		h *= HashCode(p.Car()) * HashCode(p.Cdr())
	}
	return h
}

func (o *object) pairs() Pairs {
	res := Pairs{}
	for f, r, ok := o.Split(); ok; f, r, ok = r.Split() {
		res = append(res, f.(Pair))
	}
	return res
}

func (o *object) String() string {
	var buf bytes.Buffer
	buf.WriteString("{")
	for i, p := range sortedPairs(o.pairs()) {
		if i > 0 {
			buf.WriteString(" ")
		}
		buf.WriteString(MaybeQuoteString(p.Car()))
		buf.WriteString(" ")
		buf.WriteString(MaybeQuoteString(p.Cdr()))
	}
	buf.WriteString("}")
	return buf.String()
}

func sortedPairs(p Pairs) Pairs {
	sort.Slice(p, func(l, r int) bool {
		ls := fmt.Sprintf("%s", p[l].Car().String())
		rs := fmt.Sprintf("%s", p[r].Car().String())
		return ls < rs
	})
	return p
}
