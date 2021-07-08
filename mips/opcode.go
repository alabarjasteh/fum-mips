package mips

import "errors"

type OpcodeType byte

var ErrInvalidOpcode = errors.New("invalid opcode")

const (
	OpcodeTypeInvalid OpcodeType = iota
	OpcodeTypeR
	OpcodeTypeI
	OpcodeTypeJ
)

func getOpcodeType(opcode int) (OpcodeType, error) {
	if opcode == 0 {
		return OpcodeTypeR, nil
	}
	if _, ok := FunctionTypeIMap[opcode]; ok {
		return OpcodeTypeI, nil
	}
	if opcode == 0b000010 {
		return OpcodeTypeJ, nil
	}
	return OpcodeTypeInvalid, ErrInvalidOpcode
}

// type OpcodeFunctionalType byte

// const (
// 	FunctionalTypeInvalid OpcodeFunctionalType = iota
// 	ALU
// 	ALUI
// 	LoadStore
// 	Branch
// 	Jump
// )

// func getFuncionalType(opcode int) (OpcodeFunctionalType, error) {
// 	if opcode == 0 {
// 		return ALU, nil
// 	}
// 	insName, ok := InstructionName[opcode]
// 	if !ok {
// 		return FunctionalTypeInvalid, ErrInvalidOpcode
// 	}
// 	if insName == "LW" || insName == "SW" {
// 		return LoadStore, nil
// 	}
// 	if insName == "BEQ" && insName == "BNE" {
// 		return Branch, nil
// 	}
// 	if insName == "J" {
// 		return Jump, nil
// 	}
// 	if _, ok := FunctionTypeIMap[opcode]; ok {
// 		return ALUI, nil
// 	}
// 	return FunctionalTypeInvalid, ErrInvalidOpcode
// }
