package memory

import (
	"t1-so/src/so"
)

const memorySize = 16000

var Data [memorySize]any

func InitMemory(kernel *so.Kernel) {
	Data[0] = kernel.TimeIsr
	Data[1] = kernel.IoIsr
}
