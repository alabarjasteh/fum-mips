package mips

// R Type Instructions
type FunctionTypeR func(cpu *CPU, rs int32, rt int32, rd int32) (aluOut int32)

func Add(cpu *CPU, rs int32, rt int32, rd int32) (aluOut int32) {
	aluOut = rs + rt
	return aluOut
}

func And(cpu *CPU, rs int32, rt int32, rd int32) (aluOut int32) {
	aluOut = rs & rt
	return aluOut
}

func Slt(cpu *CPU, rs int32, rt int32, rd int32) (aluOut int32) {
	if rs < rt {
		aluOut = 1
	} else {
		aluOut = 0
	}
	return aluOut
}

func Sub(cpu *CPU, rs int32, rt int32, rd int32) (aluOut int32) {
	aluOut = rs - rt
	return aluOut
}

func Or(cpu *CPU, rs int32, rt int32, rd int32) (aluOut int32) {
	aluOut = rs | rt
	return aluOut
}

// I Type Instruction
type FunctionTypeI func(cpu *CPU, rs int32, rt int32, imm int) (aluOut int32)

func Addi(cpu *CPU, rs int32, rt int32, imm int) (aluOut int32) {
	aluOut = rs + int32(int16(imm))
	return aluOut
}

func Andi(cpu *CPU, rs int32, rt int32, imm int) (aluOut int32) {
	aluOut = rs & int32(int16(imm))
	return aluOut
}

func LwAddrCalc(cpu *CPU, rs int32, rt int32, imm int) (aluOut int32) {
	aluOut = rs + int32(int16(imm))
	return aluOut
}

func Ori(cpu *CPU, rs int32, rt int32, imm int) (aluOut int32) {
	aluOut = rs | int32(int16(imm))
	return aluOut
}

func Slti(cpu *CPU, rs int32, rt int32, imm int) (aluOut int32) {
	if rs < int32(int16(imm)) {
		aluOut = 1
	} else {
		aluOut = 0
	}
	return aluOut
}

func SwAddrCalc(cpu *CPU, rs int32, rt int32, imm int) (aluOut int32) {
	aluOut = rs + int32(int16(imm))
	return aluOut
}

func BeqBranchCond(cpu *CPU, rs int32, rt int32, imm int) (aluOut int32) {
	if rs == rt {
		return 1
	}
	return 0
}

func BneBranchCond(cpu *CPU, rs int32, rt int32, imm int) (aluOut int32) {
	if rs == rt {
		return 1
	}
	return 0
}
