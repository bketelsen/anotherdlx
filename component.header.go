package main

import (
	"fmt"
	"os"
	"time"

	"github.com/mackerelio/go-osstat/cpu"
	"github.com/mackerelio/go-osstat/memory"

	"github.com/yuriizinets/kyoto"
)

type ComponentHeader struct {
	MemoryTotal  int
	MemoryUsed   int
	MemoryCached int
	MemoryFree   int
	CPUUser      float64
	CPUSystem    float64
	CPUIdle      float64
}

func (c *ComponentHeader) Init() {
}
func (c *ComponentHeader) Async(p kyoto.Page) error {
	memory, err := memory.Get()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		return err
	}
	fmt.Printf("memory total: %d bytes\n", memory.Total)
	fmt.Printf("memory used: %d bytes\n", memory.Used)
	fmt.Printf("memory cached: %d bytes\n", memory.Cached)
	fmt.Printf("memory free: %d bytes\n", memory.Free)
	c.MemoryTotal = int(memory.Total / 1024 / 1024)
	c.MemoryUsed = int(memory.Used / 1024 / 1024)
	c.MemoryCached = int(memory.Cached / 1024 / 1024)
	c.MemoryFree = int(memory.Free / 1024 / 1024)
	before, err := cpu.Get()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		return err
	}
	time.Sleep(time.Duration(1) * time.Second)
	after, err := cpu.Get()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		return err
	}
	total := float64(after.Total - before.Total)
	fmt.Printf("cpu user: %f %%\n", float64(after.User-before.User)/total*100)
	fmt.Printf("cpu system: %f %%\n", float64(after.System-before.System)/total*100)
	fmt.Printf("cpu idle: %f %%\n", float64(after.Idle-before.Idle)/total*100)
	c.CPUUser = float64(after.User-before.User) / total * 100
	c.CPUSystem = float64(after.System-before.System) / total * 100
	c.CPUIdle = float64(after.Idle-before.Idle) / total * 100

	return nil
}
