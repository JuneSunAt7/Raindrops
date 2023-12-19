package main

import (
	"crypto/tls"
	"flag"
	"log"
	"net"
	"time"

	server "github.com/JuneSunAt7/netMg/0server"
	"github.com/JuneSunAt7/netMg/logger"
)

const (
	PORT = "2121"
)

func run() (err error) {

	var lstnr net.Listener

	boolTSL := flag.Bool("tls", false, "Set tls")
	flag.Parse()
	if !*boolTSL {

		lstnr, err = net.Listen("tcp", ":"+PORT)
		if err != nil {
			return err
		}

		logger.Println("TCP server is UP @ localhost: " + PORT)

	} else {

		cer, err := tls.LoadX509KeyPair("certs/server.crt", "certs/server.key")
		if err != nil {
			log.Println(err)
			return err
		}

		config := &tls.Config{Certificates: []tls.Certificate{cer}}

		lstnr, err = tls.Listen("tcp", ":"+PORT, config)
		if err != nil {
			return err
		}

		logger.Println("TCP TLS Server is UP @ localhost: " + PORT)

	}

	defer lstnr.Close()

	for {
		connection, err := lstnr.Accept()
		connection.SetDeadline(time.Now().Add(time.Minute * 2))
		if err != nil {
			logger.Println("Client Connection failed")
			continue
		}

		go server.HandleServer(connection)
	}

	// return nil
}

func main() {

	if err := run(); err != nil {
		log.Fatal(err)
	}

}
