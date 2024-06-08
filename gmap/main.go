package main

import (
	"fmt"
	"gmap/utils"
	"os/signal"
    "syscall"
)

func main() {
	// SIGINT handling 
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func(){
		<-c 
		fmt.Printf("%s%s%s\n", utils.RED , "[*] Exiting...", Reset)
		os.Exit(1)
	}()


	// Show banner
	printBanner()

	// Check arguments
	args := parseArguments()

	if args.Help {
		printHelp()
		return
	}

	// Check that target is provided
	if args.Target == "" {
		utils.PrintError("[ERROR] target must be provided")
		printHelp()
		return
	}

	// Parse ports
	ports, err := parsePorts(args.Ports)
	if err != nil {
		printHelp()
		return
	}

	// Parse target
	target, err := parseTarget(args.Target)
	if err != nil {
		printHelp()
		return
	}

	// Parse Scan Type 
	scanType, err := parseScanType(args.ScanType)
	if err != nil {
		printHelp()
		return 
	}


	// Validate output and format if output flag is set
	if outputFlagSet {
		if err := parseFormat(args.OutputFile, args.Format); err != nil {
			printHelp()
			return
		}
	}

	// Set the scan parameters 
	scanParams := utils.ScanParameters{
		Target:  target,
		Ports:   ports,
		Timeout: args.Timeout, 
	}


	var results[] Port 

	// Perform Scan 
	switch scanType {
	// Perform UDP Scan 
	case "udp":
		results = scanner.udpScan(scanParams)
	// Perform TCP Scan 
	case "tcp":
		results = scanner.tcpScan(scanParams)
	}

	// Export results if necessary 
	if args.OutputFile != "" {
		if err := utils.ExportResults(results, args.OutputFile, args.Format); err != nil {
			fmt.Println(err)
			printHelp()
		}
	}
	
	// Succesfull exit 
	os.Exit(0)
}
