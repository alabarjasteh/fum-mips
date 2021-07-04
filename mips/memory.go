package mips

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
)

type Memory map[int]int // map[PC]Data

func NewMemory() Memory {
	mem := Memory{}
	mem.loadMemoryFromFile()
	return mem
}

func (mem Memory) loadMemoryFromFile() {
	file, err := os.Open("./memory_state1.txt")
	if err != nil {
		log.Fatalf("failed to open")

	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	i := 0
	for scanner.Scan() {
		val, err := strconv.ParseInt(scanner.Text(), 2, 64)
		if err != nil {
			log.Fatal("cannot parse memory state from file")
		}
		mem[i] = int(val)
		fmt.Printf("mem: %b\n", mem[i])
		i += 4
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}
