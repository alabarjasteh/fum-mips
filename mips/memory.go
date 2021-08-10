package mips

import (
	"bufio"
	"log"
	"os"
	"strconv"
)

type Memory map[int]int // map[PC]Data

func NewMemory(fileName string) Memory {
	mem := Memory{}
	mem.loadMemoryFromFile(fileName)
	return mem
}

func (mem Memory) loadMemoryFromFile(fileName string) {
	file, err := os.Open("./memory-initial-state/" + fileName)
	if err != nil {
		log.Fatalf("failed to open")

	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	i := 0
	for scanner.Scan() {
		text := scanner.Text()
		if len(text) == 0 {
			i += 4
			continue
		}
		val, err := strconv.ParseInt(text, 2, 64)
		if err != nil {
			log.Fatal("cannot parse memory state from file")
		}
		mem[i] = int(val)
		i += 4
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}
