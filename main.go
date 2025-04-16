package main
import (
	"net"
	"log"
	"crypto/tls"
	"strings"
)

const (
	PORT = "42069"
	MAX_HEADER_LEN = 1024*4
)

func parseHeader(header string) {
	
}

func isCompleteHeader(header string) bool {
	if strings.Contains(header, "\r\n\r\n") || strings.Contains(header, "\n\n") {
		return true
	}
	return false
}

func server(client net.Conn) {
	log.Printf("New client: %s\n", client.RemoteAddr())
	buffer := make([]byte, MAX_HEADER_LEN)

	n, err := client.Read(buffer)
	if err != nil {
		log.Printf("Error reading from client: %s\n", err.Error())
		client.Close()
		return
	}
	if n == MAX_HEADER_LEN {
		log.Printf("Malformed Header: Too Large!\n")
		client.Close()
		return
	}

	header := string(buffer)

	if !isCompleteHeader(header) {
		log.Printf("Malformed Header: Not properly terminated")
		client.Close()
		return
	}

	log.Printf("%s\n", header)
	response := "Hello, World!"
	n, err = client.Write([]byte(response))
	
	if err != nil {
		log.Printf("Error writing back to client: %s\n",err.Error())
		client.Close()
		return
	}


	client.Close()
}


func main() {

	cert, err := tls.LoadX509KeyPair("cert/cert.pem", "cert/key.pem")

	if err != nil {
		log.Printf("Error loading certificates: %s\n", err.Error())
	}

	config := &tls.Config {
		Certificates: []tls.Certificate{cert},
	}

	listener, err := tls.Listen("tcp", ":"+PORT, config)

	if err != nil {
		log.Printf("Error opening listen socket: %s\n", err.Error())
	}

	log.Printf("Listening to connections on port: %s\n", PORT);

	for {
		client, err := listener.Accept()
		if err != nil {
			log.Printf("Error accepting connection: %s\n",err.Error())
			continue
		}

		go server(client)

	}


}
