package utils

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// Constant vALUES
const (
	Red          = "\033[31m"
	Green        = "\033[32m"
	Reset        = "\033[0m"
	Purple       = "\033[35m"
	Blue         = "\033[34m"
	LightGreen   = "\033[92m"
	DarkGreen    = "\033[32m"
	Yellow       = "\033[33m"
	Cyan         = "\033[36m"
	Magenta      = "\033[35m"
	White        = "\033[97m"
	Black        = "\033[30m"
	BrightRed    = "\033[91m"
	BrightBlue   = "\033[94m"
	BrightPurple = "\033[95m"
	BrightCyan   = "\033[96m"
	BrightYellow = "\033[93m"
	BrightWhite  = "\033[97m"
	Lines        = "--------------------------------"
	CommonPorts  = "20,21,22,23,25,53,67,68,69,80,110,119,123,135,137,138,139,143,161,162,179,194,443,445,465,587,993,995,1433,3306"
)

// Common Services
var CommonServices = map[int]string{
	20:    "FTP Data",
	21:    "FTP Control",
	22:    "SSH",
	23:    "Telnet",
	25:    "SMTP",
	53:    "DNS",
	67:    "DHCP Server",
	68:    "DHCP Client",
	69:    "TFTP",
	80:    "HTTP",
	110:   "POP3",
	119:   "NNTP",
	123:   "NTP",
	135:   "MS RPC",
	137:   "NetBIOS Name Service",
	138:   "NetBIOS Datagram Service",
	139:   "NetBIOS Session Service",
	143:   "IMAP",
	161:   "SNMP",
	162:   "SNMP Trap",
	179:   "BGP",
	194:   "IRC",
	443:   "HTTPS",
	445:   "Microsoft-DS",
	465:   "SMTPS",
	514:   "Syslog",
	515:   "LPD",
	587:   "Submission",
	631:   "IPP",
	636:   "LDAPS",
	993:   "IMAPS",
	995:   "POP3S",
	1080:  "SOCKS",
	1194:  "OpenVPN",
	1433:  "MSSQL",
	1434:  "MSSQL Monitor",
	1521:  "Oracle",
	1723:  "PPTP",
	1883:  "MQTT",
	1900:  "SSDP",
	2049:  "NFS",
	2082:  "cPanel",
	2083:  "cPanel Secure",
	2483:  "Oracle DB Secure",
	2484:  "Oracle DB",
	3306:  "MySQL",
	3389:  "RDP",
	3690:  "Subversion",
	4444:  "Metasploit",
	4848:  "GlassFish Admin",
	5432:  "PostgreSQL",
	5632:  "pcAnywhere",
	5900:  "VNC",
	5984:  "CouchDB",
	6379:  "Redis",
	8000:  "Common HTTP Alt",
	8080:  "HTTP Proxy",
	8086:  "InfluxDB",
	8181:  "HTTP Proxy",
	8443:  "HTTPS Alt",
	8888:  "HTTP Proxy",
	9000:  "SonarQube",
	9092:  "Kafka",
	9200:  "Elasticsearch",
	9300:  "Elasticsearch",
	11211: "Memcached",
	27017: "MongoDB",
	27018: "MongoDB",
	27019: "MongoDB",
	50000: "SAP",
	50070: "Hadoop",
}

// Type definitions
type Arguments struct {
	Help          bool
	Ports         string
	Target        string
	Output        string
	Open          bool
	Timeout       time.Duration
	Format        string
	ScanType      string
	HostDiscovery bool
	// TODO ADD MORE OPTIONS
	/**
	NOTE: Options to filter by
		-nmap
		--min-rate
		--ip-range
		-vvv

	*/
}

type Port struct {
	Port    int
	Status  string
	Service string
}

type ScanParameters struct {
	Target  string
	Ports   []int
	Timeout time.Duration
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

func exportToTxt(results []Port, file *os.File) error {
	// Dump results
	for _, result := range results {
		line := fmt.Sprintf("Port: %d, Status: %s, Service: %s\n", result.Port, result.Status, result.Service)
		if _, err := file.WriteString(line); err != nil {
			return fmt.Errorf("could not write to file: %v", err)
		}
	}

	PrintSuccess("[!] Results successfully exported to .txt file")
	return nil
}

func exportToCsv(results []Port, file *os.File) error {

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write header
	header := []string{"Port", "Status", "Service"}
	if err := writer.Write(header); err != nil {
		return PrintError(fmt.Sprintf("[ERROR] could not write header to file: %v", err))
	}

	// Dump results
	for _, result := range results {
		record := []string{fmt.Sprintf("%d", result.Port), result.Status, result.Service}

		if err := writer.Write(record); err != nil {
			return PrintError(fmt.Sprintf("[ERROR] could not write record to file: %v", err))
		}
	}

	PrintSuccess("[!] Results successfully exported to .csv file")
	return nil
}

func exportToJson(results []Port, file *os.File) error {

	encoder := json.NewEncoder(file)
	if err := encoder.Encode(results); err != nil {
		return PrintError(fmt.Sprintf("[ERROR] could not encode results to JSON: %v", err))
	}

	PrintSuccess("[!] Results successfully exported to .json file")
	return nil
}

// Export to file
func ExportResults(results []Port, file string, format string) error {

	var fileName string = fmt.Sprintf("%s.%s", file, format)

	// Create file
	f, err := os.Create(fileName)
	if err != nil {
		return PrintError(fmt.Sprintf("[ERROR] could not create file: %v", err))
	}

	defer f.Close()

	// Handle formats to export to
	switch format {
	case "txt":
		return exportToTxt(results, f)
	case "csv":
		return exportToCsv(results, f)
	case "json":
		return exportToJson(results, f)
	}

	return nil
}
