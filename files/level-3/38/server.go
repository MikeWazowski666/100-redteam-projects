package main

// Run: go run server.go [args] | tee <datafile>

import (
	"bufio"
	"crypto/rand"
	"encoding/hex"
	"flag"
	"fmt"
	"net"
	"os"
	"time"
)

type User struct {
	ip   string
	date string
	key  string
}

// Generates a key
func generateKey() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

func writeToFile(saveFile string, data User) {
	ip := data.ip
	date := data.date
	key := data.key

	f, err := os.OpenFile(saveFile, os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		fmt.Println(err)
		return
	}

	_, err = f.Write([]byte(date + " " + ip + " " + key + "\n"))
	if err != nil {
		fmt.Println(err)
		return
	}
	f.Close()
}

func Readln(r *bufio.Reader) (string, error) {
	var (
		isPrefix bool  = true
		err      error = nil
		line, ln []byte
	)

	for isPrefix && err == nil {
		line, isPrefix, err = r.ReadLine()
		ln = append(ln, line...)
	}
	print(ln)
	return string(ln), err
}

func readFile(saveFile string, ip string) {
	file, err := os.Open(saveFile) // Open file
	if err != nil {
		fmt.Printf("no such file: %s", file)
	}
	defer file.Close()

	r := bufio.NewReader(file)
	s, e := Readln(r)
	print(s)

	for e == nil {
		fmt.Println(s)
		s, e = Readln(r)
	}
}

func main() {
	port := flag.String("p", "443", "port to run on")
	saveFile := flag.String("f", "data", "savefile")
	flag.Parse()

	_, chkFile := os.Open(*saveFile)
	if chkFile != nil {
		os.Create(*saveFile)
	}

	// Start tcp server
	l, err := net.Listen("tcp", "127.0.0.1:"+*port)
	if err != nil {
		os.Exit(1)
	}

	defer l.Close()

	for {
		conn, err := l.Accept()
		if err != nil {
			os.Exit(1)
		}
		strRemoteAddr := conn.RemoteAddr().String()

		key, err := generateKey()
		if err != nil {
			fmt.Printf("error generating key: %s", err)
		}

		// readFile(*saveFile, strRemoteAddr) // find if ip exists
		date := time.Now().String()
		// date to readable format
		fmt.Println(date, strRemoteAddr, key)

		// Write key
		go func(c net.Conn) {
			conn.Write([]byte(string(key)))
			writeToFile(*saveFile, User{strRemoteAddr, date, key})
			// Shut down the connection.
			c.Close()
		}(conn)
	}

}
