package main

type Instruction struct {
	Label string
	Code  any
}

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

// var Program1 = []any{
// 	1,
// 	2,
// 	3,
// 	4,
// 	"data string",
// 	"A",
// 	"B",
// 	"C",
// 	func() { cpu.Registers["$t0"] = cpu.Registers["$t0"].(int) + 1 },
// 	func() { cpu.Registers["$t1"] = cpu.Registers["$t1"].(int) + 1 },
// 	func() { cpu.Registers["$t2"] = cpu.Registers["$t2"].(int) - 1 },
// 	func() { cpu.Registers["$t3"] = cpu.Registers["$t3"].(int) - 1 },
// 	func() { cpu.Registers["$t0"] = cpu.Registers["$t0"].(int) + 1 },
// 	func() { cpu.Registers["$t1"] = cpu.Registers["$t1"].(int) + 1 },
// 	func() { cpu.Registers["$t2"] = cpu.Registers["$t2"].(int) - 1 },
// 	func() { cpu.Registers["$t3"] = cpu.Registers["$t3"].(int) - 1 },
// 	func() { cpu.Registers["$t0"] = cpu.Registers["$t0"].(int) + 1 },
// 	func() { cpu.Registers["$t1"] = cpu.Registers["$t1"].(int) + 1 },
// 	func() { cpu.Registers["$t2"] = cpu.Registers["$t2"].(int) - 1 },
// 	func() { cpu.Registers["$t3"] = cpu.Registers["$t3"].(int) - 1 },
// 	func() { cpu.Registers["$t0"] = cpu.Registers["$t0"].(int) + 1 },
// 	func() { cpu.Registers["$t1"] = cpu.Registers["$t1"].(int) + 1 },
// 	func() { cpu.Registers["$t2"] = cpu.Registers["$t2"].(int) - 1 },
// 	func() { cpu.Registers["$t3"] = cpu.Registers["$t3"].(int) - 1 },
// 	func() { cpu.Registers["$t0"] = cpu.Registers["$t0"].(int) + 1 },
// 	func() { cpu.Registers["$t1"] = cpu.Registers["$t1"].(int) + 1 },
// 	func() { cpu.Registers["$t2"] = cpu.Registers["$t2"].(int) - 1 },
// 	func() { cpu.Registers["$t3"] = cpu.Registers["$t3"].(int) - 1 },
// 	func() { cpu.Registers["$t0"] = cpu.Registers["$t0"].(int) + 1 },
// 	func() { cpu.Registers["$t1"] = cpu.Registers["$t1"].(int) + 1 },
// 	func() { cpu.Registers["$t2"] = cpu.Registers["$t2"].(int) - 1 },
// 	func() { cpu.Registers["$t3"] = cpu.Registers["$t3"].(int) - 1 },
// 	func() { cpu.Registers["$t0"] = cpu.Registers["$t0"].(int) + 1 },
// 	func() { cpu.Registers["$t1"] = cpu.Registers["$t1"].(int) + 1 },
// 	func() { cpu.Registers["$t2"] = cpu.Registers["$t2"].(int) - 1 },
// 	func() { cpu.Registers["$t3"] = cpu.Registers["$t3"].(int) - 1 },
// 	func() { cpu.Registers["$t0"] = cpu.Registers["$t0"].(int) + 1 },
// 	func() { cpu.Registers["$t1"] = cpu.Registers["$t1"].(int) + 1 },
// 	func() { cpu.Registers["$t2"] = cpu.Registers["$t2"].(int) - 1 },
// 	func() { cpu.Registers["$t3"] = cpu.Registers["$t3"].(int) - 1 },
// }

var Program2 = []Instruction{
	{"SET $t0", func() { cpu.Registers["$t0"] = 1 }},
	{"SET $t1", func() { cpu.Registers["$t1"] = 1 }},
	{"SET $t2", func() { cpu.Registers["$t2"] = 1 }},
	{"SET $t3", func() { cpu.Registers["$t3"] = 1 }},
	{"MUL $t0", func() { cpu.Registers["$t0"] = cpu.Registers["$t0"].(int) * 2 }},
	{"MUL $t1", func() { cpu.Registers["$t1"] = cpu.Registers["$t1"].(int) * 3 }},
	{"MUL $t2", func() { cpu.Registers["$t2"] = cpu.Registers["$t2"].(int) * 4 }},
	{"MUL $t3", func() { cpu.Registers["$t3"] = cpu.Registers["$t3"].(int) * 5 }},
}

var Program3 = []Instruction{
	{"SET $t0", func() { cpu.Registers["$t0"] = "A" }},
	{"SET $t1", func() { cpu.Registers["$t1"] = "A" }},
	{"SET $t2", func() { cpu.Registers["$t2"] = "A" }},
	{"SET $t3", func() { cpu.Registers["$t3"] = "A" }},
	{"SET $t0", func() { cpu.Registers["$t0"] = "A" }},
	{"SET $t1", func() { cpu.Registers["$t1"] = "A" }},
	{"SET $t2", func() { cpu.Registers["$t2"] = "A" }},
	{"SET $t3", func() { cpu.Registers["$t3"] = "A" }},
}

var Program4 = []Instruction{
	{"SET $t0", func() { cpu.Registers["$t0"] = 13 }},
	{"SET $t1", func() { cpu.Registers["$t1"] = 13 }},
	{"SET $t2", func() { cpu.Registers["$t2"] = 13 }},
	{"SET $t3", func() { cpu.Registers["$t3"] = 13 }},
	{"SET $t0", func() { cpu.Registers["$t0"] = 13 }},
	{"SET $t1", func() { cpu.Registers["$t1"] = 13 }},
	{"SET $t2", func() { cpu.Registers["$t2"] = 13 }},
	{"SET $t3", func() { cpu.Registers["$t3"] = 13 }},
}

var ProgramIO = []Instruction{
	{"IO REQUEST", IORequestFunc("Digite seu nome:")},
	{"CONCAT", func() { cpu.Registers["$t0"] = "Resposta: " + cpu.Registers["$t0"].(string) }},
	{"SET $t1", func() { cpu.Registers["$t1"] = 42 }},
}

func IORequestFunc(message string) func() {
	return func() {
		cpu.RequestIO(message)
	}
}
