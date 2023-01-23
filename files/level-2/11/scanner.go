package main

import (
	"flag"
	"fmt"
	"net"
	"strconv"
	"strings"
)

func TCP_Scan(target string, port_range int) []int {
	var out []int
	var i int = 1
	for i <= port_range {
		current_try := target + ":" + strconv.Itoa(i)
		conn, err := net.Dial("tcp", current_try)
		if err == nil {
			fmt.Println("Found:", conn.RemoteAddr().String()+"/tcp")
			out = append(out, i)
			conn.Close()
		}
		i++
	}
	return out
}

func UDP_Scan(target string, port_range int) []int {
	var out []int
	var i int = 1
	for i <= port_range {
		conn, err := net.DialUDP("udp", nil, &net.UDPAddr{
			IP: net.ParseIP(target), Port: i})
		if err == nil {
			_, err := conn.Write([]byte{181})
			if err == nil {
				fmt.Println("Found:", conn.RemoteAddr().String()+"/udp")
				out = append(out, i)
				conn.Close()
			}
		}
		i++
	}
	return out
}

func main() {
	target := flag.String("host", "127.0.0.1", "host to scan")
	port_range := flag.Int("p", 1024, "ports to scan")
	udp_FLAG := flag.Bool("udp", false, "udp scan flag (does not work 100%)")
	flag.Parse()
	var out []int // Array of open ports
	// TODO: Check if machine is online

	if *udp_FLAG {
		out = UDP_Scan(*target, *port_range)
	} else {
		out = TCP_Scan(*target, *port_range)
	}

	fmt.Printf("\nScan of %s complete:\n", *target)
	for _, l := range out {
		fmt.Printf("Port %d is open\n", l)
	}
	fmt.Println(strings.Repeat("-", 20))
}
