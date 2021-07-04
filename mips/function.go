package mips

import (
	"errors"
	"fmt"
)

// R Type Instructions
type FunctionTypeR func(cpu *CPU, rs int, rt int, rd int) error

func Add(cpu *CPU, rs int, rt int, rd int) error {

	cpu.RegFile[rd] = cpu.RegFile[rs] + cpu.RegFile[rt]

	return nil
}

func And(cpu *CPU, rs int, rt int, rd int) error {
	fmt.Printf("(\"and\" not implemented)\n")
	return errors.New("not implemented: and")
}

func Slt(cpu *CPU, rs int, rt int, rd int) error {
	fmt.Printf("(\"slt\" not implemented)\n")
	return errors.New("not implemented: slt")
}

func Sub(cpu *CPU, rs int, rt int, rd int) error {
	fmt.Printf("(\"sub\" not implemented)\n")
	return errors.New("not implemented: sub")
}

func Or(cpu *CPU, rs int, rt int, rd int) error {

	cpu.RegFile[rd] = cpu.RegFile[rs] | cpu.RegFile[rt]

	return nil
}

// I Type Instruction
type FunctionTypeI func(cpu *CPU, rs int, rt int, imm int) error

func Addi(cpu *CPU, rs int, rt int, imm int) error {

	cpu.RegFile[rt] = cpu.RegFile[rs] + int32(int16(imm))

	return nil
}

func Andi(cpu *CPU, rs int, rt int, imm int) error {
	fmt.Printf("(\"andi\" not implemented)\n")
	return errors.New("not implemented: andi")
}

func Lw(cpu *CPU, rs int, rt int, imm int) error {

	value := int32(cpu.Memory[int(int16(imm))+int(cpu.RegFile[rs])])
	cpu.RegFile[rt] = value

	return nil
}

func Ori(cpu *CPU, rs int, rt int, imm int) error {
	fmt.Printf("(\"ori\" not implemented)\n")
	return errors.New("not implemented: ori")
}

func Slti(cpu *CPU, rs int, rt int, imm int) error {

	if cpu.RegFile[rs] < int32(int16(imm)) {
		cpu.RegFile[rt] = 1
	} else {
		cpu.RegFile[rt] = 0
	}

	return nil
}

func Sw(cpu *CPU, rs int, rt int, imm int) error {

	cpu.Memory[int(int16(imm))+int(cpu.RegFile[rs])] = int(cpu.RegFile[rt])

	return nil
}
