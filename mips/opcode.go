package mips

import "errors"

type OpcodeType byte

var ErrInvalidOpcode = errors.New("invalid opcode")

const (
	OpcodeTypeInvalid OpcodeType = iota
	OpcodeTypeR
	OpcodeTypeI
)

func getOpcodeType(opcode int) (OpcodeType, error) {
	if opcode == 0 {
		return OpcodeTypeR, nil
	}
	if _, ok := FunctionTypeIMap[opcode]; ok {
		return OpcodeTypeI, nil
	}
	return OpcodeTypeInvalid, ErrInvalidOpcode
}
