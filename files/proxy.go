package main

import (
	"bytes"
	"fmt"
	"io"
	"net"
	"os"
)

func main() {
	fmt.Println(os.Args)
	if len(os.Args) != 4 {
		fmt.Println("prog listen_port remote_host, remote_port")
		os.Exit(3)
	}
	listen_port := os.Args[1]
	rmt_host := os.Args[2]
	rmt_port := os.Args[3]
	fmt.Println("localhost:" + listen_port + "  -> " + rmt_host + ":" + rmt_port)
	ln, err := net.Listen("tcp", ":"+listen_port)
	if err != nil {
		panic(err)
	}

	for {
		conn, err := ln.Accept()
		if err != nil {
			panic(err)
		}

		go portFWD(conn, rmt_host, rmt_port)
	}
}

func portFWD(conn net.Conn, rmt_host string, rmt_port string) {

	proxy, err := net.Dial("tcp", rmt_host+":"+rmt_port)
	if err != nil {
		panic(err)
	}

	go func() {
		defer conn.Close()
		defer proxy.Close()

		//experiment part
		var buf bytes.Buffer
		len, err := buf.ReadFrom(proxy)
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Printf("Message size = %d\n", len)
		}
		fmt.Println("total size:", buf.Len())
		_, err = conn.Write(buf.Bytes())

		//io.Copy(conn,proxy) //working part
	}()

	go func() {
		defer conn.Close()
		defer proxy.Close()
		io.Copy(proxy, conn)
	}()
}
