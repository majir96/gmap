package scanner

import (
	"fmt"
	"gmap/utils"
	"net"
	"time"
)

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

// Function to perform a basic TCP Scan
func TcpScan(scan utils.ScanParameters) []utils.Port {
	var results []utils.Port

	fmt.Printf("%s[*] Starting TCP scan on host %s%s\n", utils.Blue, scan.Target, utils.Reset)
	fmt.Println(utils.Lines)

	for _, port := range scan.Ports {
		var state, service string

		address := fmt.Sprintf("%s:%d", scan.Target, port)
		conn, err := net.DialTimeout("udp", address, scan.Timeout)

		// Closed Port
		if err != nil {
			// Check if it is really closed
			if opErr, ok := err.(*net.OpError); ok && opErr.Err.Error() == "connect: connection refused" {
				state = "closed"
			} else { // We asumme it is filtered
				state = "filtered"
			}
			service = utils.CommonServices[port]
		} else {
			service = bannerGrab(conn)
			state = "open"
		}

		// Making sure the connection will be closed afterwards
		defer conn.Close()

		service = checkService(service)
		results = append(results, utils.Port{port, state, service})
	}

	fmt.Println(utils.Lines)
	fmt.Printf("%s[*] TCP Scan finished on host %s%s\n", utils.Blue, scan.Target, utils.Reset)
	fmt.Printf("%s[*] %d ports scanned %d up %s%s\n", utils.Blue, len(scan.Ports), countOpenPorts(results), utils.Reset)

	return results
}

// Function to perform an UDP Scan
func UdpScan(scan utils.ScanParameters) []utils.Port {
	var results []utils.Port

	fmt.Printf("%s[*] Starting UDP scan on host %s%s\n", utils.Blue, scan.Target, utils.Reset)
	fmt.Println(utils.Lines)

	for _, port := range scan.Ports {
		var state, service string

		address := fmt.Sprintf("%s:%d", scan.Target, port)
		conn, err := net.DialTimeout("udp", address, scan.Timeout)

		// Closed port
		if err != nil {
			service = utils.CommonServices[port]
			results = append(results, utils.Port{port, "closed", service})
			continue
		}

		// Making sure the connection will be closed afterwards
		defer conn.Close()

		// Sending a ping through connection
		_, err = conn.Write([]byte("Ping"))

		// Consider error in ping to be a closed port
		if err != nil {
			results = append(results, utils.Port{port, "closed", ""})
			continue
		}

		// Set deadline to wait for response based on set timeout
		conn.SetReadDeadline(time.Now().Add(scan.Timeout))

		// Read response
		buff := make([]byte, 1024)
		n, err := conn.Read(buff)

		if err != nil { // If there is no response it is either filtered or opened

			service = utils.CommonServices[port]
			state = "open/filtered"

		} else { // If there is response it is opened
			service = string(buff[:n])
			state = "open"
		}

		service = checkService(service)
		results = append(results, utils.Port{port, state, service})

	}

	fmt.Println(utils.Lines)
	fmt.Printf("%s[*] UDP Scan finished on host %s%s\n", utils.Blue, scan.Target, utils.Reset)
	fmt.Printf("%s[*] %d ports scanned %d up %s%s\n", utils.Blue, len(scan.Ports), countOpenPorts(results), utils.Reset)

	return results
}
