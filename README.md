# gmap

A basic project to learn **Go** which I will be updating through time. The tool is designed to perform port scanning on hosts by leveraging Go's quick run times as well as its portabilities. 

---
---

# Installation   

To install gmap, ensure you have [Go](https://golang.org/dl/) installed on your machine, then follow these steps:

### Clone the repository: 
```sh 
git clone https://github.com/majir96/gmap.git
cd gmap 
```



## Compile using Go 

To compile gmap you can use the Go build tool. Here are the steps to do so, including options to reduce the size of the executable:

### **Build the Executable**

```sh
go build -ldflags="-s -w" -o gmap main.go
```
#### Options for other OS
You can cross-compile the executable for different operating systems. 

##### Windows 

```sh
GOOS=windows GOARCH=amd64 go build -ldflags="-s -w" -o gmap.exe main.go
```

----

##### macOS 

```sh 
GOOS=darwin GOARCH=amd64 go build -ldflags="-s -w" -o gmap main.go
```

### Reduce the size of the executable 

If you have **[upx](https://github.com/upx/upx)** installed you can use it to reduce the binary size to make it lighter. 

1. **Install upx**
```
sudo apt install upx 
```

2.**Reduce binary size**
```sh
upx gmap 
```

## Dockerfile 

To build and run gmap using Docker:

1. **Build the Docker Image**

```sh
docker build -t gmap .
```

2. **Run the Docker Container**

```sh 
docker run --rm -it gmap -t 127.0.0.1 -p 80,443 -o output.txt -f txt
```


---

# Features

- Scan specific ports or ranges of ports
- Scan all ports (0-65535)
- Perform TCP and UDP scans
- Export scan results to text, CSV, or JSON files
- Filter results to show only open ports
- Set custom timeout for scan operations

# Usage 

```go
go run main.go -t <IP> -p <PORTS> -s <scan> -o <OUTPUT FILE> -f <FORMAT>
```

## Options
- **-t, --target \<IP>**: Target to scan (required)
- **-p, --port \<PORTS>**: Port(s) to scan. Default set to common ports. Separate multiple ports with commas (e.g., -p 22,80,443) or specify a range (e.g., -p 0-100).
- **-p-**: Scan all ports (0-65535)
- **-s, --scan \<SCAN>**: Type of scan to perform. Options:
    - **tcp**: Perform a TCP scan (default)
    - **udp**: Perform a UDP scan
- **-h, --help**: Display the help message
- **-o, --output \<FILE>**: Export output to a file (default format: .txt)
- **-f, --format \<FORMAT>**: Format to export the file to. Formats:
    - **txt**: Export to text file (default)
    - **csv**: Export to CSV file
    - **json**: Export to JSON file
- **--open**: Filter by open ports on output
- **--timeout \<TIMEOUT>**: Timeout for packets when scanning (e.g., 500ms, 2s, 1m)

### Examples

#### Scan a single IP for common ports

```sh
go run main.go -t 127.0.0.1
```

---

#### Scan a range of ports on a target IP

```sh
go run main.go -t 192.168.1.1 -p 20-80
```

---

#### Perform a UDP scan and export the results to a JSON file

```sh
go run main.go -t 10.0.0.1 -p 53,123,161 -s udp -o scan_results -f json
 
```


# Future Implementations 
- Use of Docker to deploy the tool 
- Including nmap support 
- Including new types of scan 
- Including verbose mode 
- Including proxy/IP spoofing 
- Including vulnerability scanning 

# Contributing
Contributions are welcome! Please fork the repository and create a pull request to add features, improve documentation, or fix bugs.

# License
This project is licensed under the **GPL License**. See the [LICENSE](LICENSE) file for details.

# Ethical Considerations and Responsability 

**Using gmap to scan networks without proper authorization is illegal and unethical**. Always obtain explicit permission from the owner of the network before performing any scans. Unauthorized scanning can lead to severe legal consequences and negatively impact the security and stability of networks. **The developers of gmap are not responsible for any misuse of this tool. Use it responsibly and ethically.**