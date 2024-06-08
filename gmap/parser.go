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


var outputFlagSet bool 

func parseArguments() utils.Arguments {
	var args utils.Arguments

	// Define arguments
	flag.BoolVar(&args.Help, "h", false, "Display help panel")
	flag.BoolVar(&args.Help, "help", false, "Display help panel")

	flag.StringVar(&args.Output, "o", "", "Export output to file")
	flag.StringVar(&args.Output, "output", "", "Export output to file")

	flag.BoolVar(&args.open, "open", false, "Show only opened ports")

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
	flag.StringVar(&timeout, "timeout", "1s","Delaty timeout for packets being sent (e.g., 500ms, 2s, 1m)")

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


func parseFormat(output string, format string) error{

	if output == "" {
		return utils.PrintError("[ERROR] output filename must be provided")
	}
	
	if format != "txt" && format != "json" && format != "csv"{
		return utils.PrintError("[ERROR] unsupported file provided")
	}

	return nil 
}

func parseScanType(scan string) string, error {
	
	if scan != "udp" && scan != "tcp" {
		return "", utils.PrintError("[ERROR] unsupported scan type")
	}

	return scan, nil 
} 


func printHelp() {
	fmt.Println("Help panel for gomap:")
	fmt.Println(utils.Lines)
	fmt.Println("Usage")
	fmt.Println("./gomap -t <IP> -p <PORTS> -o ")
	fmt.Println(utils.Lines)
	fmt.Println("Options:")
	fmt.Printf("   -p, --port  <PORTS>   Port(s) to scan. Default set to %s\n", utils.CommonPorts)
	fmt.Println("If various ports are to be scanned separate by commas, i.e -p 22,23")
	fmt.Println("If a range is to be scanned separate by hyphen, i.e -p 0-400")
	fmt.Println("Services will automatically be scanned or obtained for all ports")
	fmt.Println("  -p-              	 All ports are to be scanned 0-65535")
	fmt.Println("  -t, --target <IP>     Target to scan (required)")
	fmt.Println("  -s, --scan 	<SCAN>	 Type of scan to perform. Options:")
	fmt.Println("  tcp: 				 Perform a TCP Scan (default)")
	fmt.Println("  udp: 				 Perform a UDP Scan ")
	fmt.Println("  -h, --help       	 Display this help message")
	fmt.Println("  -o, --output <FILE>   Export output to file, default format .txt")
	fmt.Println("  -f, --format <FORMAT> Format to export the file to. Formats:")
	fmt.Println("  txt: 				 Export to text file (default)")
	fmt.Println("  csv: 				 Export to csv file")
	fmt.Println("  json: 				 Export to json file")
	fmt.Println("  --open 				 Filter by open ports on output ")
	fmt.Println(" --timeout <TIMEOUT>	 Timeout to be set for packets when scanning (e.g., 500ms, 2s, 1m)")
	fmt.Println(utils.Lines)
	fmt.Println("Example of use:")
	fmt.Println("./gomap -t 127.0.0.1 -p 0-65535 -o test")
	fmt.Println(utils.Lines)

}
