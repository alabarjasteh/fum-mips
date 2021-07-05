package mips

import (
	"errors"
	"log"
	"time"
)

type CPU struct {
	PC      int
	Mem     Memory
	RegFile [32]int32
}

func NewCPU(mem Memory) *CPU {
	return &CPU{
		PC:      0,
		RegFile: [32]int32{},
		Mem:     mem,
	}
}

func (cpu *CPU) Fetch(clockChan chan string, outChan chan IfDec) (npc int, instData int) {
	for {
		<-clockChan

		ins, ok := cpu.Mem[cpu.PC]
		if !ok {
			log.Printf("Fetch -> PC: %v, no instruction to fetch", cpu.PC)
			continue
		}
		log.Printf("Fetch -> PC: %v, instruction: %32b", cpu.PC, ins)
		cpu.PC += 4
		npc = cpu.PC
		out := IfDec{NPC: npc, IR: ins}
		time.Sleep(time.Millisecond * 100) // wait till Decode brings its old data from chan first, then write new value on it.
		outChan <- out
	}
}

func (cpu *CPU) Decode(clockChan chan string, inChan chan IfDec, outChan chan DecExc) error {
	for {
		<-clockChan

		select {
		case in := <-inChan:
			IR := in.IR
			opcode := IR >> 26
			opcodeType, err := getOpcodeType(opcode)
			if err != nil {
				return err
			}

			ins := &Instruction{}
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
				log.Printf("Decode -> (R Type) instruction funcCode: %b", insTypeR.FuncCode)
				if !ok {
					return errors.New("unsupported funcCode")
				}
				insTypeR.AluFunc = f
				ins.TypeR = insTypeR

				outChan <- DecExc{instruction: ins}

			case OpcodeTypeI:
				insTypeI := &InstructionTypeI{}
				insTypeI.Immediate = IR & 0xFFFF
				insTypeI.TargetRegister = (IR >> 16) & 0b11111
				insTypeI.SourceRegister = (IR >> 21) & 0b11111
				f, ok := FunctionTypeIMap[opcode]
				log.Printf("Decode -> (I Type) instruction opcode: %b", opcode)
				if !ok {
					return errors.New("unsupported opcode")
				}
				insTypeI.AluFunc = f
				ins.TypeI = insTypeI

				outChan <- DecExc{instruction: ins}
			}

		default:
			log.Println("Decode -> NOP")
		}

	}
}

func (cpu *CPU) Execute(clockChan chan string, inChan chan DecExc, outChan chan ExMem) {
	for {
		<-clockChan

		select {
		case in := <-inChan:
			instruction := in.instruction
			var aluOut int32
			switch instruction.OpcodeType {
			case OpcodeTypeR:
				r := instruction.TypeR
				aluOut = r.AluFunc(cpu, r.SourceRegister, r.TargetRegister, r.DestinationRegister)
			case OpcodeTypeI:
				i := instruction.TypeI
				aluOut = i.AluFunc(cpu, i.SourceRegister, i.TargetRegister, i.Immediate)
			}
			time.Sleep(time.Millisecond * 100)
			outChan <- ExMem{instruction, aluOut}
			log.Printf("Execute -> aluOut: %v", aluOut)

		default:
			log.Println("Execute -> NOP")
		}
	}
}

func (cpu *CPU) Memory(clockChan chan string, inChan chan ExMem, outChan chan MemWB) {
	for {
		<-clockChan

		select {
		case in := <-inChan:
			instruction := in.instruction
			aluOut := in.aluOut
			log.Printf("Memory -> instruction opcode: %b", instruction.Opcode)

			if ok, v := instruction.IsLoadStore(); ok {
				if v == "LW" {
					memoryData := int32(cpu.Mem[int(aluOut)])
					time.Sleep(time.Millisecond * 100)
					outChan <- MemWB{instruction: instruction, data: memoryData, needsToWrite: true}
					continue
				}
				if v == "SW" {
					cpu.Mem[int(aluOut)] = int(cpu.RegFile[instruction.TypeI.TargetRegister])
					time.Sleep(time.Millisecond * 100)
					outChan <- MemWB{instruction: instruction, data: 0, needsToWrite: false}
					continue
				}
			}

			time.Sleep(time.Millisecond * 100)
			outChan <- MemWB{instruction: instruction, data: aluOut, needsToWrite: true}

		default:
			log.Println("Memory -> NOP")
		}
	}
}

func (cpu *CPU) Writeback(clockChan chan string, inChan chan MemWB) {
	for {
		<-clockChan

		select {
		case in := <-inChan:
			needsToWrite := in.needsToWrite
			instruction := in.instruction
			data := in.data

			if !needsToWrite {
				continue
			}

			switch instruction.OpcodeType {
			case OpcodeTypeR:
				cpu.RegFile[instruction.TypeR.DestinationRegister] = data
				log.Printf("Writeback -> (R Type) destination register (rd): %v, data: %v", instruction.TypeR.DestinationRegister, data)
			case OpcodeTypeI:
				cpu.RegFile[instruction.TypeI.TargetRegister] = data
				log.Printf("Writeback -> (I Type) target register (rt): %v, data: %v", instruction.TypeI.TargetRegister, data)
			}

		default:
			log.Println("Writeback -> NOP")
		}
	}
}
