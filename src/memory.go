package main

const MemorySize = 1024

var Data [MemorySize]any

func loadKernelIntoMemory() {
	Data[0] = k.TimeIsr
	Data[1] = k.IoIsr
}
