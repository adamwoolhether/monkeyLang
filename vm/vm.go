package vm

import (
	"fmt"

	"github.com/adamwoolhether/monkeyLang/code"
	"github.com/adamwoolhether/monkeyLang/compiler"
	"github.com/adamwoolhether/monkeyLang/object"
)

const StackSize = 2048

var (
	// True and False allow implementation of immutable, unique
	// values. Defined as global vars gives a performance increase
	// meaning we won't have to allocate and unwrap different vars
	// each for comparison, because true is always true and false
	// is always false.
	True  = &object.Boolean{Value: true}
	False = &object.Boolean{Value: false}
)

// VM defines our virtual machine. It holds constants and instructions
// generated by the compiler, and has a stack which will be pre-allocated
// to have `StackSize` number of elements, and a stack pointer, which
// will increment or decremented to grow/shrink the stack.
type VM struct {
	constants    []object.Object
	instructions code.Instructions

	stack []object.Object
	sp    int // Always points to the next value. Top of stack is stack[sp-1]
}

func New(bytecode *compiler.Bytecode) *VM {
	return &VM{
		constants:    bytecode.Constants,
		instructions: bytecode.Instructions,
		stack:        make([]object.Object, StackSize),
		sp:           0,
	}
}

// Run turns VM into a virtual machine. It contains the heartbeat,
// main loop, and fetch-decode-execute cycle.
func (vm *VM) Run() error {

	// Increment over the instruction pointer, fetching the current
	// instruction by accessing vm.instructions, turning the byte
	// into an Opcode.
	for ip := 0; ip < len(vm.instructions); ip++ {
		op := code.Opcode(vm.instructions[ip])

		switch op {
		case code.OpConstant:
			// decode
			constIndex := code.ReadUint16(vm.instructions[ip+1:]) // decode operands into bytecode.
			ip += 2

			// execute
			if err := vm.push(vm.constants[constIndex]); err != nil { // push the const onto the stack.
				return err
			}
		case code.OpAdd, code.OpSub, code.OpMul, code.OpDiv:
			if err := vm.executeBinaryOperation(op); err != nil {
				return err
			}
		case code.OpPop:
			vm.pop()
		case code.OpTrue:
			if err := vm.push(True); err != nil {
				return err
			}
		case code.OpFalse:
			if err := vm.push(False); err != nil {
				return err
			}
		case code.OpEqual, code.OpNotEqual, code.OpGreaterThan:
			if err := vm.executeComparison(op); err != nil {
				return err
			}
		}
	}

	return nil
}

// StackTop returns the element at the top of the stack.
func (vm *VM) StackTop() object.Object {
	if vm.sp == 0 {
		return nil
	}

	return vm.stack[vm.sp-1]
}

// LastPoppedStackElem allows a sanity check about
// what element should have been on stack immediately
// before being popped off.
func (vm *VM) LastPoppedStackElem() object.Object {
	return vm.stack[vm.sp]
}

// push checks the stack size and adds the object to the stack
// and increments the stack pointer.
func (vm *VM) push(o object.Object) error {
	if vm.sp >= StackSize {
		return fmt.Errorf("stack overflow")
	}

	vm.stack[vm.sp] = o
	vm.sp++

	return nil
}

// pop return the element located at the top of the stack and
// decrements vm.sp, allowing it to eventually be overwritten.
func (vm *VM) pop() object.Object {
	o := vm.stack[vm.sp-1]
	vm.sp--

	return o
}

func (vm *VM) executeBinaryOperation(op code.Opcode) error {
	right := vm.pop()
	left := vm.pop()

	leftType := left.Type()
	rightType := right.Type()

	if leftType == object.INTEGER_OBJ && rightType == object.INTEGER_OBJ {
		return vm.executeBinaryIntegerOperation(op, left, right)
	}

	return fmt.Errorf("unsupported types for binary operation: %s %s", leftType, rightType)
}

func (vm *VM) executeBinaryIntegerOperation(op code.Opcode, left, right object.Object) error {
	leftValue := left.(*object.Integer).Value
	rightValue := right.(*object.Integer).Value

	var result int64

	switch op {
	case code.OpAdd:
		result = leftValue + rightValue
	case code.OpSub:
		result = leftValue - rightValue
	case code.OpMul:
		result = leftValue * rightValue
	case code.OpDiv:
		result = leftValue / rightValue
	default:
		return fmt.Errorf("unknown integer operator: %d", op)
	}

	return vm.push(&object.Integer{Value: result})
}

// executeComparison determines whether two operands are integers, pops
// them off the stack, and turns them nto *object.Booleans before
// pushing the result back on to the stack.
func (vm *VM) executeComparison(op code.Opcode) error {
	right := vm.pop()
	left := vm.pop()

	// If both operands are integers
	if left.Type() == object.INTEGER_OBJ && right.Type() == object.INTEGER_OBJ {
		return vm.executeIntegerComparison(op, left, right)
	}

	switch op {
	case code.OpEqual:
		return vm.push(nativeBoolToBooleanObject(right == left))
	case code.OpNotEqual:
		return vm.push(nativeBoolToBooleanObject(right != left))
	default:
		return fmt.Errorf("unknown operator: %d (%s %s)", op, left.Type(), right.Type())
	}
}

// executeIntegerComparison unwraps values contained in left & right, compares them
// and turns them into a resulting True or False.
func (vm *VM) executeIntegerComparison(op code.Opcode, left, right object.Object) error {
	leftValue := left.(*object.Integer).Value
	rightValue := right.(*object.Integer).Value

	switch op {
	case code.OpEqual:
		return vm.push(nativeBoolToBooleanObject(rightValue == leftValue))
	case code.OpNotEqual:
		return vm.push(nativeBoolToBooleanObject(rightValue != leftValue))
	case code.OpGreaterThan:
		return vm.push(nativeBoolToBooleanObject(leftValue > rightValue))
	default:
		return fmt.Errorf("unknown operator: %d", op)
	}
}

// nativeBoolToBooleanObject turns the inputted boolean into a True or False.
func nativeBoolToBooleanObject(input bool) *object.Boolean {
	if input {
		return True
	}

	return False
}
