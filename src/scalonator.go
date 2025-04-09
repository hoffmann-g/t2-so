package main

type Scalonator struct {
}

func (s *Scalonator) Init() {}

func (s *Scalonator) Scalonate(processor *Processor, pm *ProcessManager) {
	// fmt.Println("Scalonating processes 😈")
	processor.Pc += 100
	// increment pc to simulate process execution
	// access cpu registers and pc
}
