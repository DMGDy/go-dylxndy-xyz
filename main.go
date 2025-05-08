/* go-server
 * Copyright (C) 2025 Dylan Dy <dylangarza1909@gmail.com>
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Affero General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Affero General Public License for more details.
 *
 * You should have received a copy of the GNU Affero General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
*/

package main
import (
	"errors"
	"fmt"
	"net"
	"log"
	"crypto/tls"
	"strings"
	"os"
)

const (
	PORT = "42069"
	MAX_HEADER_LEN = 1024*4
)

type RequestHeader struct {
	method string
	path int
	version string
	ua string
	accept string
	ref string
	keep_alive bool
}

type RequestLine  int
const (
	Method = iota
	File
	HTTPVersion
)

type RequestFields int
const (
	Skip = iota
	UserAgent
	Accept
	Referer
	KeepAlive
)

var files = [...]string{
	"error.html",
	"",
	"index.html" ,
	"styles.css",
	"favicon.ico",
	"assets/android-chrome-192x192.png" ,
	"assets/android-chrome-512x512.png",
	"assets/apple-touch-icon.png",
	"assets/favicon-16x16.png",
	"assets/favicon-32x32.png",
	"assets/favicon.ico",
	"assets/trollface-drift-phonk.gif",
	"assets/buttons/agplv3.png",
	"assets/buttons/archlinux.gif",
	"assets/buttons/linux_powered.gif",
	"assets/buttons/vim.gif",
	"assets/buttons/wget.gif",
}

// types mapped to the same index as file

var mime_types = [...]string{
	"text/html",                                                                              
	"text/html",
	"text/html",
	"text/css",
	"image/x-icon",
	"image/png",
	"image/png",
	"image/png",
	"image/png",
	"image/png",
	"image/x-icon",
	"image/gif",
	"image/png",
	"image/gif",
	"image/gif",
	"image/gif",
	"image/gif",
}

var debug = false
var log_file_path = ""

func logStdout(message string) {
	if debug {
		log.Println(message)
	}
}

func logFile(messages <- chan string) {
	file, err := os.OpenFile(log_file_path, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)

	if err != nil {
		logStdout(fmt.Sprintf("Error opening %s: %s", log_file_path, err.Error()))
		return
	}

	defer file.Close()

	for message := range messages {
		_, err := file.WriteString(message + "\n")
		if err != nil {
			logStdout(fmt.Sprintf("Error writing to %s: %s", log_file_path, err.Error))
		}
	}
	
}

func sendResponse(client net.Conn, rq_info RequestHeader) (int, error) {
	/* open file
	 * construct response
	 * send response
	 * send file
	*/

	file_bytes, err := os.ReadFile("dylxndy.xyz/" + files[rq_info.path])
	if err != nil {
			return 1, err
	}

	content_length := len(file_bytes)
	
	response := fmt.Sprintf("HTTP/1.1 200 OK\n"+
							"Server: epic-server v420.69 (Linux)\n"+
							"Accept-Ranges: Bytes\n"+
							"Connection: Keep-Alive\n"+
							"Content-Type: %s\n"+
							"Content-Length: %d"+
							"\r\n\r\n", mime_types[rq_info.path], content_length)
	_, err = client.Write([]byte(response))

	written, err := client.Write(file_bytes)

	return written, nil
}

func getFile(requested string) int {
	for n, file := range files {
		if debug {
			log.Printf("%s ?= %s\n", requested, file)
		}
		if requested[1:] == file {
			if n == 1 { 
				n = n + 1
			}
			return n
		}
	}
	return 0
}

func parseHeader(header string) (*RequestHeader, error) {
	lines := strings.Split(header, "\n")

	// if less than 2 lines or more than 32, something is not right
	if len(lines) < 2 || len(lines) > 32 {
		return nil, errors.New("Malformed or incorrect Header\n") 
	}

	rq := RequestHeader{}
	file_index := 0

	for line_n, line := range lines {
		// first line is request line
		if line_n == 0 {
			rql_values := strings.Split(line, " ")

			for n, value := range rql_values {
				switch n {
				case Method:
					if value != "GET" {
						return nil, errors.New("Not Supported Request Method\n")
					}
					rq.method = value
				case File:
					file_index = getFile(value)
					rq.path = file_index
					if debug{
						fmt.Printf("Requested Path: %s\n", value)
					}
				case HTTPVersion:
					rq.version = value
					if debug {
						fmt.Printf("HTTP Version: %s\n", value)
					}
				}
			}
			continue
		}
		fields := strings.SplitN(line, ": ", 2)
		header_field := 0
		for field_n, field := range fields {
			
			// first part tells what im looking at
			if field_n == 0 {
				switch field {
				case "User-Agent":
					header_field = UserAgent
				case "Accept":
					header_field = Accept
				case "Connection":
					header_field = KeepAlive
				default:
					header_field = Skip
				}
				continue
			}
			switch header_field {
			case Skip:
				continue
			case UserAgent: 
				rq.ua = field
			case Accept:
				rq.accept = mime_types[rq.path]
			case KeepAlive:
				if strings.ToLower(field) == "keep-alive" {
					rq.keep_alive = true
				} else {
					rq.keep_alive = false
				}
			}
		}
	}
	return &rq, nil
}

func isCompleteHeader(header string) bool {
	// who knows what a web browser sends today
	if strings.Contains(header, "\r\n\r\n") || strings.Contains(header, "\n\n") {
		return true
	}
	return false
}

func Server(client net.Conn) {
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

	rq_info, err := parseHeader(header)

	if err != nil  || rq_info == nil {
		log.Printf("Error Parsing request: %s\n",err.Error())
		client.Close()
		return
	}

	bytes, err := sendResponse(client, *rq_info)
	
	if err != nil {
		log.Printf("Error writing back to client: %s\n",err.Error())
		client.Close()
		return
	}

	log.Printf("Wrote back %d bytes", bytes)

	client.Close()
}

func parseCliArgs(argv []string) {

	if len(argv) > 1 {
		for argc, arg := range argv {
			if argc == 0 {
				continue
			}

			// all cli arguments have '=' 
			arg_parts := strings.Split(arg, "=")
			if len(arg_parts) < 2 {
				fmt.Errorf(`"%s" is not a propper argument, must have '='\n`)
				os.Exit(1)
			}

			switch(arg_parts[0]) {
			case "debug":
				switch(strings.ToLower(arg_parts[1])) {
				case "true":
					fmt.Printf("Debugging logs to stdout enabled.\n")
					debug = true
				case "false":
					debug = false
					fmt.Printf("Debugging logs to stdout disabled.\n")
				default:
					fmt.Errorf(`"%s": invalid argument to debug. Must be "true" or "false" (case-insensitive).`)
					os.Exit(1)
				}
			case "log-file":
				log_file_path = arg_parts[1]
			}

		}
	}
}

func main() {

	parseCliArgs(os.Args)


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

		go Server(client)

	}

}
