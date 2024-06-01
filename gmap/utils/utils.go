package utils

import (
	"fmt"
	"time"
)

// Lines Const
const Lines string = "--------------------------------"

// Color codes
const Red = "\033[31m"
const Green = "\033[32m"
const Reset = "\033[0m"
const Purple = "\033[35m"
const Blue = "\033[34m"

// Common Ports
const CommonPorts string = "20,21,22,23,25,53,67,68,69,80,110,119,123,135,137,138,139,143,161,162,179,194,443,445,465,587,993,995,1433,3306"

// Type definition
type Arguments struct {
	Help   bool
	Ports  string
	Target string
	Output bool
	// TODO ADD MORE OPTIONS
}

type Port struct {
	port    int
	status  string
	service string
}

type ScanParameters struct {
	target    string
	ports     []int
	scan_type string
	timeout   time.Duration
}

// Auxiliary functions

// Print error in red
func PrintError(msg string) error {
	return fmt.Errorf("%s%s%s", Red, msg, Reset)
}

// Print success in green
func PrintSuccess(msg string) {
	fmt.Printf("%s%s%s\n", Green, msg, Reset)
}
