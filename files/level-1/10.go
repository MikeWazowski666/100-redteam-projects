package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"net"
	"os"
)

var conn_host string
var buf bytes.Buffer
var proto string

func listen(proto string, host string, port string) {
	l, err := net.Listen(proto /*change this ->*/, host+":"+port)
	if err != nil {
		panic(err)
	}
	go func() {
		conn, err := l.Accept()
		defer conn.Close()
		if err != nil {
			panic(err)
		}
		var data []byte
		_, err = conn.Read(data)
		if err != nil {
			panic(err)
		}
		fmt.Println(data)
	}()

}

func connect(proto string, host string, port string) {
	c, err := net.Dial(proto, host+":"+port)
	if err != nil {
		c.Close()
		panic(err)
	}
	defer c.Close()
	for {
		go func() {
			var v []byte
			len, err := c.Read(v)
			if err != nil {
				panic(err)
			}
			if len > 0 {
				fmt.Println(v)
			}
		}()
		go func() {
			input, err := bufio.NewReader(os.Stdin).ReadString('\n')
			if err != nil {
				fmt.Print(err)
			}
			c.Write([]byte(input))
		}()
	}

	// go func() {
	// 	if err != nil {
	// 		c.Close()
	// 		panic(err)
	// 	}
	// 	data, err := in_scanner.ReadByte()
	// 	if err != nil {
	// 		c.Close()
	// 		panic(err)
	// 	}
	// 	_, err := c.Write(data)
	// 	if err != nil {
	// 		c.Close()
	// 		panic(err)
	// 	}
	// 	fmt.Println(buf)
	// }()
}

func main() {
	listen_bool := flag.Bool("l", false, "listen")
	sel_proto := flag.Bool("u", false, "select udp protocol")
	flag.Parse()

	port := &os.Args[2]

	if *sel_proto == false {
		proto = "tcp"
	} else {
		proto = "udp"
	}
	if *listen_bool != false {
		conn_host = os.Args[1]
	}
	conn_host := &os.Args[1]

	fmt.Println(*conn_host)
	if *listen_bool == true {
		listen(proto, *conn_host, *port)
	} else {
		connect(proto, *conn_host, *port)
	}
}
