package main

import (
	"flag"
	"fmt"
	"regexp"
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

	flag.Parse()

	return args

}

func castSlice(portString []string) ([]int, error) {
	portInt := make([]int, len(portString))

	for i, p := range portString {
		port, err := strconv.Atoi(p)

		// Check for errors
		if err != nil {
			return nil, fmt.Errorf("arguments must be integers %s", p)
		}

		portInt[i] = port
	}

	return portInt, nil
}

func generateRange(start int, end int) ([]int, error) {

	// Check for errors in ranges
	if start > end || end > 65535 || start < 0 || end < 0 {
		return nil, fmt.Errorf("invalid port range: %d - %d", start, end)
	}

	ports := make([]int, end-start+1)

	// Generate range list
	for i := start; i <= end; i++ {
		ports[i-start] = i
	}

	return ports, nil

}

func parsePorts(portString string) ([]int, error) {
	var portList []string

	switch {
	// Parse a range of ports
	case strings.Contains(portString, "-"):
		portList = strings.Split(portString, "-")

		// Check for errors
		if len(portList) == 2 {
			start, err1 := strconv.Atoi(portList[0])
			end, err2 := strconv.Atoi(portList[1])

			if err1 != nil || err2 != nil {
				return nil, fmt.Errorf("range arguments must be integers: %d - %d", start, end)
			}

			// Generate slice between ranges start and end
			return generateRange(start, end)
		} else {
			return nil, fmt.Errorf("exactly 2 arguments are needed for port range")
		}

		// Parse a list of ports
	case strings.Contains(portString, ","):
		portList := strings.Split(portString, ",")
		// Cast the string slice to int slice
		return castSlice(portList)

	// Check if single character is integer and do
	default:
		if _, err := strconv.Atoi(portString); err == nil {
			port, err := strconv.Atoi(portString)

			if err != nil {
				return nil, fmt.Errorf("arguments must be Integers %s", portString)
			}

			return []int{port}, nil
		}

		// Unhandled cases
		return nil, fmt.Errorf("invalid format provided for ports at: %s", portString)
	}

}

func parseTarget(targetString string) (string, error) {
	// Regex for IPv4
	const pattern = `^((25[0-5]|2[0-4][0-9]|1[0-9]{2}|[1-9]?[0-9])\.){3}(25[0-5]|2[0-4][0-9]|1[0-9]{2}|[1-9]?[0-9])$`

	re := regexp.MustCompile(pattern)

	if !re.MatchString(targetString) {
		return "", fmt.Errorf("invalid target host %s", targetString)
	}

	return targetString, nil
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
