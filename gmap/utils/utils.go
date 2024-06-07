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

// Common Services 
const commonServices = map[int]string{
	20:   "FTP Data",
	21:   "FTP Control",
	22:   "SSH",
	23:   "Telnet",
	25:   "SMTP",
	53:   "DNS",
	67:   "DHCP Server",
	68:   "DHCP Client",
	69:   "TFTP",
	80:   "HTTP",
	110:  "POP3",
	119:  "NNTP",
	123:  "NTP",
	135:  "MS RPC",
	137:  "NetBIOS Name Service",
	138:  "NetBIOS Datagram Service",
	139:  "NetBIOS Session Service",
	143:  "IMAP",
	161:  "SNMP",
	162:  "SNMP Trap",
	179:  "BGP",
	194:  "IRC",
	443:  "HTTPS",
	445:  "Microsoft-DS",
	465:  "SMTPS",
	514:  "Syslog",
	515:  "LPD",
	587:  "Submission",
	631:  "IPP",
	636:  "LDAPS",
	993:  "IMAPS",
	995:  "POP3S",
	1080: "SOCKS",
	1194: "OpenVPN",
	1433: "MSSQL",
	1434: "MSSQL Monitor",
	1521: "Oracle",
	1723: "PPTP",
	1883: "MQTT",
	1900: "SSDP",
	2049: "NFS",
	2082: "cPanel",
	2083: "cPanel Secure",
	2483: "Oracle DB Secure",
	2484: "Oracle DB",
	3306: "MySQL",
	3389: "RDP",
	3690: "Subversion",
	4444: "Metasploit",
	4848: "GlassFish Admin",
	5432: "PostgreSQL",
	5632: "pcAnywhere",
	5900: "VNC",
	5984: "CouchDB",
	6379: "Redis",
	8000: "Common HTTP Alt",
	8080: "HTTP Proxy",
	8086: "InfluxDB",
	8181: "HTTP Proxy",
	8443: "HTTPS Alt",
	8888: "HTTP Proxy",
	9000: "SonarQube",
	9092: "Kafka",
	9200: "Elasticsearch",
	9300: "Elasticsearch",
	11211: "Memcached",
	27017: "MongoDB",
	27018: "MongoDB",
	27019: "MongoDB",
	50000: "SAP",
	50070: "Hadoop",
}


// Type definitions
type Arguments struct {
	Help   bool
	Ports  string
	Target string
	Output bool // TODO : -> Write to file 
	Open bool 
	Timeout time.Duration 
	// TODO ADD MORE OPTIONS
	/**
	NOTE: Options to filter by 
		-sS 
		-nmap 
		--min-rate 
		--timeout 
		--ip-range 
		-vvv
	
	*/
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
