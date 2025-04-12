package main

import (
	"bufio"
	"fmt"
	"os"
	"reflect"
	"runtime"
	"strconv"
	"strings"
)

func HandleShell() {
	fmt.Println("Welcome to the shell!")
	fmt.Println("Type 'help' for a list of commands.")
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print("> ")

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		tokens := strings.Fields(line)

		if len(tokens) == 0 {
			fmt.Print("> ")
			continue
		}

		command := tokens[0]
		args := tokens[1:]

		switch command {
		case "help":
			printHelp()

		case "create":
			handleCreate(args)

		case "kill":
			handleKill(args)

		case "ps":
			handlePs(args)

		case "dump":
			handleDump(args)

		case "exec":
			handleExec(args)

		case "log":
			handleLog(args)

		case "exit":
			return

		default:
			fmt.Println("Unknown command:", command)
		}

		fmt.Print("> ")
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading shell input:", err)
	}
}

func printHelp() {
	fmt.Println("Available commands:")
	fmt.Println("  help - Show this help message")
	fmt.Println("  create <program-name> - Create a new process from a given program")
	fmt.Println("  kill <pid> - Kill a process with the given PID")
	fmt.Println("  ps -a - List all processes")
	fmt.Println("  ps -p <pid> - List a process with the given PID")
	fmt.Println("  dump -p <pid> - Show the PCB with the given PID")
	fmt.Println("  dump -memory <start> <end> - Show memory contents from start to end")
	fmt.Println("  exec <pid> - Execute a process with the given PID")
	fmt.Println("  log <level> - Change logging mode")
	fmt.Println("  exit - Exit the shell")
}

func handleCreate(args []string) {
	if len(args) < 1 {
		fmt.Println("Usage: create <program-name>")
		return
	}
	programName := args[0]
	pid, err := kernel.PMU.CreateProcess(programName)
	if err != nil {
		fmt.Println("Error creating process:", err)
		return
	}
	fmt.Println("Created process from program '"+programName+"' with PID:", pid)
}

func handleKill(args []string) {
	if len(args) < 1 {
		fmt.Println("Usage: kill <pid>")
		return
	}
	pid, err := strconv.Atoi(args[0])
	if err != nil {
		fmt.Println("Invalid PID:", args[0])
		return
	}
	cpu.Registers["$v0"] = pid
	cpu.InterruptionBits[2] = true
	// kernel.PMU.DestroyProcess()
}

func handlePs(args []string) {
	if len(args) < 1 {
		fmt.Println("Usage: ps <-p|-a>")
		return
	}

	switch args[0] {
	case "-a":
		for _, process := range kernel.PMU.ReadyProcesses {
			fmt.Printf("Process PID: %d, Status: %s\n", process.PID, process.Status)
		}
	case "-p":
		if len(args) < 2 {
			fmt.Println("Usage: ps -p <pid>")
			return
		}
		pid, err := strconv.Atoi(args[1])
		if err != nil {
			fmt.Println("Invalid PID:", args[1])
			return
		}
		for _, process := range kernel.PMU.ReadyProcesses {
			if process.PID == pid {
				fmt.Printf("Process PID: %d, Status: %s\n", process.PID, process.Status)
				return
			}
		}
		fmt.Println("No process found with PID:", pid)
	default:
		fmt.Println("Usage: ps <-p|-a>")
	}
}

func handleDump(args []string) {
	if len(args) < 1 {
		fmt.Println("Usage: dump <-p|-memory>")
		return
	}

	switch args[0] {
	case "-p":
		if len(args) < 2 {
			fmt.Println("Usage: dump -p <pid>")
			return
		}
		pid, err := strconv.Atoi(args[1])
		if err != nil {
			fmt.Println("Invalid PID:", args[1])
			return
		}
		for _, process := range kernel.PMU.ReadyProcesses {
			if process.PID == pid {
				fmt.Printf("Process Control Block:\n")
				fmt.Printf("  PID: %d\n", process.PID)
				fmt.Printf("  Status: %s\n", process.Status)
				fmt.Printf("  Program: %v\n", process.Program)
				fmt.Printf("  Program Length: %d\n", process.ProgramLength)
				fmt.Printf("  PC: %d\n", process.Pc)
				fmt.Printf("  Registers: %v\n", process.Registers)
				return
			}
		}
		fmt.Println("No process found with PID:", pid)

	case "-memory":
		if len(args) < 3 {
			fmt.Println("Usage: dump -memory <start> <end>")
			return
		}
		start, err1 := strconv.Atoi(args[1])
		end, err2 := strconv.Atoi(args[2])
		if err1 != nil || err2 != nil {
			fmt.Println("Invalid memory range:", args[1], args[2])
			return
		}
		if start > end || start < 0 || end >= len(Data) {
			fmt.Println("Invalid memory range:", args[1], args[2])
			return
		}
		for i := start; i <= end; i++ {
			if Data[i] != nil {
				val := reflect.ValueOf(Data[i])
				if val.Kind() == reflect.Func {
					funcName := runtime.FuncForPC(val.Pointer()).Name()
					funcName = strings.TrimPrefix(funcName, "main.")
					fmt.Printf("Address %d: %v\n", i, funcName)
				} else {
					fmt.Printf("Address %d: %v\n", i, Data[i])
				}
			} else {
				fmt.Printf("Address %d: nil\n", i)
			}
		}
	default:
		fmt.Println("Usage: dump <-p|-memory>")
	}
}

func handleExec(args []string) {
	if len(args) < 1 {
		fmt.Println("Usage: exec <pid>")
		return
	}
	pid, err := strconv.Atoi(args[0])
	if err != nil {
		fmt.Println("Invalid PID:", args[0])
		return
	}
	err = kernel.PMU.Execute(pid)
	if err != nil {
		fmt.Println("Error executing process:", err)
	} else {
		fmt.Println("Executing process with PID:", pid)
	}
}

func handleLog(args []string) {
	if len(args) < 1 {
		fmt.Println("Usage: log <info|trace|debug>")
		return
	}
	level := args[0]
	switch level {
	case "info":
		LogLevel = "INFO"
	case "trace":
		LogLevel = "TRACE"
	case "debug":
		LogLevel = "DEBUG"
	default:
		fmt.Println("Unknown log level:", level)
		return
	}
	fmt.Println("Log level set to:", LogLevel)
}
