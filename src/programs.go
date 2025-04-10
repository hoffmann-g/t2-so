package main

var Program1 = []func(){
	func() { processor.Registers["$t0"] = 1 },
	func() { processor.Registers["$t1"] = 2 },
	func() {
		processor.Registers["$t2"] = processor.Registers["$t0"].(int) + processor.Registers["$t1"].(int)
	},
	func() { processor.Registers["$t3"] = processor.Registers["$t2"].(int) * 2 },
	func() { processor.Registers["$t0"] = 1 },
	func() { processor.Registers["$t1"] = 2 },
	func() {
		processor.Registers["$t2"] = processor.Registers["$t0"].(int) + processor.Registers["$t1"].(int)
	},
	func() { processor.Registers["$t3"] = processor.Registers["$t2"].(int) * 2 },
	func() { processor.Registers["$t0"] = 1 },
	func() { processor.Registers["$t1"] = 2 },
	func() {
		processor.Registers["$t2"] = processor.Registers["$t0"].(int) + processor.Registers["$t1"].(int)
	},
	func() { processor.Registers["$t3"] = processor.Registers["$t2"].(int) * 2 },
	func() { processor.Registers["$t0"] = 1 },
	func() { processor.Registers["$t1"] = 2 },
	func() {
		processor.Registers["$t2"] = processor.Registers["$t0"].(int) + processor.Registers["$t1"].(int)
	},
	func() { processor.Registers["$t3"] = processor.Registers["$t2"].(int) * 2 },
	func() { processor.Registers["$t0"] = 1 },
	func() { processor.Registers["$t1"] = 2 },
	func() {
		processor.Registers["$t2"] = processor.Registers["$t0"].(int) + processor.Registers["$t1"].(int)
	},
	func() { processor.Registers["$t3"] = processor.Registers["$t2"].(int) * 2 },
	func() { processor.Registers["$t0"] = 1 },
	func() { processor.Registers["$t1"] = 2 },
	func() {
		processor.Registers["$t2"] = processor.Registers["$t0"].(int) + processor.Registers["$t1"].(int)
	},
	func() { processor.Registers["$t3"] = processor.Registers["$t2"].(int) * 2 },
}

var Program2 = []func(){
	func() { processor.Registers["$a0"] = 10 },
	func() { processor.Registers["$a1"] = 20 },
	func() {
		processor.Registers["$a2"] = processor.Registers["$a0"].(int) - processor.Registers["$a1"].(int)
	},
	func() { processor.Registers["$a3"] = processor.Registers["$a2"].(int) / 2 },
}
