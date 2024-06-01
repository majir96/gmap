package main

import "fmt"

/*

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"time"
	"gmap/scanner"
	"strconv"
	"bufio"
	"encoding/xml"
	"os/exec"
	"errors"
)


*/

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
		fmt.Errorf("target must be provided")
		printHelp()
		return
	}

	// Parse ports
	ports, err := parsePorts(args.Ports)
	if err != nil {
		fmt.Printf("[ERROR] %v", err)
		printHelp()
		return
	}

	// Parse target
	target, err := parseTarget(args.Target)
	if err != nil {
		fmt.Printf("[ERROR] %v", err)
		printHelp()
		return
	}

	// TODO Scan
}
