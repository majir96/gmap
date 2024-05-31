package main

import (
	"flag"
	"fmt"
	"strconv"
	"strings"
)

type Arguments struct {
	Help   bool
	Ports  string
	Target string
	// TODO ADD MORE OPTIONS
}

const CommonPorts string = "20,21,22,23,25,53,67,68,69,80,110,119,123,135,137,138,139,143,161,162,179,194,443,445,465,587,993,995,1433,3306"

func parseArguments() Arguments {
	var args Arguments

	// Define arguments
	flag.BoolVar(&args.Help, "h", false, "Display help panel")
	flag.BoolVar(&args.Help, "help", false, "Display help panel")

	flag.StringVar(&args.Ports, "p", CommonPorts, "Port(s) to scan")
	flag.StringVar(&args.Ports, "port", CommonPorts, "Port(s) to scan")
	flag.StringVar(&args.Ports, "p-", "0-65535", "Scan all ports (0-65535)")

	flag.StringVar(&args.Target, "t", "", "Target to scan")
	flag.StringVar(&args.Target, "target", "", "Target to scan")

	return args

}

func castSlice(portString []string) ([]int, error) {
	return nil, nil
}

func generateRange(start int, end int) []int {

	if start > end || end > 65535 || start < 0 || end < 0 {
		// TODO control de errores
		return nil
	}

	portRange := end - start + 1

	ports := make([]int, portRange)

	for i := 0; i <= end; i++ {
		ports[i] = i
	}

	return ports

}

func parsePorts(portString string) []int {
	var ports []int
	var portList []string

	// Parse a range
	if strings.Contains(portString, "-") {
		portList = strings.Split(portString, "-")

		// Generate range of ports
		if len(portList) == 2 {
			start, err1 := strconv.Atoi(portList[0])
			end, err2 := strconv.Atoi(portList[1])

			if err1 != nil || err2 != nil {
				fmt.Println("[ERROR] Arguments must be Integers")
				printHelp()
				return nil
			}

			ports = generateRange(start, end)

		} else {
			fmt.Println("[ERROR] Just 2 arguments are needed for ports")
			printHelp()
			return nil
		}

	} else if strings.Contains(portString, ",") { // Parse a list
		portList := strings.Split(portString, ",")

		ports, _ = castSlice(portList)

	} else { // Parse a single port
		port, err := strconv.Atoi(portString)

		if err == nil {
			ports = append(ports, port)
		}

	}

	return ports
}

func printHelp() {
	fmt.Println("Help panel for gomap:")
	fmt.Println(Lines)
	fmt.Println("Usage")
	fmt.Println("./gomap -t <ip> -p <port(s)>")
	fmt.Println(Lines)
	fmt.Println("Options:")
	fmt.Printf("  -p, --port    Port(s) to scan. Default set to %s\n", CommonPorts)
	fmt.Println("If various ports are to be scanned separate by commas, i.e -p 22,23")
	fmt.Println("If a range is to be scanned separate by hyphen, i.e -p 0-400")
	fmt.Println(" -p-           All ports are to be scanned 0-65535")
	fmt.Println("  -t, --target    Target to scan (required)")
	fmt.Println("  -h, --help      Display this help message")
	fmt.Println(Lines)
	fmt.Println("Example of use:")
	fmt.Println("./gomap -t 127.0.0.1 -p 0-65535")
	fmt.Println(Lines)

}
