package main

import (
	"fmt"
	"time"
)

var pc int
var registers map[int]any

func main() {
	for {
		exec(pc)
		pc++
	}
}

func exec(pc int) {
	fmt.Printf("%06x\n", pc)
	// fmt.Println(Data[pc])

	time.Sleep(1 * time.Second)
}
