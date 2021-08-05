package mips

import (
	"log"
	"time"
)

type CPU struct { // no mutex is needed, all accesses to share data are sequential.
	PC                 int
	Mem                Memory
	RegFile            [32]int32
	ScoreBoard         [32]bool // shows if a register is pending
	BranchesInPipeline []int    // list of all BEQ and BNE instructions in the pipeline
	Stall              bool
}

func NewCPU(mem Memory) *CPU {
	return &CPU{
		Mem:        mem,
		RegFile:    [32]int32{},
		ScoreBoard: [32]bool{},
	}
}

func (cpu *CPU) Fetch(clockChan chan string, outChan chan IfDec) {
	for {
		<-clockChan

		time.Sleep(time.Millisecond * 50)

		IR, ok := cpu.Mem[cpu.PC]
		assembly := assemblyEquivalent(IR)
		instName := InstructionName[IR>>26]

		if cpu.Stall {
			log.Println("Fetch	-> Stall")
			continue
		}

		if len(cpu.BranchesInPipeline) != 0 {
			log.Println("Fetch	-> Killed <Branch in the pipeline>")
			continue
		}

		if !ok {
			log.Printf("Fetch	-> PC: %v, no instruction to fetch", cpu.PC)
			continue
		}
		log.Printf("Fetch	-> PC: %v, %v", cpu.PC, assembly)

		// if instruction is 'unconditional jump', change PC to target
		if instName == "J" {
			immediate := IR & 0x3FFFFFF // 26-bit immediate
			target := ((cpu.PC + 4) & 0xF0000000) | (immediate << 2)
			cpu.PC = target
			log.Printf("[*Fetch*] -> jump to PC: %v", cpu.PC)
			continue
		} else {
			cpu.PC = cpu.PC + 4
		}

		out := IfDec{NPC: cpu.PC, IR: IR}
		time.Sleep(time.Millisecond * 100) // wait till Decode brings its old data from chan first, then write new value on it.
		outChan <- out
	}
}

func (cpu *CPU) Decode(clockChan chan string, inChan chan IfDec, outChan chan DecExc) error {
	for {
		<-clockChan

		select {
		case in := <-inChan:
			NPC, IR := in.NPC, in.IR
			assembly := assemblyEquivalent(IR)
			opcode := IR >> 26
			opcodeType, err := getOpcodeType(opcode)
			if err != nil {
				return err
			}

			ins, err := ConstructInstruction(IR)
			if err != nil {
				return err
			}
			var decExc DecExc

			log.Printf("Decode 	-> %v", assembly)

			switch opcodeType {
			case OpcodeTypeR:
				rs := cpu.RegFile[ins.TypeR.SourceRegister]
				rt := cpu.RegFile[ins.TypeR.TargetRegister]
				rd := cpu.RegFile[ins.TypeR.DestinationRegister]
				decExc = DecExc{NPC: NPC, Instruction: ins, Rs: rs, Rt: rt, Rd: rd}

				if cpu.readingPendingRegs(ins) {
					inChan <- IfDec{NPC: NPC, IR: IR} // redo the Decode
					cpu.Stall = true
					continue
				}
				cpu.updateScoreBoard(ins, "decode")

			case OpcodeTypeI:
				rs := cpu.RegFile[ins.TypeI.SourceRegister]
				rt := cpu.RegFile[ins.TypeI.TargetRegister]
				decExc = DecExc{NPC: NPC, Instruction: ins, Rs: rs, Rt: rt}

				if cpu.readingPendingRegs(ins) {
					inChan <- IfDec{NPC: NPC, IR: IR} // redo the decode
					cpu.Stall = true
					continue
				}
				cpu.updateScoreBoard(ins, "decode")

				instName := InstructionName[ins.Opcode]
				if instName == "BEQ" || instName == "BNE" {
					cpu.BranchesInPipeline = append(cpu.BranchesInPipeline, in.IR)
				}
			}

			cpu.Stall = false
			time.Sleep(time.Millisecond * 100)
			outChan <- decExc

		default:
			log.Println("Decode 	-> NOP")
		}

	}
}

func (cpu *CPU) Execute(clockChan chan string, inChan chan DecExc, outChan chan ExMem) {
	for {
		<-clockChan

		select {
		case in := <-inChan:
			NPC, instruction, rs, rt, rd := in.NPC, in.Instruction, in.Rs, in.Rt, in.Rd
			assembly := assemblyEquivalent(instruction.IR)
			var aluOut int32
			var branchAddr int

			switch instruction.OpcodeType {
			case OpcodeTypeR:
				r := instruction.TypeR
				aluOut = r.AluFunc(cpu, rs, rt, rd)
			case OpcodeTypeI:
				i := instruction.TypeI
				aluOut = i.AluFunc(cpu, rs, rt, i.Immediate)
				branchAddr = NPC + (instruction.TypeI.Immediate << 2)
			}

			time.Sleep(time.Millisecond * 100)
			outChan <- ExMem{NPC: NPC, Instruction: instruction, Rt: rt, AluOut: aluOut, BranchAddr: branchAddr}
			log.Printf("Execute	-> %v, aluOut: %v", assembly, aluOut)

		default:
			log.Println("Execute	-> NOP")
		}
	}
}

func (cpu *CPU) Memory(clockChan chan string, inChan chan ExMem, outChan chan MemWB) {
	for {
		<-clockChan

		select {
		case in := <-inChan:
			instruction, aluOut, branchAddr := in.Instruction, in.AluOut, in.BranchAddr
			assembly := assemblyEquivalent(instruction.IR)
			instName := InstructionName[instruction.Opcode]
			memWB := MemWB{Instruction: instruction}

			switch instName {
			case "LW":
				memoryData := int32(cpu.Mem[int(aluOut)])
				memWB.Data = memoryData
				log.Printf("Memory	-> %v, load %v from memory loction(%v)", assembly, memoryData, aluOut)

			case "SW":
				cpu.Mem[int(aluOut)] = int(in.Rt)
				log.Printf("Memory	-> %v, writes %v into memory location (%v)", assembly, in.Rt, aluOut)

			case "BEQ", "BNE":
				if aluOut == 1 { // branch is taken
					time.Sleep(time.Millisecond * 100)
					cpu.PC = branchAddr
					log.Printf("Memory	-> %v, branch is taken", assembly)
				} else {
					log.Printf("Memory	-> %v, branch is NOT taken", assembly)
				}
				// Pop from list
				cpu.BranchesInPipeline = cpu.BranchesInPipeline[:len(cpu.BranchesInPipeline)-1]
				memWB.Data = aluOut

			default:
				memWB.Data = aluOut
				log.Printf("Memory	-> %v", assembly)
			}

			time.Sleep(time.Millisecond * 100)
			outChan <- memWB

		default:
			log.Println("Memory	-> NOP")
		}
	}
}

func (cpu *CPU) Writeback(clockChan chan string, inChan chan MemWB) {
	for {
		<-clockChan

		select {
		case in := <-inChan:
			instruction, data := in.Instruction, in.Data
			assembly := assemblyEquivalent(instruction.IR)
			instName := InstructionName[instruction.Opcode]

			switch instruction.OpcodeType {
			case OpcodeTypeR:
				cpu.RegFile[instruction.TypeR.DestinationRegister] = data
				cpu.updateScoreBoard(instruction, "writeback")
				log.Printf("Writeback	-> %v,  writes %v into R$%v", assembly, data, instruction.TypeR.DestinationRegister)

			case OpcodeTypeI:
				if instName == "SW" || instName == "BEQ" || instName == "BNE" {
					log.Printf("Writeback	-> %v ,no write to register file", assembly)
					continue
				}
				cpu.RegFile[instruction.TypeI.TargetRegister] = data
				cpu.updateScoreBoard(instruction, "writeback")
				log.Printf("Writeback	-> %v,  writes %v into R$%v", assembly, data, instruction.TypeI.TargetRegister)

			}

		default:
			log.Println("Writeback	-> NOP")
		}
	}
}

func (cpu *CPU) readingPendingRegs(ins *Instruction) bool {
	switch ins.OpcodeType {
	case OpcodeTypeR:
		if cpu.ScoreBoard[ins.TypeR.SourceRegister] || cpu.ScoreBoard[ins.TypeR.TargetRegister] {
			return true
		}

	case OpcodeTypeI:
		instName := InstructionName[ins.Opcode]
		if cpu.ScoreBoard[ins.TypeI.SourceRegister] ||
			((instName == "SW" || instName == "BEQ" || instName == "BNE") && cpu.ScoreBoard[ins.TypeI.TargetRegister]) {
			return true
		}
	}
	return false
}

func (cpu *CPU) updateScoreBoard(ins *Instruction, stage string) {
	switch ins.OpcodeType {
	case OpcodeTypeR:
		if stage == "decode" {
			cpu.ScoreBoard[ins.TypeR.DestinationRegister] = true
		} else if stage == "writeback" {
			cpu.ScoreBoard[ins.TypeR.DestinationRegister] = false
		}
	case OpcodeTypeI:
		insName := InstructionName[ins.Opcode]
		if insName != "SW" && insName != "BEQ" && insName != "BNE" {
			if stage == "decode" {
				cpu.ScoreBoard[ins.TypeI.TargetRegister] = true
			} else if stage == "writeback" {
				cpu.ScoreBoard[ins.TypeI.TargetRegister] = false
			}
		}
	}
}
