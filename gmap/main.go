package main

import (
	"fmt"
	"gmap/utils"
)

func main() {
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
		utils.PrintError("target must be provided")
		printHelp()
		return
	}

	// Parse ports
	ports, err := parsePorts(args.Ports)
	if err != nil {
		utils.PrintError(fmt.Sprintf("[ERROR] %v", err))
		printHelp()
		return
	}

	// Parse target
	target, err := parseTarget(args.Target)
	if err != nil {
		utils.PrintError(fmt.Sprintf("[ERROR] %v", err))
		printHelp()
		return
	}

	// TODO Scan
}
