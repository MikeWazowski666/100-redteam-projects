package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"strings"
)

func usage() {
	fmt.Fprintf(os.Stderr, "usage: 1.go [PORT] [TCP/UDP]\nJust a TCP/UDP port listner.\n")
	flag.PrintDefaults()
	os.Exit(2)
}

func main() {
	flag.Usage = usage
	flag.Parse()

	args := flag.Args()
	if len(args) < 2 {
		fmt.Println("Some arguments are missing.")
		os.Exit(1)
	}
	port := args[0]
	proto := strings.ToLower(args[1])

	l, err := net.Listen(proto, ":"+port)
	if err != nil {
		os.Exit(1)
	}
	defer l.Close()

	for {
		conn, err := l.Accept()
		if err != nil {
			os.Exit(1)
		}

		go func(c net.Conn) {
			buffer := make([]byte, 1024)
			_, err := c.Read(buffer)
			if err != nil {
				panic(err)
			}

			fmt.Print(string(buffer))

			conn.Close()
		}(conn)

		// go func(c net.Conn) {
		// 	// Echo all incoming data.
		// 	io.Copy(c, c)
		// 	fmt.Println(content)
		// 	// Shut down the connection.
		// 	c.Close()
		// }(conn)
	}
}
