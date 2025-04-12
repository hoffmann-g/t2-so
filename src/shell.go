package main

import (
	"bufio"
	"fmt"
	"os"
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
		tokens := strings.Split(line, " ")

		switch tokens[0] {
		case "help":
			fmt.Println("Available commands:")
			fmt.Println("  help - Show this help message")
			fmt.Println("  create <program-name> - Create a new process from a given program")
			fmt.Println("  kill <pid> - Kill a process with the given PID")
			fmt.Println("  dump -a - Dump all processes")
			fmt.Println("  dump -p <pid> - Dump a process with the given PID")
			fmt.Println("  exec <pid> - Execute a process with the given PID")
			fmt.Println("  log <level> - Change logging mode")
			fmt.Println("  exit - Exit the shell")

		case "create":
			if len(tokens) < 2 {
				fmt.Println("Usage: create <program-name>")
				break
			}
			programName := tokens[1]

			pid, err := kernel.PMU.CreateProcess(programName)
			if err != nil {
				fmt.Println("Error creating process:", err)
				break
			}
			fmt.Println("Created process from program '"+programName+"' with PID:", pid)

		case "kill":
			if len(tokens) < 2 {
				fmt.Println("Usage: kill <pid>")
				break
			}
			pid := tokens[1]
			pidInt, convErr := strconv.Atoi(pid)
			if convErr != nil {
				fmt.Println("Invalid PID:", pid)
				break
			}

			cpu.Registers["$v0"] = pidInt
			kernel.PMU.DestroyProcess()
		case "dump":
			if len(tokens) < 2 {
				fmt.Println("Usage: dump <-p><-a>")
				break
			}

			switch tokens[1] {
			case "-a":
				for _, process := range kernel.PMU.ReadyProcesses {
					fmt.Printf("Process PID: %d, Status: %s\n", process.PID, process.Status)
				}
			case "-p":
				if len(tokens) < 3 {
					fmt.Println("Usage: dump -p <pid>")
					break
				}
				pid := tokens[2]

				fmt.Println("Dumping process with PID:", pid)
			}

		case "exec":
			if len(tokens) < 2 {
				fmt.Println("Usage: exec <pid>")
				break
			}

			pid := tokens[1]
			pidInt, convErr := strconv.Atoi(pid)
			if convErr != nil {
				fmt.Println("Invalid PID:", pid)
				break
			}
			err := kernel.PMU.Execute(pidInt)
			if err == nil {
				fmt.Println("Executing process with PID:", pid)
			} else {
				fmt.Println("Error executing process:", err)
			}

		case "log":
			if len(tokens) < 2 {
				fmt.Println("Usage: log <info|trace|debug>")
				break
			}
			level := tokens[1]
			found := true
			switch level {
			case "info":
				LogLevel = "INFO"
			case "trace":
				LogLevel = "TRACE"
			case "debug":
				LogLevel = "DEBUG"
			default:
				fmt.Println("Unknown log level:", level)
				found = false
			}
			if found {
				fmt.Println("Log level set to:", LogLevel)
			}

		case "exit":
			return

		default:
			fmt.Println("Unknown command:", tokens[0])
		}

		fmt.Print("> ")
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading shell input:", err)
	}
}
