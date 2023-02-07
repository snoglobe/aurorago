package main

import "math"

type ChunkType int

const (
	TypeProgram ChunkType = iota
	TypeFunction
	TypeSubroutine
)

type Chunk struct {
	code      []byte
	lines     []int
	constants []any
}

type AuroraFunction struct {
	name  string
	args  []string
	arity int
	body  *Chunk
}

type CallFrame struct {
	locals    [256]any
	function  AuroraFunction
	pc        int
	dest      uint8
	chunkType ChunkType
}

// register-based virtual machine
type AuroraVM struct {
	registers [256]*any
	callStack []CallFrame
	globals   map[int]any
}

func NewAuroraVM(chunk *Chunk) *AuroraVM {
	return &AuroraVM{
		registers: [256]*any{},
		callStack: []CallFrame{
			{[256]any{}, AuroraFunction{"[script]", []string{}, 0, chunk}, 0, 0, TypeProgram},
		},
		globals: map[int]any{},
	}
}

type Opcode uint8

const (
	OpLoad           Opcode = iota // LOAD <constant> <register>
	OpStore                        // STORE <register> <local>
	OpStoreGlobal                  // STOREGLOBAL <register> <global>
	OpAdd                          // ADD <register (a)> <register (b)> <register (dest)>
	OpAddTo                        // ADDTO <register (a)> <register (b)>
	OpSub                          // SUB <register (a)> <register (b)> <register (dest)>
	OpSubFrom                      // SUBFROM <register (a)> <register (b)>
	OpMul                          // MUL <register (a)> <register (b)> <register (dest)>
	OpDiv                          // DIV <register (a)> <register (b)> <register (dest)>
	OpMod                          // MOD <register (a)> <register (b)> <register (dest)>
	OpNeg                          // NEG <register (a)> <register (dest)>
	OpNot                          // NOT <register (a)> <register (dest)>
	OpEqual                        // EQUAL <register (a)> <register (b)> <register (dest)>
	OpNotEqual                     // NOTEQUAL <register (a)> <register (b)> <register (dest)>
	OpLess                         // LESS <register (a)> <register (b)> <register (dest)>
	OpLessEqual                    // LESSEQUAL <register (a)> <register (b)> <register (dest)>
	OpGreater                      // GREATER <register (a)> <register (b)> <register (dest)>
	OpGreaterEqual                 // GREATEREQUAL <register (a)> <register (b)> <register (dest)>
	OpJump                         // JUMP <short offset>
	OpJumpIfFalse                  // JUMPIFFALSE <register (a)> <short offset>
	OpJumpIfTrue                   // JUMPIFTRUE <register (a)> <short offset>
	OpJumpIfEqual                  // JUMPIFEQUAL <register (a)> <register (b)> <short offset>
	OpJumpIfNotEqual               // JUMPIFNOTEQUAL <register (a)> <register (b)> <short offset>
	OpLoop                         // LOOP <short offset>
	OpCall                         // CALL <register (func)> <register (n args)> <register (base of args)> <register (dest)>
	OpReturn                       // RETURN <register (a)>
	OpIndex                        // INDEX <register (a)> <register (b)> <register (dest)>
	OpIndexAssign                  // INDEXASSIGN <register (a)> <register (b)> <register (c)>
)

func (vm *AuroraVM) readByte() byte {
	frame := &vm.callStack[len(vm.callStack)-1]
	frame.pc++
	return frame.function.body.code[frame.pc-1]
}

func (vm *AuroraVM) readConstant() any {
	frame := &vm.callStack[len(vm.callStack)-1]
	return frame.function.body.constants[vm.readByte()]
}

func (vm *AuroraVM) readShort() uint16 {
	return uint16(vm.readByte()) | uint16(vm.readByte())<<8
}

func (vm *AuroraVM) Step() {
	frame := &vm.callStack[len(vm.callStack)-1]
	instruction := Opcode(vm.readByte())
	switch instruction {
	case OpLoad:
		constant := vm.readConstant()
		register := vm.readByte()
		vm.registers[register] = &constant
	case OpStore:
		register := vm.readByte()
		local := vm.readByte()
		frame.locals[local] = register
	case OpStoreGlobal:
		register := vm.readByte()
		global := vm.readByte()
		vm.globals[int(global)] = *vm.registers[register]
	case OpAdd:
		a := vm.readByte()
		b := vm.readByte()
		dest := vm.readByte()
		switch (*vm.registers[a]).(type) {
		case float64:
			*vm.registers[dest] = (*vm.registers[a]).(float64) + (*vm.registers[b]).(float64)
		}
	case OpAddTo:
		a := vm.readByte()
		b := vm.readByte()
		switch (*vm.registers[a]).(type) {
		case float64:
			*vm.registers[a] = (*vm.registers[a]).(float64) + (*vm.registers[b]).(float64)
		}
	case OpSub:
		a := vm.readByte()
		b := vm.readByte()
		dest := vm.readByte()
		switch (*vm.registers[a]).(type) {
		case float64:
			*vm.registers[dest] = (*vm.registers[a]).(float64) - (*vm.registers[b]).(float64)
		}
	case OpSubFrom:
		a := vm.readByte()
		b := vm.readByte()
		switch (*vm.registers[a]).(type) {
		case float64:
			*vm.registers[a] = (*vm.registers[a]).(float64) - (*vm.registers[b]).(float64)
		}
	case OpMul:
		a := vm.readByte()
		b := vm.readByte()
		dest := vm.readByte()
		switch (*vm.registers[a]).(type) {
		case float64:
			*vm.registers[dest] = (*vm.registers[a]).(float64) * (*vm.registers[b]).(float64)
		}
	case OpDiv:
		a := vm.readByte()
		b := vm.readByte()
		dest := vm.readByte()
		switch (*vm.registers[a]).(type) {
		case float64:
			*vm.registers[dest] = (*vm.registers[a]).(float64) / (*vm.registers[b]).(float64)
		}
	case OpMod:
		a := vm.readByte()
		b := vm.readByte()
		dest := vm.readByte()
		switch (*vm.registers[a]).(type) {
		case float64:
			*vm.registers[dest] = math.Mod((*vm.registers[a]).(float64), (*vm.registers[b]).(float64))
		}
	case OpNeg:
		a := vm.readByte()
		dest := vm.readByte()
		switch (*vm.registers[a]).(type) {
		case float64:
			*vm.registers[dest] = -(*vm.registers[a]).(float64)
		}
	case OpNot:
		a := vm.readByte()
		dest := vm.readByte()
		switch (*vm.registers[a]).(type) {
		case bool:
			*vm.registers[dest] = !(*vm.registers[a]).(bool)
		}
	case OpEqual:
		a := vm.readByte()
		b := vm.readByte()
		dest := vm.readByte()
		*vm.registers[dest] = *vm.registers[a] == *vm.registers[b]
	case OpNotEqual:
		a := vm.readByte()
		b := vm.readByte()
		dest := vm.readByte()
		*vm.registers[dest] = *vm.registers[a] != *vm.registers[b]
	case OpLess:
		a := vm.readByte()
		b := vm.readByte()
		dest := vm.readByte()
		switch (*vm.registers[a]).(type) {
		case float64:
			*vm.registers[dest] = (*vm.registers[a]).(float64) < (*vm.registers[b]).(float64)
		}
	case OpLessEqual:
		a := vm.readByte()
		b := vm.readByte()
		dest := vm.readByte()
		switch (*vm.registers[a]).(type) {
		case float64:
			*vm.registers[dest] = (*vm.registers[a]).(float64) <= (*vm.registers[b]).(float64)
		}
	case OpGreater:
		a := vm.readByte()
		b := vm.readByte()
		dest := vm.readByte()
		switch (*vm.registers[a]).(type) {
		case float64:
			*vm.registers[dest] = (*vm.registers[a]).(float64) > (*vm.registers[b]).(float64)
		}
	case OpGreaterEqual:
		a := vm.readByte()
		b := vm.readByte()
		dest := vm.readByte()
		switch (*vm.registers[a]).(type) {
		case float64:
			*vm.registers[dest] = (*vm.registers[a]).(float64) >= (*vm.registers[b]).(float64)
		}
	case OpJump:
		offset := vm.readShort()
		frame.pc += int(offset)
	case OpJumpIfFalse:
		offset := vm.readShort()
		register := vm.readByte()
		if !(*vm.registers[register]).(bool) {
			frame.pc += int(offset)
		}
	case OpJumpIfTrue:
		offset := vm.readShort()
		register := vm.readByte()
		if (*vm.registers[register]).(bool) {
			frame.pc += int(offset)
		}
	case OpJumpIfEqual:
		offset := vm.readShort()
		a := vm.readByte()
		b := vm.readByte()
		if *vm.registers[a] == *vm.registers[b] {
			frame.pc += int(offset)
		}
	case OpJumpIfNotEqual:
		offset := vm.readShort()
		a := vm.readByte()
		b := vm.readByte()
		if *vm.registers[a] != *vm.registers[b] {
			frame.pc += int(offset)
		}
	case OpLoop:
		offset := vm.readShort()
		frame.pc -= int(offset)
	case OpCall: // f a b d
		function := vm.readByte()
		arity := vm.readByte()
		registerBase := vm.readByte()
		dest := vm.readByte()
		funcObj := (*vm.registers[function]).(AuroraFunction)
		if int(arity) != funcObj.arity {
			panic("Arity mismatch")
		}
		frame := CallFrame{
			function:  funcObj,
			pc:        0,
			locals:    [256]any{},
			dest:      dest,
			chunkType: TypeFunction,
		}
		for i := 0; i < int(arity); i++ {
			frame.locals[i] = vm.registers[registerBase+byte(i)]
		}
		vm.callStack = append(vm.callStack, frame)
	case OpReturn:
		frame := vm.callStack[len(vm.callStack)-1]
		vm.callStack = vm.callStack[:len(vm.callStack)-1]
		*vm.registers[frame.dest] = vm.registers[vm.readByte()]
	case OpIndex:
		a := vm.readByte()
		b := vm.readByte()
		dest := vm.readByte()
		switch (*vm.registers[a]).(type) {
		case []any:
			*vm.registers[dest] = (*vm.registers[a]).([]any)[(*vm.registers[b]).(int)]
		}
	case OpIndexAssign:
		a := vm.readByte()
		b := vm.readByte()
		c := vm.readByte()
		switch (*vm.registers[a]).(type) {
		case []any:
			(*vm.registers[a]).([]any)[(*vm.registers[b]).(int)] = *vm.registers[c]
		}
	}
}
