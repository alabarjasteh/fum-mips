package mips

import (
	"errors"
	"fmt"
)

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
	0b000100: BeqBranchCond,
	0b000101: BneBranchCond,
}

type Instruction struct {
	IR         int
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

var InstructionName = map[int]string{
	0b100000: "ADD",
	0b100010: "SUB",
	0b100101: "OR",
	0b100100: "AND",
	0b101010: "SLT",
	0b101011: "SW",
	0b100011: "LW",
	0b001000: "ADDI",
	0b001010: "SLTI",
	0b001100: "ANDI",
	0b001101: "ORI",
	0b000100: "BEQ",
	0b000101: "BNE",
	0b000010: "J",
}

func ConstructInstruction(IR int) (*Instruction, error) {
	opcode := IR >> 26
	opcodeType, err := getOpcodeType(opcode)
	if err != nil {
		return &Instruction{}, err
	}
	ins := &Instruction{}
	ins.IR = IR
	ins.Opcode = opcode
	ins.OpcodeType = opcodeType

	switch opcodeType {
	case OpcodeTypeR:
		insTypeR := &InstructionTypeR{}
		insTypeR.FuncCode = IR & 0b111111
		insTypeR.DestinationRegister = (IR >> 11) & 0b11111
		insTypeR.TargetRegister = (IR >> 16) & 0b11111
		insTypeR.SourceRegister = (IR >> 21) & 0b11111
		f, ok := FunctionTypeRMap[insTypeR.FuncCode]
		if !ok {
			return &Instruction{}, errors.New("unsupported funcCode")
		}
		insTypeR.AluFunc = f
		ins.TypeR = insTypeR

	case OpcodeTypeI:
		insTypeI := &InstructionTypeI{}
		insTypeI.Immediate = IR & 0xFFFF
		insTypeI.TargetRegister = (IR >> 16) & 0b11111
		insTypeI.SourceRegister = (IR >> 21) & 0b11111
		f, ok := FunctionTypeIMap[opcode]
		if !ok {
			return &Instruction{}, errors.New("unsupported opcode")
		}
		insTypeI.AluFunc = f
		ins.TypeI = insTypeI
	}
	return ins, nil
}

func assemblyEquivalent(IR int) string {
	opCode := IR >> 26
	opType, _ := getOpcodeType(opCode)
	rs := (IR >> 21) & 0b11111
	rt := (IR >> 16) & 0b11111

	switch opType {
	case OpcodeTypeR:
		name := InstructionName[IR&0b111111]
		rd := (IR >> 11) & 0b11111
		return fmt.Sprintf("%v R$%v R$%v R$%v", name, rd, rs, rt)

	case OpcodeTypeI:
		name := InstructionName[opCode]
		imm := IR & 0xFFFF
		return fmt.Sprintf("%v R$%v R$%v (%v)", name, rt, rs, imm)

	case OpcodeTypeJ:
		imm := IR & 0x3FFFFFF
		return fmt.Sprintf("J %v", imm)
	}

	return "cannot parse"
}
