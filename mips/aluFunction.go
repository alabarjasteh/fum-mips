package mips

// R Type Instructions
type FunctionTypeR func(cpu *CPU, rs int, rt int, rd int) (aluOut int32)

func Add(cpu *CPU, rs int, rt int, rd int) (aluOut int32) {
	aluOut = cpu.RegFile[rs] + cpu.RegFile[rt]
	return aluOut
}

func And(cpu *CPU, rs int, rt int, rd int) (aluOut int32) {
	aluOut = cpu.RegFile[rs] & cpu.RegFile[rt]
	return aluOut
}

func Slt(cpu *CPU, rs int, rt int, rd int) (aluOut int32) {
	if cpu.RegFile[rs] < cpu.RegFile[rt] {
		aluOut = 1
	} else {
		aluOut = 0
	}
	return aluOut
}

func Sub(cpu *CPU, rs int, rt int, rd int) (aluOut int32) {
	aluOut = cpu.RegFile[rs] - cpu.RegFile[rt]
	return aluOut
}

func Or(cpu *CPU, rs int, rt int, rd int) (aluOut int32) {
	aluOut = cpu.RegFile[rs] | cpu.RegFile[rt]
	return aluOut
}

// I Type Instruction
type FunctionTypeI func(cpu *CPU, rs int, rt int, imm int) (aluOut int32)

func Addi(cpu *CPU, rs int, rt int, imm int) (aluOut int32) {
	aluOut = cpu.RegFile[rs] + int32(int16(imm))
	return aluOut
}

func Andi(cpu *CPU, rs int, rt int, imm int) (aluOut int32) {
	aluOut = cpu.RegFile[rs] & int32(int16(imm))
	return aluOut
}

func LwAddrCalc(cpu *CPU, rs int, rt int, imm int) (aluOut int32) {
	aluOut = cpu.RegFile[rs] + int32(int16(imm))
	return aluOut
}

func Ori(cpu *CPU, rs int, rt int, imm int) (aluOut int32) {
	aluOut = cpu.RegFile[rs] | int32(int16(imm))
	return aluOut
}

func Slti(cpu *CPU, rs int, rt int, imm int) (aluOut int32) {
	if cpu.RegFile[rs] < int32(int16(imm)) {
		aluOut = 1
	} else {
		aluOut = 0
	}
	return aluOut
}

func SwAddrCalc(cpu *CPU, rs int, rt int, imm int) (aluOut int32) {
	aluOut = cpu.RegFile[rs] + int32(int16(imm))
	return aluOut
}
