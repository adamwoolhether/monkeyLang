package vm

import (
	"github.com/adamwoolhether/monkeyLang/code"
	"github.com/adamwoolhether/monkeyLang/object"
)

// Frame defines our stack frame that
// holds execution-relevant information.
type Frame struct {
	fn *object.CompiledFunction // The compiled func referenced by the frame.
	ip int                      // The IP for this frame/function.
}

func NewFrame(fn *object.CompiledFunction) *Frame {
	return &Frame{fn: fn, ip: -1}
}

func (f *Frame) Instructions() code.Instructions {
	return f.fn.Instructions
}
