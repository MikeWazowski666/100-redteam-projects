package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"time"
)

type Client struct {
	name   string
	incMsg chan string
	outMsg chan string
	conn   net.Conn
	reader *bufio.Reader
	writer *bufio.ReadWriter
}

type ChatRoom struct {
	time   time.Time
	client *Client
	text   string
}

func writeToChannel(c chan int, x int) {
	fmt.Println(x)
	c <- x
	close(c)
}

func usage() {
	fmt.Fprintf(os.Stderr, "usage: 1.go [PORT] [TCP/UDP]\nA TCP/UDP commandline chat service.\n")
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
	ch := make(chan int)
	writeToChannel(ch, 12)
	time.Sleep(1 * time.Second)
	writeToChannel(ch, 34)

	// port := args[0]
	// proto := strings.ToLower(args[1])
	//
	// l, err := net.Listen(proto, "127.0.0.1:"+port)
	// if err != nil {
	// 	os.Exit(1)
	// }
	// defer l.Close()

	// for {
	// 	conn, err := l.Accept()
	// 	if err != nil {
	// 		os.Exit(1)
	// 	}

	// 	go func(c net.Conn) {
	// 		buffer := make([]byte, 1024)
	// 		_, err := c.Read(buffer)
	// 		if err != nil {
	// 			panic(err)
	// 		}

	// 		fmt.Print(string(buffer))

	// 		conn.Close()
	// 	}(conn)
	// }
}

func handleRequest(conn net.Conn) {
	buf, read_err := ioutil.ReadAll(conn)
	if read_err != nil {
		fmt.Println("failed:", read_err)
		return
	}
	fmt.Println("Got: ", string(buf))

	_, write_err := conn.Write([]byte("Message received.\n"))
	if write_err != nil {
		fmt.Println("failed:", write_err)
		return
	}
	conn.Close()
}
