package main

import (
	"flag"
	"fmt"
	"gmap/utils"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var outputFlagSet bool

func parseArguments() utils.Arguments {
	var args utils.Arguments

	// Define arguments
	flag.BoolVar(&args.Help, "h", false, "Display help panel")
	flag.BoolVar(&args.Help, "help", false, "Display help panel")

	flag.StringVar(&args.Output, "o", "", "Export output to file")
	flag.StringVar(&args.Output, "output", "", "Export output to file")

	flag.BoolVar(&args.Open, "open", false, "Show only opened ports")

	flag.StringVar(&args.Ports, "p", utils.CommonPorts, "Port(s) to scan")
	flag.StringVar(&args.Ports, "port", utils.CommonPorts, "Port(s) to scan")
	flag.StringVar(&args.Ports, "p-", "0-65535", "Scan all ports (0-65535)")

	flag.StringVar(&args.Target, "t", "", "Target to scan")
	flag.StringVar(&args.Target, "target", "", "Target to scan")

	flag.StringVar(&args.ScanType, "s", "tcp", "Type of scan to perform")
	flag.StringVar(&args.ScanType, "scan", "tcp", "Type of scan to perform")

	flag.StringVar(&args.Format, "f", ".txt", "Format to export the file to, default to txt")
	flag.StringVar(&args.Format, "format", ".txt", "Format to export the file to, default to txt")

	var timeout string
	flag.StringVar(&timeout, "timeout", "1s", "Delaty timeout for packets being sent (e.g., 500ms, 2s, 1m)")

	flag.Parse()

	// Parse and check if timeout format is correct
	parsedTimeout, err := time.ParseDuration(timeout)
	if err != nil {
		fmt.Println(utils.PrintError(fmt.Sprintf("Invalid timeout value: %s, defaulting to 2s", timeout)))
		parsedTimeout = 2 * time.Second
	}

	args.Timeout = parsedTimeout

	// Check if the output flag was explicitly set by the user
	flag.Visit(func(f *flag.Flag) {
		if f.Name == "o" || f.Name == "output" {
			outputFlagSet = true
		}
	})

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
			return nil, utils.PrintError(fmt.Sprintf("exactly 2 arguments are needed for port range, %d provided", len(portList)))
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

func parseFormat(output string, format string) error {

	if output == "" {
		return utils.PrintError("[ERROR] output filename must be provided")
	}

	if format != "txt" && format != "json" && format != "csv" {
		return utils.PrintError("[ERROR] unsupported file provided")
	}

	return nil
}

func parseScanType(scan string) (string, error) {

	if scan != "udp" && scan != "tcp" {
		return "", utils.PrintError("[ERROR] unsupported scan type")
	}

	return scan, nil
}

func printHelp() {
	fmt.Println("Help panel for gomap:")
	fmt.Println(utils.Lines)
	fmt.Printf("%sUsage%s\n", utils.LightGreen, utils.Reset)
	fmt.Printf("%s./gomap -t <IP> -p <PORTS> -o %s\n", utils.LightGreen, utils.Reset)
	fmt.Println(utils.Lines)
	fmt.Printf("%sOptions:%s\n", utils.LightGreen, utils.Reset)
	fmt.Printf("  %s-p, --port <PORTS>        %s%sPort(s) to scan. Default set to %s%s\n", utils.LightGreen, utils.Reset, utils.BrightWhite, utils.CommonPorts, utils.Reset)
	fmt.Printf("                            %sIf various ports are to be scanned, separate by commas, e.g., -p 22,23%s\n", utils.LightGreen, utils.Reset)
	fmt.Printf("                            %sIf a range is to be scanned, separate by hyphen, e.g., -p 0-400%s\n", utils.LightGreen, utils.Reset)
	fmt.Printf("                            %sServices will automatically be scanned or obtained for all ports%s\n", utils.LightGreen, utils.Reset)
	fmt.Printf("  %s-p-                       All ports are to be scanned 0-65535%s\n", utils.LightGreen, utils.Reset)
	fmt.Printf("  %s-t, --target <IP>         Target to scan (required)%s\n", utils.LightGreen, utils.Reset)
	fmt.Printf("  %s-s, --scan <SCAN>         Type of scan to perform. Options:%s\n", utils.LightGreen, utils.Reset)
	fmt.Printf("                            %stcp: Perform a TCP Scan (default)%s\n", utils.LightGreen, utils.Reset)
	fmt.Printf("                            %sudp: Perform a UDP Scan%s\n", utils.LightGreen, utils.Reset)
	fmt.Printf("  %s-h, --help                Display this help message%s\n", utils.LightGreen, utils.Reset)
	fmt.Printf("  %s-o, --output <FILE>       Export output to file, default format .txt%s\n", utils.LightGreen, utils.Reset)
	fmt.Printf("  %s-f, --format <FORMAT>     Format to export the file to. Formats:%s\n", utils.LightGreen, utils.Reset)
	fmt.Printf("                            %stxt: Export to text file (default)%s\n", utils.LightGreen, utils.Reset)
	fmt.Printf("                            %scsv: Export to csv file%s\n", utils.LightGreen, utils.Reset)
	fmt.Printf("                            %sjson: Export to json file%s\n", utils.LightGreen, utils.Reset)
	fmt.Printf("  %s--open                    Filter by open ports on output%s\n", utils.LightGreen, utils.Reset)
	fmt.Printf("  %s--timeout <TIMEOUT>       Timeout to be set for packets when scanning (e.g., 500ms, 2s, 1m)%s\n", utils.LightGreen, utils.Reset)
	fmt.Println(utils.Lines)
	fmt.Printf("%sExample of use:%s\n", utils.LightGreen, utils.Reset)
	fmt.Printf("%s./gomap -t 127.0.0.1 -p 0-65535 -o test%s\n", utils.LightGreen, utils.Reset)
	fmt.Println(utils.Lines)
}
