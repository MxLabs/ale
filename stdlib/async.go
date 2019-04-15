package stdlib

import (
	"fmt"
	"runtime"
	"sync"
	"sync/atomic"

	"gitlab.com/kode4food/ale/api"
)

type (
	// Emitter is an interface that is used to emit values to a Channel
	Emitter interface {
		Writer
		Closer
		Error(interface{})
	}

	// Promise represents a Value that will eventually be resolved
	Promise interface {
		api.Value
		Deliver(api.Value) api.Value
		Resolve() api.Value
	}

	channelResult struct {
		value api.Value
		error interface{}
	}

	channelWrapper struct {
		seq    chan channelResult
		status uint32
	}

	channelEmitter struct {
		ch *channelWrapper
	}

	channelSequence struct {
		once Do
		ch   *channelWrapper

		isSeq  bool
		result channelResult
		rest   api.Sequence
	}

	promise struct {
		cond  *sync.Cond
		state uint32
		val   api.Value
	}
)

const (
	promiseUndelivered uint32 = iota
	promiseDelivered
)

// Error messages
const (
	ExpectedUndelivered = "can't deliver a promise twice"
)

const (
	channelReady uint32 = iota
	channelCloseRequested
	channelClosed
)

var emptyResult = channelResult{value: api.Nil, error: nil}

func (ch *channelWrapper) Close() {
	if status := atomic.LoadUint32(&ch.status); status != channelClosed {
		atomic.StoreUint32(&ch.status, channelClosed)
		close(ch.seq)
	}
}

// NewChannel produces a Emitter and Sequence pair
func NewChannel() (Emitter, api.Sequence) {
	seq := make(chan channelResult, 0)
	ch := &channelWrapper{
		seq:    seq,
		status: channelReady,
	}
	return NewChannelEmitter(ch), NewChannelSequence(ch)
}

// NewChannelEmitter produces an Emitter for sending values to a Go chan
func NewChannelEmitter(ch *channelWrapper) Emitter {
	r := &channelEmitter{
		ch: ch,
	}
	runtime.SetFinalizer(r, func(e *channelEmitter) {
		defer func() { recover() }()
		if s := atomic.LoadUint32(&ch.status); s != channelClosed {
			e.Close()
		}
	})
	return r
}

// Write will send a Value to the Go chan
func (e *channelEmitter) Write(v api.Value) {
	if s := atomic.LoadUint32(&e.ch.status); s == channelReady {
		e.ch.seq <- channelResult{v, nil}
	}
	if s := atomic.LoadUint32(&e.ch.status); s == channelCloseRequested {
		e.Close()
	}
}

// Error will send an Error to the Go chan
func (e *channelEmitter) Error(err interface{}) {
	if s := atomic.LoadUint32(&e.ch.status); s == channelReady {
		e.ch.seq <- channelResult{api.Nil, err}
	}
	e.Close()
}

// Close will Close the Go chan
func (e *channelEmitter) Close() {
	runtime.SetFinalizer(e, nil)
	e.ch.Close()
}

func (e *channelEmitter) Type() api.Name {
	return "channel-emitter"
}

func (e *channelEmitter) String() string {
	return api.DumpString(e)
}

// NewChannelSequence produces a new Sequence whose values come from a Go chan
func NewChannelSequence(ch *channelWrapper) api.Sequence {
	r := &channelSequence{
		once:   Once(),
		ch:     ch,
		result: emptyResult,
		rest:   api.EmptyList,
	}
	runtime.SetFinalizer(r, func(c *channelSequence) {
		defer func() { recover() }()
		if s := atomic.LoadUint32(&c.ch.status); s == channelReady {
			atomic.StoreUint32(&c.ch.status, channelCloseRequested)
			<-ch.seq // consume whatever is there
		}
	})
	return r
}

func (c *channelSequence) resolve() *channelSequence {
	c.once(func() {
		runtime.SetFinalizer(c, nil)
		ch := c.ch
		if result, isSeq := <-ch.seq; isSeq {
			c.isSeq = isSeq
			c.result = result
			c.rest = NewChannelSequence(ch)
		}
	})
	if e := c.result.error; e != nil {
		panic(e)
	}
	return c
}

func (c *channelSequence) IsSequence() bool {
	return c.resolve().isSeq
}

func (c *channelSequence) First() api.Value {
	return c.resolve().result.value
}

func (c *channelSequence) Rest() api.Sequence {
	return c.resolve().rest
}

func (c *channelSequence) Split() (api.Value, api.Sequence, bool) {
	r := c.resolve()
	return r.result.value, r.rest, r.isSeq
}

func (c *channelSequence) Prepend(v api.Value) api.Sequence {
	return &channelSequence{
		once:   Never(),
		isSeq:  true,
		result: channelResult{value: v, error: nil},
		rest:   c,
	}
}

func (c *channelSequence) Type() api.Name {
	return "channel-sequence"
}

func (c *channelSequence) String() string {
	return api.DumpString(c)
}

// NewPromise instantiates a new Promise
func NewPromise() Promise {
	return &promise{
		cond:  sync.NewCond(new(sync.Mutex)),
		state: promiseUndelivered,
	}
}

func (p *promise) Caller() api.Call {
	return func(args ...api.Value) api.Value {
		if len(args) > 0 {
			return p.Deliver(args[0])
		}
		return p.Resolve()
	}
}

func (p *promise) Resolve() api.Value {
	if atomic.LoadUint32(&p.state) == promiseDelivered {
		return p.val
	}

	cond := p.cond
	cond.L.Lock()
	defer cond.L.Unlock()
	for atomic.LoadUint32(&p.state) != promiseDelivered {
		cond.Wait()
	}
	return p.val
}

func (p *promise) checkNewValue(v api.Value) api.Value {
	if v == p.val {
		return p.val
	}
	panic(fmt.Errorf(ExpectedUndelivered))
}

func (p *promise) Deliver(v api.Value) api.Value {
	if atomic.LoadUint32(&p.state) == promiseDelivered {
		return p.checkNewValue(v)
	}

	cond := p.cond
	cond.L.Lock()
	defer cond.L.Unlock()

	if p.state == promiseUndelivered {
		p.val = v
		atomic.StoreUint32(&p.state, promiseDelivered)
		cond.Broadcast()
		return v
	}

	cond.Wait()
	return p.checkNewValue(v)
}

func (p *promise) Type() api.Name {
	return "promise"
}

func (p *promise) String() string {
	return api.DumpString(p)
}
