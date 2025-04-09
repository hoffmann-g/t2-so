package modules

import (
	"fmt"
	"t1-so/src/processor"
)

type Scalonator struct {
}

func (s *Scalonator) Scalonate(processor *processor.Processor, pm *ProcessManager) {
	fmt.Println("Scalonating processes 😈")
	processor.Pc += 100
	// increment pc to simulate process execution
	// access cpu registers and pc
}
