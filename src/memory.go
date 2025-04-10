package main

const MemorySize = 1024

var Data [MemorySize]any

func loadKernelIntoMemory() {
	// load ISR
	Data[0] = kernel.TimeISR
	Data[1] = kernel.IOISR

	// load scalonator instructions
}
