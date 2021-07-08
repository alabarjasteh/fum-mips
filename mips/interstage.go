package mips

type IfDec struct {
	NPC int
	IR  int
}

type DecExc struct {
	NPC         int
	Instruction *Instruction
	Rs          int32
	Rt          int32
	Rd          int32
}

type ExMem struct {
	NPC         int
	Instruction *Instruction
	Rt          int32
	AluOut      int32
	BranchAddr  int
}

type MemWB struct {
	Instruction *Instruction
	Data        int32
}
