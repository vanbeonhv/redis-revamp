package main

import (
	"fmt"
	"log"
	"net"
)

func main(){
	listener, err := net.Listen("tcp",":8080")
	if err != nil {
		log.Fatal("can't open port 8080: %v\n", err)
	}

	defer listener.Close()

	fmt.Println("Server is listening in port 8080")

	for{
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Error accept new connection %v\n", err)
			continue
		}

		fmt.Printf("[+] new Client: %s\n", conn.RemoteAddr().String())
		conn.Close()
	}
}