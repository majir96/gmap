package scanner

import (
	"fmt"
	"gmap/utils"
	"net"
	"sync"
	"time"

	"github.com/go-ping/ping"
)

// Check Availability of host
func HostUp(target string, timeout time.Duration) bool {
	pinger, err := ping.NewPinger(target)

	if err != nil {
		utils.PrintError(fmt.Sprintf("[ERROR] Failed to create pinger for target %s", target))
		return false
	}

	// Establish parameters for pinger
	pinger.Count = 3
	pinger.Timeout = timeout
	pinger.SetPrivileged(true)

	// Block until finished
	err = pinger.Run()

	if err != nil {
		utils.PrintError(fmt.Sprintf("[ERROR] Ping failed for target %s", target))
		return false
	}

	stats := pinger.Statistics()

	return stats.PacketsRecv > 0
}

// Auxiliary function to check if service is known
func checkService(service string) string {

	if service == "" {
		service = "unknown"
	}

	return service
}

// Auxiliary function to count all opened and filtered ports
func countOpenPorts(results []utils.Port) int {
	count := 0

	for _, result := range results {
		if result.Status == "open" || result.Status == "open/filtered" {
			count++
		}
	}

	return count
}

// Auxiliary function for TCP banner grab for service detection
func bannerGrab(conn net.Conn) string {
	// Set a short read deadline for banner grabbing
	conn.SetReadDeadline(time.Now().Add(2 * time.Second))

	buffer := make([]byte, 1024)
	n, err := conn.Read(buffer)

	if err != nil {
		return ""
	}

	return string(buffer[:n])
}

// TCP Worker for go routine Multithreading
func tcpWorker(target string, port int, timeout time.Duration, results chan<- utils.Port, wg *sync.WaitGroup) {
	// Defer the call to Done to ensure the WaitGroup counter is decremented when the function completes
	defer wg.Done()

	// Format address string
	address := fmt.Sprintf("%s:%d", target, port)

	// Try to establish connection
	conn, err := net.DialTimeout("tcp", address, timeout)

	var state, service string

	if err != nil {
		// Check wether the port is closed or filtered based on error
		if opErr, ok := err.(*net.OpError); ok && opErr.Err.Error() == "connect: connection refused" {
			state = "closed"
		} else {
			state = "filtered"
		}
		service = utils.CommonServices[port]
	} else { // If no error, port is opened
		// We try to check for banner to determien service
		service = bannerGrab(conn)
		state = "open"
		conn.Close()
	}

	// Check service
	service = checkService(service)

	// Send back the result through results channel
	results <- utils.Port{Port: port, Status: state, Service: service}
}

// UDP Worker for go routine multithreading
func udpWorker(target string, port int, timeout time.Duration, results chan<- utils.Port, wg *sync.WaitGroup) {
	// Defer the call to Done to ensure the WaitGroup counter is decremented when the function completes
	defer wg.Done()

	// Format address
	address := fmt.Sprintf("%s:%d", target, port)

	// Try to establish connection
	conn, err := net.DialTimeout("udp", address, timeout)

	var state, service string
	// If an error occurs the port is closed
	if err != nil {
		service = utils.CommonServices[port]
		results <- utils.Port{Port: port, Status: "closed", Service: service}
		return
	}
	// Ensure the connection is closed
	defer conn.Close()

	// Send a ping to check if port is closed
	_, err = conn.Write([]byte("Ping"))
	if err != nil {
		results <- utils.Port{Port: port, Status: "closed", Service: ""}
		return
	}

	// Set a read timeline for the response
	conn.SetReadDeadline(time.Now().Add(timeout))
	buff := make([]byte, 1024)
	n, err := conn.Read(buff)

	if err != nil {
		// If no response, port is either opened or filtered
		service = utils.CommonServices[port]
		state = "open/filtered"
	} else {
		// If there is response, port is opened
		service = string(buff[:n])
		state = "open"
	}

	// Check service
	service = checkService(service)

	// Send back the result through results channel
	results <- utils.Port{Port: port, Status: state, Service: service}
}

// Function to perform a basic TCP Scan
func TcpScan(scan utils.ScanParameters) []utils.Port {
	var results []utils.Port
	resultChan := make(chan utils.Port, len(scan.Ports))
	var wg sync.WaitGroup

	fmt.Printf("%s[*] Starting TCP scan on host %s%s\n", utils.Blue, scan.Target, utils.Reset)
	fmt.Println(utils.Lines)

	for _, port := range scan.Ports {
		wg.Add(1)
		go tcpWorker(scan.Target, port, scan.Timeout, resultChan, &wg)
	}

	wg.Wait()
	close(resultChan)

	for result := range resultChan {
		results = append(results, result)
	}

	fmt.Println(utils.Lines)
	fmt.Printf("%s[*] TCP Scan finished on host %s%s\n", utils.Blue, scan.Target, utils.Reset)
	fmt.Printf("%s[*] %d ports scanned %d up %s\n", utils.Blue, len(scan.Ports), countOpenPorts(results), utils.Reset)

	return results
}

// Function to perform an UDP Scan
func UdpScan(scan utils.ScanParameters) []utils.Port {
	var results []utils.Port
	resultChan := make(chan utils.Port, len(scan.Ports))
	var wg sync.WaitGroup

	fmt.Printf("%s[*] Starting UDP scan on host %s%s\n", utils.Blue, scan.Target, utils.Reset)
	fmt.Println(utils.Lines)

	for _, port := range scan.Ports {
		wg.Add(1)
		go udpWorker(scan.Target, port, scan.Timeout, resultChan, &wg)
	}

	wg.Wait()
	close(resultChan)

	for result := range resultChan {
		results = append(results, result)
	}

	fmt.Println(utils.Lines)
	fmt.Printf("%s[*] UDP Scan finished on host %s%s\n", utils.Blue, scan.Target, utils.Reset)
	fmt.Printf("%s[*] %d ports scanned %d up %s\n", utils.Blue, len(scan.Ports), countOpenPorts(results), utils.Reset)

	return results
}
