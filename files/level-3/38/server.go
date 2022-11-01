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
	"strconv"
	"time"
)

type User struct {
	ip   string
	date int64
	key  string
}

// Generates a key
func generateKey() (string, error) {
	bytes := make([]byte, 64)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

func writeToFile(saveFile string, data User) {
	ip := data.ip
	date := strconv.Itoa(int(data.date)) // converting int64 to str
	key := data.key

	fmt.Println(date, ip, key) // log to stdout

	f, err := os.OpenFile(saveFile, os.O_APPEND, 0600)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer f.Close()
	f.Write([]byte(date + " " + ip + " " + key + "\n"))
	if err != nil {
		fmt.Println(err)
		return
	}
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

func readFile(saveFile string, ip string) { //(bool, line int)
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

	// scanner := bufio.NewScanner(file)
	// print(scanner.Err())
	// // Loop over all lines in the file and print them.
	// for scanner.Scan() {
	// 	line := scanner.Text()
	// 	fmt.Println(line)
	// }

	// var line int = 0

	// return true, line

}

func main() {
	port := flag.String("p", "80", "port to run on")
	saveFile := flag.String("f", "./data", "savefile")

	flag.Parse()

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

		readFile(*saveFile, strRemoteAddr) // find if ip exists
		date := time.Now().Unix()
		// _ = date
		fmt.Println(date, strRemoteAddr, key)

		// Write key
		go func(c net.Conn) {
			conn.Write([]byte(string(key)))
			// Shut down the connection.
			c.Close()
		}(conn)
	}
}
