package mips

import (
	"errors"
	"fmt"
)

type CPU struct {
	PC      int
	Memory  Memory
	RegFile [32]int32
}

func NewCPU(mem Memory) *CPU {
	return &CPU{
		PC:      0,
		RegFile: [32]int32{},
		Memory:  mem,
	}
}

func (cpu *CPU) Fetch() (npc int, instData int) {
	ins := cpu.Memory[cpu.PC]
	fmt.Printf("instruction: %b\n", ins)
	cpu.PC += 4
	npc = cpu.PC
	return npc, ins
}

func (cpu *CPU) Decode(insData int) (*Instruction, error) {
	fmt.Printf("inst: %b\n", insData)
	opcode := insData >> 26
	fmt.Printf("opcode: %b\n", opcode)
	opcodeType, err := getOpcodeType(opcode)
	if err != nil {
		return nil, err
	}

	ins := &Instruction{}
	ins.Opcode = opcode
	ins.OpcodeType = opcodeType

	switch opcodeType {
	case OpcodeTypeR:
		insTypeR := &InstructionTypeR{}
		insTypeR.FuncCode = insData & 0b111111
		insTypeR.DestinationRegister = (insData >> 11) & 0b11111
		insTypeR.TargetRegister = (insData >> 16) & 0b11111
		insTypeR.SourceRegister = (insData >> 21) & 0b11111
		f, ok := FunctionTypeRMap[insTypeR.FuncCode]
		if !ok {
			return nil, errors.New("unsupported funcCode")
		}
		insTypeR.Function = f

		ins.TypeR = insTypeR
	case OpcodeTypeI:
		insTypeI := &InstructionTypeI{}
		insTypeI.Immediate = insData & 0xFFFF
		insTypeI.TargetRegister = (insData >> 16) & 0b11111
		insTypeI.SourceRegister = (insData >> 21) & 0b11111
		f, ok := FunctionTypeIMap[opcode]
		if !ok {
			return nil, errors.New("unsupported opcode")
		}
		insTypeI.Function = f

		ins.TypeI = insTypeI
	}

	return ins, nil
}

func (cpu *CPU) Execute(ins *Instruction) error {
	var err error
	switch ins.OpcodeType {
	case OpcodeTypeR:
		r := ins.TypeR
		err = r.Function(cpu, r.SourceRegister, r.TargetRegister, r.DestinationRegister)
	case OpcodeTypeI:
		i := ins.TypeI
		err = i.Function(cpu, i.SourceRegister, i.TargetRegister, i.Immediate)
	}
	return err
}
