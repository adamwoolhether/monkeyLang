package vm

import (
	"github.com/adamwoolhether/monkeyLang/code"
	"github.com/adamwoolhether/monkeyLang/object"
)

// Frame defines our stack frame that
// holds execution-relevant information.
type Frame struct {
	fn          *object.CompiledFunction // The compiled func referenced by the frame.
	ip          int                      // The IP for this frame/function.
	basePointer int                      // Points to the bottom of the stack of current call frame.
}

func NewFrame(fn *object.CompiledFunction, basePointer int) *Frame {
	f := &Frame{
		fn:          fn,
		ip:          -1,
		basePointer: basePointer,
	}

	return f
}

func (f *Frame) Instructions() code.Instructions {
	return f.fn.Instructions
}
