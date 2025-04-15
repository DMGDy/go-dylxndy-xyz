package main
import (
	"net"
	"log"
)

const (
	Port = "8080"
)

func server(client net.Conn) {
	log.Printf("New client: %s\n", client.RemoteAddr())
}


func main() {
	listener, err := net.Listen("tcp", ":"+Port)
	if err != nil {
		log.Printf("Error opening listen socket: %s\n", err.Error())
	}
	log.Printf("Listening to connections on port: %s\n", Port);



	for {
		client, err := listener.Accept()
		if err != nil {
			log.Printf("Error accepting connection: %s\n",err.Error())
			continue
		}

		go server(client)

	}


}
