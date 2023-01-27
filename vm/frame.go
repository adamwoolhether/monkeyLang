package vm

import (
	"github.com/adamwoolhether/monkeyLang/code"
	"github.com/adamwoolhether/monkeyLang/object"
)

// Frame defines our stack frame that
// holds execution-relevant information.
type Frame struct {
	cl          *object.Closure // The compiled func referenced by the frame.
	ip          int             // The IP for this frame/function.
	basePointer int             // Points to the bottom of the stack of current call frame.
}

func NewFrame(cl *object.Closure, basePointer int) *Frame {
	f := &Frame{
		cl:          cl,
		ip:          -1,
		basePointer: basePointer,
	}

	return f
}

func (f *Frame) Instructions() code.Instructions {
	return f.cl.Fn.Instructions
}
