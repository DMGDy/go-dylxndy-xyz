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

type HeaderField int
const (
	RequestLine = iota + 1
	Host
	UserAgent
	Accept
	Connection
	Referer
)


func parseHeader(header string) {
	lines := strings.Split(header, "\n")
	for line_n , line := range lines {
		if line_n == 0 {
			log.Printf("%s\n", line)
		}
		parts := strings.SplitN(line, ": ", 2)
		for part_n, part := range parts {
			if part_n == 1 {
				vals := strings.Split(part, ",")
				for _, val := range vals {
					log.Printf("\t%s\n", val)
				}
			}

		}
	}
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

	parseHeader(header)

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
