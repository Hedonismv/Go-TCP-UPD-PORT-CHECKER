package main

import (
	"flag"
	"fmt"
	"net"
	"sync"
	"time"
)

const timeout = 1 * time.Second

func scanPort(protocol, hostname string, port int, wg *sync.WaitGroup) {
	defer wg.Done()

	address := fmt.Sprintf("%s:%d", hostname, port)
	var conn net.Conn
	var err error

	if protocol == "tcp" {
		conn, err = net.DialTimeout("tcp", address, timeout)
	} else if protocol == "udp" {
		conn, err = net.DialTimeout("udp", address, timeout)
	}

	if err != nil {
		return //! port is closed
	}
	conn.Close()

	fmt.Printf("[%s] Port %d open to penetration \n", protocol, port)
}

func scanPorts(protocol, hostname string, startPort, endPort int) {
	var wg sync.WaitGroup

	for port := startPort; port <= endPort; port++ {
		wg.Add(1)
		go scanPort(protocol, hostname, port, &wg)
	}

	wg.Wait()
}

func main() {

	hostname := flag.String("host", "localhost", "IP-address or hostname for scanning")
	startPort := flag.Int("start", 1, "Start port for scanning")
	endPort := flag.Int("end", 1024, "End port for scanning")
	protocol := flag.String("proto", "tcp", "Protocol for scanning (tcp or udp)")

	flag.Parse()

	if *startPort <= 0 || *endPort <= 0 || *startPort > *endPort {
		fmt.Println("Wrong port range")
		return
	}
	if *protocol != "tcp" && *protocol != "udp" {
		fmt.Println("Wrong protocol, use tcp or udp")
		return
	}

	fmt.Printf("Scanning %s ports from %d to %d for %s...\n", *protocol, *startPort, *endPort, *hostname)
	scanPorts(*protocol, *hostname, *startPort, *endPort)
	fmt.Println("Scanning completed!")
}
