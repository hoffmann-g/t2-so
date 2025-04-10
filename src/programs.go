package main

var Program1 = []func(p *Processor){
	func(p *Processor) { p.Registers["$t0"] = 1 },
	func(p *Processor) { p.Registers["$t1"] = 2 },
	func(p *Processor) {
		p.Registers["$t2"] = p.Registers["$t0"].(int) + p.Registers["$t1"].(int)
	},
	func(p *Processor) { p.Registers["$t3"] = p.Registers["$t2"].(int) * 2 },
	func(p *Processor) { p.Registers["$t0"] = 1 },
	func(p *Processor) { p.Registers["$t1"] = 2 },
	func(p *Processor) {
		p.Registers["$t2"] = p.Registers["$t0"].(int) + p.Registers["$t1"].(int)
	},
	func(p *Processor) { p.Registers["$t3"] = p.Registers["$t2"].(int) * 2 },
	func(p *Processor) { p.Registers["$t0"] = 1 },
	func(p *Processor) { p.Registers["$t1"] = 2 },
	func(p *Processor) {
		p.Registers["$t2"] = p.Registers["$t0"].(int) + p.Registers["$t1"].(int)
	},
	func(p *Processor) { p.Registers["$t3"] = p.Registers["$t2"].(int) * 2 },
	func(p *Processor) { p.Registers["$t0"] = 1 },
	func(p *Processor) { p.Registers["$t1"] = 2 },
	func(p *Processor) {
		p.Registers["$t2"] = p.Registers["$t0"].(int) + p.Registers["$t1"].(int)
	},
	func(p *Processor) { p.Registers["$t3"] = p.Registers["$t2"].(int) * 2 },
	func(p *Processor) { p.Registers["$t0"] = 1 },
	func(p *Processor) { p.Registers["$t1"] = 2 },
	func(p *Processor) {
		p.Registers["$t2"] = p.Registers["$t0"].(int) + p.Registers["$t1"].(int)
	},
	func(p *Processor) { p.Registers["$t3"] = p.Registers["$t2"].(int) * 2 },
	func(p *Processor) { p.Registers["$t0"] = 1 },
	func(p *Processor) { p.Registers["$t1"] = 2 },
	func(p *Processor) {
		p.Registers["$t2"] = p.Registers["$t0"].(int) + p.Registers["$t1"].(int)
	},
	func(p *Processor) { p.Registers["$t3"] = p.Registers["$t2"].(int) * 2 },
}

var Program2 = []func(p *Processor){
	func(p *Processor) { p.Registers["$a0"] = 10 },
	func(p *Processor) { p.Registers["$a1"] = 20 },
	func(p *Processor) {
		p.Registers["$a2"] = p.Registers["$a0"].(int) - p.Registers["$a1"].(int)
	},
	func(p *Processor) { p.Registers["$a3"] = p.Registers["$a2"].(int) / 2 },
}
