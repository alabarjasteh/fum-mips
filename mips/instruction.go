package mips

var FunctionTypeRMap = map[int]FunctionTypeR{
	0b100000: Add,
	0b100010: Sub,
	0b100101: Or,
	0b100100: And,
	0b101010: Slt,
}

var FunctionTypeIMap = map[int]FunctionTypeI{
	0b101011: SwAddrCalc,
	0b100011: LwAddrCalc,
	0b001000: Addi,
	0b001010: Slti,
	0b001100: Andi,
	0b001101: Ori,
}

type Instruction struct {
	Opcode     int
	OpcodeType OpcodeType
	TypeR      *InstructionTypeR
	TypeI      *InstructionTypeI
}

type InstructionTypeR struct {
	SourceRegister      int
	TargetRegister      int
	DestinationRegister int
	FuncCode            int
	AluFunc             FunctionTypeR
}

type InstructionTypeI struct {
	SourceRegister int
	TargetRegister int
	Immediate      int
	AluFunc        FunctionTypeI
}

func (ins *Instruction) IsLoadStore() (bool, string) {
	if ins.OpcodeType == OpcodeTypeI {
		if ins.Opcode == 0b101011 {
			return true, "SW"
		}
		if ins.Opcode == 0b100011 {
			return true, "LW"
		}
	}
	return false, ""
}
