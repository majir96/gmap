package main

import (
	"flag"
	"fmt"
	"time"
	"gmap/utils"
	"regexp"
	"strconv"
	"strings"
)

func parseArguments() utils.Arguments {
	var args utils.Arguments

	// Define arguments
	flag.BoolVar(&args.Help, "h", false, "Display help panel")
	flag.BoolVar(&args.Help, "help", false, "Display help panel")

	flag.BoolVar(&args.Output, "o", false, "Export output to file")
	flag.BoolVar(&args.Output, "output", false, "Export output to file")

	flag.BoolVar(&args.open, "open", false, "Show only opened ports")

	flag.StringVar(&args.Ports, "p", utils.CommonPorts, "Port(s) to scan")
	flag.StringVar(&args.Ports, "port", utils.CommonPorts, "Port(s) to scan")
	flag.StringVar(&args.Ports, "p-", "0-65535", "Scan all ports (0-65535)")

	flag.StringVar(&args.Target, "t", "", "Target to scan")
	flag.StringVar(&args.Target, "target", "", "Target to scan")

	var timeout string 
	flag.StringVar(&timeout, "timeout", "1s","Delaty timeout for packets being sent (e.g., 500ms, 2s, 1m)")

	flag.Parse()

	parsedTimeout, err := time.ParseDuration(timeout)

	if err != nil {
		fmt.Println(utils.PrintError(fmt.Sprintf("Invalid timeout value: %s, defaulting to 2s", timeout)))
		parsedTimeout = 2 * time.Second
	}

	args.Timeout = parsedTimeout


	return args

}

func castSlice(portString []string) ([]int, error) {
	portInt := make([]int, len(portString))

	for i, p := range portString {
		port, err := strconv.Atoi(p)

		// Check for errors
		if err != nil {
			return nil, utils.PrintError(fmt.Sprintf("arguments must be integers %s", p))
		}

		portInt[i] = port
	}

	return portInt, nil
}

func generateRange(start int, end int) ([]int, error) {

	// Check for errors in ranges
	if start > end || end > 65535 || start < 0 || end < 0 {
		return nil, utils.PrintError(fmt.Sprintf("invalid port range: %d - %d", start, end))
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
				return nil, utils.PrintError(fmt.Sprintf("range arguments must be integers: %d - %d", start, end))
			}

			// Generate slice between ranges start and end
			return generateRange(start, end)
		} else {
			return nil, utils.PrintError(fmt.Sprintf("exactly 2 arguments are needed for port range"))
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
				return nil, utils.PrintError(fmt.Sprintf("arguments must be Integers %s", portString))
			}

			return []int{port}, nil
		}

		// Unhandled cases
		return nil, utils.PrintError(fmt.Sprintf("invalid format provided for ports at: %s", portString))
	}

}

func parseTarget(targetString string) (string, error) {
	// Regex for IPv4
	const pattern = `^((25[0-5]|2[0-4][0-9]|1[0-9]{2}|[1-9]?[0-9])\.){3}(25[0-5]|2[0-4][0-9]|1[0-9]{2}|[1-9]?[0-9])$`

	re := regexp.MustCompile(pattern)

	if !re.MatchString(targetString) {
		return "", utils.PrintError(fmt.Sprintf("invalid target host %s", targetString))
	}

	return targetString, nil
}

func printHelp() {
	fmt.Println("Help panel for gomap:")
	fmt.Println(utils.Lines)
	fmt.Println("Usage")
	fmt.Println("./gomap -t <IP> -p <PORTS> -o ")
	fmt.Println(utils.Lines)
	fmt.Println("Options:")
	fmt.Printf("   -p, --port  <PORTS>  Port(s) to scan. Default set to %s\n", utils.CommonPorts)
	fmt.Println("If various ports are to be scanned separate by commas, i.e -p 22,23")
	fmt.Println("If a range is to be scanned separate by hyphen, i.e -p 0-400")
	fmt.Println("Services will automatically be scanned or obtained for all ports")
	fmt.Println("  -p-              	All ports are to be scanned 0-65535")
	fmt.Println("  -t, --target <IP>    Target to scan (required)")
	fmt.Println("  -h, --help       	Display this help message")
	fmt.Println("  -o, --output <FILE>  Export output to file, default format .txt")
	fmt.Println("  --open 				Filter by open ports on output ")
	fmt.Println(" --timeout <TIMEOUT>	Timeout to be set for packets when scanning (e.g., 500ms, 2s, 1m)")
	fmt.Println(utils.Lines)
	fmt.Println("Example of use:")
	fmt.Println("./gomap -t 127.0.0.1 -p 0-65535 -o test")
	fmt.Println(utils.Lines)

}
