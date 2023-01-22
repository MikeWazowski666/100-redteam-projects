package main

import (
	"flag"
	"fmt"
	"io/fs"
	"os"
)

func enc(fn string) {
	content, err := os.ReadFile(fn)
	if err != nil {
		panic(err)
	}
	var out []byte
	fmt.Print(content)
	// Simple encryption
	for l, r := range content {
		out = append(out, (r ^ (content[0] * byte(l))))
	}
	// Write to file
	os.WriteFile(fn, out, fs.FileMode(02))
}

func decrypt(fn string) {
	content, err := os.ReadFile(fn)
	if err != nil {
		panic(err)
	}
	var out []byte
	var first byte

	// Simple decryption
	for l, r := range content {
		if l == 0 {
			first = r
		}
		out = append(out, (r ^ (first * byte(l))))
	}

	// Write to file
	os.WriteFile(fn, out, fs.FileMode(02))
}

func main() {
	fn := flag.String("f", "temp.txt", "file to en/decrypt")
	decrypt_FLAG := flag.Bool("de", false, "decrypt flag")
	flag.Parse()
	if !*decrypt_FLAG {
		enc(*fn)
	} else {
		decrypt(*fn)
	}
	fmt.Println("done...")
}
