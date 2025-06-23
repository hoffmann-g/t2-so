package main

type Instruction struct {
	Label string
	Code  any
}

// alphabet
var Program1 = []Instruction{
	{"A", "A"},
	{"B", "B"},
	{"C", "C"},
	{"D", "D"},
	{"E", "E"},
	{"F", "F"},
	{"G", "G"},
	{"H", "H"},
	{"I", "I"},
	{"J", "J"},
	{"K", "K"},
	{"L", "L"},
	{"M", "M"},
	{"N", "N"},
	{"O", "O"},
	{"P", "P"},
	{"Q", "Q"},
	{"R", "R"},
	{"S", "S"},
	{"T", "T"},
	{"U", "U"},
	{"V", "V"},
	{"W", "W"},
	{"X", "X"},
	{"Y", "Y"},
	{"Z", "Z"},
}

// calc
var Program2 = []Instruction{
	{"SET $t0 1", func() { cpu.Registers["$t0"] = 1 }},
	{"SET $t1 2", func() { cpu.Registers["$t1"] = 2 }},
	{"SET $t2 3", func() { cpu.Registers["$t2"] = 3 }},
	{"SET $t3 4", func() { cpu.Registers["$t3"] = 4 }},
	{"MUL $t0 2", func() { cpu.Registers["$t0"] = cpu.Registers["$t0"].(int) * 2 }},
	{"MUL $t1 2", func() { cpu.Registers["$t1"] = cpu.Registers["$t1"].(int) * 2 }},
	{"MUL $t2 2", func() { cpu.Registers["$t2"] = cpu.Registers["$t2"].(int) * 2 }},
	{"MUL $t3 2", func() { cpu.Registers["$t3"] = cpu.Registers["$t3"].(int) * 2 }},
}

// set letters
var Program3 = []Instruction{
	{"SET $t0 A", func() { cpu.Registers["$t0"] = "A" }},
	{"SET $t1 B", func() { cpu.Registers["$t1"] = "B" }},
	{"SET $t2 C", func() { cpu.Registers["$t2"] = "C" }},
	{"SET $t3 D", func() { cpu.Registers["$t3"] = "D" }},
	{"SET $t0 E", func() { cpu.Registers["$t0"] = "E" }},
	{"SET $t1 F", func() { cpu.Registers["$t1"] = "F" }},
	{"SET $t2 G", func() { cpu.Registers["$t2"] = "G" }},
	{"SET $t3 H", func() { cpu.Registers["$t3"] = "H" }},
}

// set numbers
var Program4 = []Instruction{
	{"SET $t0 0", func() { cpu.Registers["$t0"] = 0 }},
	{"SET $t1 1", func() { cpu.Registers["$t1"] = 1 }},
	{"SET $t2 2", func() { cpu.Registers["$t2"] = 2 }},
	{"SET $t3 3", func() { cpu.Registers["$t3"] = 3 }},
	{"SET $t0 4", func() { cpu.Registers["$t0"] = 4 }},
	{"SET $t1 5", func() { cpu.Registers["$t1"] = 5 }},
	{"SET $t2 6", func() { cpu.Registers["$t2"] = 6 }},
	{"SET $t3 7", func() { cpu.Registers["$t3"] = 7 }},
}

var ProgramIO = []Instruction{
	{"SET $t0 0", func() { cpu.Registers["$t0"] = 0 }},
	{"SET $t1 1", func() { cpu.Registers["$t1"] = 1 }},
	{"SET $t2 2", func() { cpu.Registers["$t2"] = 2 }},
	{"SET $t3 3", func() { cpu.Registers["$t3"] = 3 }},
	{"SYSCALL $t0", IORequestFunc("Digite seu nome:")},
	{"SET $t3 4", func() { cpu.Registers["$t1"] = 4 }},
	{"SET $t3 5", func() { cpu.Registers["$t2"] = 5 }},
	{"SET $t3 6", func() { cpu.Registers["$t3"] = 6 }},
	{"CONCAT", func() { cpu.Registers["$t0"] = cpu.Registers["$t0"].(string) }},
	{"SET $t1 OK", func() { cpu.Registers["$t1"] = "OK" }},
}

func IORequestFunc(message string) func() {
	return func() {
		cpu.RequestIO(message)
	}
}
