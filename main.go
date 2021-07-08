package main

import (
	"flag"
	"log"
	"time"

	"github.com/alabarjasteh/mips-simulator/mips"
)

type Instruction int64

func main() {
	memFile := flag.String("file", "array-max-min.txt", "initiating memory state")
	flag.Parse()
	log.Printf("Load memory from: %v\n", *memFile)

	mem := mips.NewMemory(*memFile)
	cpu := mips.NewCPU(mem)

	ticker := time.NewTicker(time.Millisecond * 200)
	done := make(chan bool)

	fetchClockChan := make(chan string) // blocking channels, for synchronization of pipeline stages
	decodeClockChan := make(chan string)
	executeClockChan := make(chan string)
	memoryClockChan := make(chan string)
	writebackClockChan := make(chan string)

	ifDecChan := make(chan mips.IfDec, 1) // non-blocking channels with buffers size = 1. These act as inter-stage's registers.
	decExcChan := make(chan mips.DecExc, 1)
	exMemChan := make(chan mips.ExMem, 1)
	memWBChan := make(chan mips.MemWB, 1)

	go func() {
		for {
			<-ticker.C
			log.Println("\n\nTik...")

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
