package main

var currentChunk *Chunk

func beginCompile() {
	currentChunk = &Chunk{
		code:      []byte{},
		lines:     []int{},
		constants: []any{},
	}
}

func emitByte(b byte) {
	currentChunk.code = append(currentChunk.code, b)
}

func emitOp(op Opcode) {
	emitByte(byte(op))
}

func emitConstant(value any) {
	emitOp(OpLoad)
	emitByte(byte(len(currentChunk.constants)))
	currentChunk.constants = append(currentChunk.constants, value)
}

func emitJump(op Opcode) int {
	emitOp(op)
	emitByte(0xff)
	emitByte(0xff)
	return len(currentChunk.code) - 2
}

func emitLoop(start int) {
	emitOp(OpLoop)
	offset := len(currentChunk.code) - start + 2
	if offset > 0xffff {
		panic("Loop body too large.")
	}
	emitByte(byte(offset >> 8))
	emitByte(byte(offset))
}

func patchJump(offset int) {
	jump := len(currentChunk.code) - offset - 2
	if jump > 0xffff {
		panic("Too much code to jump over.")
	}
	currentChunk.code[offset] = byte(jump >> 8)
	currentChunk.code[offset+1] = byte(jump)
}

func (i If) compile() {
	i.Cond.compile()
	falseJump := emitJump(OpJumpIfFalse)
	for _, n := range i.Then {
		n.compile()
	}
	endJump := emitJump(OpJump)
	patchJump(falseJump)
	for _, n := range i.Else {
		n.compile()
	}
	patchJump(endJump)
}

func (w While) compile() {
	loopStart := len(currentChunk.code)
	w.Cond.compile()
	exitJump := emitJump(OpJumpIfFalse)
	for _, n := range w.Body {
		n.compile()
	}
	emitLoop(loopStart)
	patchJump(exitJump)
}

func (f For) compile() {

}

func (f Func) compile() {

}

func (s Sub) compile() {

}

func (r Return) compile() {

}

func (b Break) compile() {

}

func (c Continue) compile() {

}

func (f FuncCall) compile() {

}

func (a Assignment) compile() {

}

func (a AssignIndex) compile() {

}

func (u Unary) compile() {

}

func (b Binary) compile() {

}

func (n Number) compile() {

}

func (s String) compile() {

}

func (b Bool) compile() {

}

func (l List) compile() {

}

func (v Variable) compile() {

}

func (i Index) compile() {

}
