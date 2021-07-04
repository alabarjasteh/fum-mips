package main

import (
	"fmt"

	"github.com/alabarjasteh/mips-simulator/mips"
)

type Instruction int64

func main() {
	mem := mips.NewMemory()
	cpu := mips.NewCPU(mem)

	for {
		_, instData := cpu.Fetch()
		instruction, err := cpu.Decode(instData)
		if err != nil {
			fmt.Printf("error: %v", err)
			return
		}
		err = cpu.Execute(instruction)
		if err != nil {
			fmt.Printf("error: %v", err)
			return
		}
	}
}

// func fetchBurst(pc *int, mem *memory.Mem, reg *register.IfDec) {
// 	inst, err := mem.FetchInstruction(*pc)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	reg.IR = inst
// 	log.Printf("IR : %v", reg.IR)
// 	*pc += 4
// 	reg.NPC = *pc
// }

// func decodeBurst(rf *register.File, ifDec *register.IfDec, decEx *register.DecEx) {
// 	r1, r2 := getReadingRegisters(rf, ifDec.IR)
// 	imm := extractImmediate(ifDec.IR)
// 	decEx.IR = ifDec.IR
// 	decEx.NPC = ifDec.NPC
// 	decEx.R1 = r1
// 	decEx.R2 = r2
// 	decEx.Imm = imm
// }

// func executeBurst() {

// }

// func memoryBurst()    {}
// func writebackBurst() {}
