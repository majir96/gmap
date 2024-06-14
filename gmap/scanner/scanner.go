package scanner

import (
	"fmt"
	"gmap/utils"
	"net"
	"sync"
	"time"

	"github.com/go-ping/ping"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
)

// Auxiliary function to get local IP
func getLocalIp() (*net.IP, error) {
	// Ping google to check the local ip
	conn, err := net.Dial("udp", "8.8.8.8:80")

	if err != nil {
		return nil, err
	}

	// Ensure connection is closed
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	// Return IP as a string
	return &localAddr.IP, nil

}

// Auxiliary function to get interface based on IP
func getInterface(ip *net.IP) (*net.Interface, error) {
	// Get all interfaces
	interfaces, err := net.Interfaces()

	if err != nil {
		return nil, err
	}

	// For each interface
	for _, iface := range interfaces {
		addrs, err := iface.Addrs()

		if err != nil {
			continue
		}

		// For each IP in an interface
		for _, addr := range addrs {
			var ifaceIP net.IP

			switch v := addr.(type) {
			case *net.IPNet:
				ifaceIP = v.IP
			case *net.IPAddr:
				ifaceIP = v.IP
			}

			// If it is the same return the interface
			if ifaceIP.Equal(*ip) {
				return &iface, nil
			}
		}
	}

	return nil, fmt.Errorf("no interface found for IP %s", ip)
}

// Auxiliary function to send RST packet and close connection
func sendRST(srcIP, dstIP net.IP, srcPort, dstPort layers.TCPPort, iface *net.Interface) error {
	ethLayer := &layers.Ethernet{
		SrcMAC:       iface.HardwareAddr,
		EthernetType: layers.EthernetTypeIPv4,
	}

	ipLayer := &layers.IPv4{
		SrcIP:    srcIP,
		DstIP:    dstIP,
		Protocol: layers.IPProtocolTCP,
	}

	tcpLayer := &layers.TCP{
		SrcPort: srcPort,
		DstPort: dstPort,
		RST:     true,
		Window:  14600,
	}
	tcpLayer.SetNetworkLayerForChecksum(ipLayer)

	buffer := gopacket.NewSerializeBuffer()
	options := gopacket.SerializeOptions{FixLengths: true, ComputeChecksums: true}
	gopacket.SerializeLayers(buffer, options, ethLayer, ipLayer, tcpLayer)
	outgoingPacket := buffer.Bytes()

	handle, err := pcap.OpenLive(iface.Name, 65536, true, pcap.BlockForever)
	if err != nil {
		return err
	}
	defer handle.Close()

	return handle.WritePacketData(outgoingPacket)
}

// Syn worker for TCP Syn Scan
func synWorker(target string, port int, timeout time.Duration, results chan<- utils.Port, wg *sync.WaitGroup) {
	// Ensure worker is done
	defer wg.Done()

	// Get the local IP
	srcIp, err := getLocalIp()
	if err != nil {
		utils.PrintError("[ERROR] caused getting local IP")
	}

	iface, err := getInterface(srcIp)
	if err != nil {
		utils.PrintError(err.Error())
	}

	// Interface in which we will receive the traffic
	handle, err := pcap.OpenLive(iface.Name, 65536, true, pcap.BlockForever)
	if err != nil {
		utils.PrintError(err.Error())
	}
	// Ensure connection is being closed
	defer handle.Close()

	// Set packet parameters
	dstIp := net.ParseIP(target)
	dstPort := layers.TCPPort(port)
	//! Modifiable src port if necessary -> current 12345
	srcPort := layers.TCPPort(12345)

	// Define package layers
	ethLayer := &layers.Ethernet{
		SrcMAC:       iface.HardwareAddr,
		EthernetType: layers.EthernetTypeIPv4,
	}

	ipLayer := &layers.IPv4{
		SrcIP:    *srcIp,
		DstIP:    dstIp,
		Protocol: layers.IPProtocolTCP,
	}

	tcpLayer := &layers.TCP{
		SrcPort: srcPort,
		DstPort: dstPort,
		SYN:     true,
		Seq:     1105024978,
		Window:  14600,
	}

	tcpLayer.SetNetworkLayerForChecksum(ipLayer)

	buffer := gopacket.NewSerializeBuffer()
	options := gopacket.SerializeOptions{FixLengths: true, ComputeChecksums: true}
	// Assemble the package
	gopacket.SerializeLayers(buffer, options, ethLayer, ipLayer, tcpLayer)
	outgoingPacket := buffer.Bytes()

	err = handle.WritePacketData(outgoingPacket)

	if err != nil {
		utils.PrintError(err.Error())
	}

	// Set Berkeley Packet Filter to obtain TCP Packets of the target on the desired port
	handle.SetBPFFilter(fmt.Sprintf("tcp and src host %s and src port %d and dst port %d", target, port, srcPort))
	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
	timeoutChan := time.After(timeout)

	var service, state string

	for {
		select {
		case packet := <-packetSource.Packets():
			// We check the tcp layer
			if tcpLayer := packet.Layer(layers.LayerTypeTCP); tcpLayer != nil {
				tcp, _ := tcpLayer.(*layers.TCP)

				// Opened port (SYN + ACK received)
				if tcp.SYN && tcp.ACK {
					err = sendRST(*srcIp, dstIp, srcPort, dstPort, iface)

					if err != nil {
						utils.PrintError(fmt.Sprintf("[ERROR] Fail to close connection on port %d", port))
					}

					state = "open"

					// Attempt banner grabbing to determine service
					address := fmt.Sprintf("%s:%d", target, port)
					conn, err := net.DialTimeout("tcp", address, timeout)
					if err == nil {
						service = bannerGrab(conn)
						conn.Close()
					}

				} else if tcp.RST {
					state = "closed"
					service = ""
				}
			}

		// Case in which we run into a timeout -> filtered
		case <-timeoutChan:
			state = "filtered"
			service = ""
		}
	}
	service = checkService(service)
	results <- utils.Port{Port: port, Status: state, Service: service}

}

func SynScan(scan utils.ScanParameters) []utils.Port {
	var results []utils.Port
	resultChan := make(chan utils.Port, len(scan.Ports))
	var wg sync.WaitGroup

	fmt.Printf("%s[*] Starting SYN scan on host %s%s\n", utils.Blue, scan.Target, utils.Reset)
	fmt.Println(utils.Lines)

	for _, port := range scan.Ports {
		wg.Add(1)
		go synWorker(scan.Target, port, scan.Timeout, resultChan, &wg)
	}

	wg.Wait()
	close(resultChan)

	for result := range resultChan {
		results = append(results, result)
	}

	fmt.Println(utils.Lines)
	fmt.Printf("%s[*] SYN Scan finished on host %s%s\n", utils.Blue, scan.Target, utils.Reset)
	fmt.Printf("%s[*] %d ports scanned %d up %s\n", utils.Blue, len(scan.Ports), countOpenPorts(results), utils.Reset)

	return results
}

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
