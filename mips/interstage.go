package mips

type IfDec struct {
	NPC int
	IR  int
}

type DecExc struct {
	instruction *Instruction
}

type ExMem struct {
	instruction *Instruction
	aluOut      int32
}

type MemWB struct {
	instruction  *Instruction
	data         int32
	needsToWrite bool
}
