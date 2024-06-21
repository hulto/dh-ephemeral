package main

import (
	"bufio"
	"crypto/ecdh"
	"crypto/rand"
	"fmt"
	"io"
	"net"
	"os"
)

func crypto() {
	x22519 := ecdh.X25519()
	priv_key, err := x22519.GenerateKey(rand.Reader)
	if err != nil {
		fmt.Println("Error: ", err)
	}
	fmt.Println("priv_key ", priv_key)
	public_key, err := x22519.NewPublicKey(priv_key.Bytes())
	if err != nil {
		fmt.Println("Error: ", err)
	}
	fmt.Println("public_key ", public_key)

	// client_priv_key, err := x22519.GenerateKey(rand.Reader)
	// if err != nil {
	// 	fmt.Println("Error: ", err)
	// }
	// fmt.Println("client_priv_key ", client_priv_key)
	// client_public_key, err := x22519.NewPublicKey(priv_key.Bytes())
	// if err != nil {
	// 	fmt.Println("Error: ", err)
	// }
	// fmt.Println("client_public_key ", client_public_key)

	// shared_key, err := priv_key.ECDH(client_public_key)
	// if err != nil {
	// 	fmt.Println("Error: ", err)
	// }

	// fmt.Println("shared_key ", shared_key)
}

func generate_key_pairs() (*ecdh.PublicKey, *ecdh.PrivateKey) {
	x22519 := ecdh.X25519()
	// priv_key, err := x22519.GenerateKey(rand.Reader)
	priv_key_bytes := []byte{103, 189, 40, 137, 41, 70, 3, 210, 150, 37, 15, 119, 12, 213, 209, 38, 202, 101, 9, 88, 26, 61, 171, 10, 143, 207, 15, 244, 121, 243, 234, 122}
	priv_key, err := x22519.NewPrivateKey(priv_key_bytes)
	if err != nil {
		fmt.Println("Error: ", err)
	}
	fmt.Println("priv_key ", priv_key)

	// ---

	public_key, err := x22519.NewPublicKey(priv_key.PublicKey().Bytes())
	if err != nil {
		fmt.Println("Error: ", err)
	}
	fmt.Println("public_key ", public_key)

	return public_key, priv_key
}

func main() {
	generate_key_pairs()
	// obtain the port and prefix via program arguments
	port := "127.0.0.1:9090"
	prefix := ""

	// create a tcp listener on the given port
	listener, err := net.Listen("tcp", port)
	if err != nil {
		fmt.Println("failed to create listener, err:", err)
		os.Exit(1)
	}
	fmt.Printf("listening on %s, prefix: %s\n", listener.Addr(), prefix)

	// listen for new connections
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("failed to accept connection, err:", err)
			continue
		}

		// pass an accepted connection to a handler goroutine
		go handleConnection(conn, prefix)
	}

}

// handleConnection handles the lifetime of a connection
func handleConnection(conn net.Conn, prefix string) {
	defer conn.Close()
	reader := bufio.NewReader(conn)
	for {
		// read client request data
		bytes := make([]byte, 1024)
		_, err := io.ReadFull(reader, bytes)
		if err != nil {
			panic(err)
		}
		fmt.Printf("request: %s", bytes)

		// prepend prefix and send as response
		line := fmt.Sprintf("%s%s", prefix, bytes)
		fmt.Printf("response: %s", line)
		conn.Write([]byte(line))
	}
}
