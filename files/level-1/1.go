package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"strings"
	"time"
)

var delimiter []byte = []byte(strings.Repeat("-", 50) + "\n")

type User struct {
	ip       string
	username string
}

func usage() {
	fmt.Fprintf(os.Stderr, "usage: 1.go [PORT] [TCP/UDP]\nA TCP/UDP CLI chat service.\n")
	flag.PrintDefaults()
	os.Exit(2)
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

/*
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
*/

func authenticate(conn net.Conn) User {
	username := make([]byte, 512)
	var len int
	_, err := conn.Write([]byte("Username: "))
	for len == 0 {
		len, err = conn.Read(username)
		check(err)
	}
	user := User{conn.RemoteAddr().String(), strings.Replace(string(username), "\n", "", 1)}
	return user
}

func search_user(ip string, users []User) User {
	for i := range users {
		if users[i].ip == ip {
			return users[i]
		}
	}
	return User{"127.0.0.1", "anon"}
}

func listen(proto string, port string) {
	l, err := net.Listen(proto, "127.0.0.1:"+port)
	check(err)
	defer l.Close()
	var users []User
	for {
		conn, err := l.Accept()
		check(err)
		users = append(users, authenticate(conn))
		go func(c net.Conn) {
		LOOP:
			buffer := make([]byte, 1024)
			recv_ip := c.RemoteAddr().String()
			len, _ := c.Read(buffer)
			check(err)
			if len > 0 { // check is someone actually sent data
				user := search_user(recv_ip, users)
				time := time.Now().UTC().String()
				message := []byte("[" + time + "]: " + user.username + "\n")
				// send data
				c.Write(message)
				c.Write(buffer)
				c.Write(delimiter)
				fmt.Printf("[%s] %s: %s ", time, user.username, string(buffer))
			}
			goto LOOP
		}(conn)
	}
}

func main() {
	flag.Usage = usage
	flag.Parse()

	args := flag.Args()
	if len(args) < 2 {
		fmt.Println("Some arguments are missing.\nusage: prog [port] [tcp/udp]")
		os.Exit(1)
	}
	port := args[0]
	proto := strings.ToLower(args[1])
	listen(proto, port)
}
