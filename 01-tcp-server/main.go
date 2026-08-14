package main

import (
	"fmt"
	"io"
	"log"
	"net"
)

func main() {
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalf("can't open port 8080: %v\n", err)
	}

	defer listener.Close()

	fmt.Println("Server is listening in port 8080")

	for {
		connection, err := listener.Accept()
		if err != nil {
			log.Printf("Error accept new connection %v\n", err)
			continue
		}
		
		go handleConnection(connection)

		fmt.Printf("[+] new Client: %s\n", connection.RemoteAddr().String())
	}
}


func handleConnection(conn net.Conn){

	defer conn.Close()

	buffer := make([]byte, 1024)

	for {
		
		// lấy ra tổng đọc số byte thực tế có data thôi. 
		noOfByteActualRead, err := conn.Read(buffer)
		
		if(err == io.EOF) {
 			log.Printf("Client disconnected %v\n", err)
			break
		}
		
		if (err != nil){
			log.Printf("Network error %v\n", err)
			break
		}
		
		// Cắt ra data bằng slice, cái đống element còn lại toàn null hay gì đó, đại khái data vô nghĩa.
		dataReceived := buffer[:noOfByteActualRead]
		// dùng %s để converse sang string, chứ để %v nó in moẹ ra số
		log.Printf("dataReceived %s\n", dataReceived)
	
	}
}
