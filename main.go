package main

import (
	"log"
	"time"

	"github.com/alabarjasteh/mips-simulator/mips"
)

type Instruction int64

func main() {
	mem := mips.NewMemory()
	cpu := mips.NewCPU(mem)

	ticker := time.NewTicker(time.Second)
	done := make(chan bool)

	fetchClockChan := make(chan string)
	decodeClockChan := make(chan string)
	executeClockChan := make(chan string)
	memoryClockChan := make(chan string)
	writebackClockChan := make(chan string)

	ifDecChan := make(chan mips.IfDec, 1) // write to chan is not blocking.
	decExcChan := make(chan mips.DecExc, 1)
	exMemChan := make(chan mips.ExMem, 1)
	memWBChan := make(chan mips.MemWB, 1)

	go func() {
		for {
			<-ticker.C
			log.Print("\n\nTik...")

			writebackClockChan <- "tik"
			memoryClockChan <- "tik"
			executeClockChan <- "tik"
			decodeClockChan <- "tik"
			fetchClockChan <- "tik"
		}
	}()

	// run stages in parallel
	go cpu.Fetch(fetchClockChan, ifDecChan)
	go cpu.Decode(decodeClockChan, ifDecChan, decExcChan)
	go cpu.Execute(executeClockChan, decExcChan, exMemChan)
	go cpu.Memory(memoryClockChan, exMemChan, memWBChan)
	go cpu.Writeback(writebackClockChan, memWBChan)

	<-done
}
