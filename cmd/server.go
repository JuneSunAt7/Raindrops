package main

import (
	"crypto/tls"
	"flag"
	"log"
	"net"
	"time"

	server "github.com/JuneSunAt7/Raindrops/0server"
	"github.com/JuneSunAt7/Raindrops/logger"

	
)

const (
	PORT = "2121"
)

func run() (err error) {

	var lstnr net.Listener

	boolTSL := flag.Bool("tls", true, "Set tls")
	flag.Parse()
	if !*boolTSL {

		lstnr, err = net.Listen("tcp", ":"+PORT)
		if err != nil {
			return err
		}

		logger.Println("TCP server is UP @ localhost without ssl: " + PORT)

	} else {

		cer, err := tls.LoadX509KeyPair("cert.pem", "key.pem")
		if err != nil {
			log.Println(err)
			return err
		}

		config := &tls.Config{Certificates: []tls.Certificate{cer}}

		lstnr, err = tls.Listen("tcp", ":"+PORT, config)
		if err != nil {
			return err
		}

		logger.Println("TCP TLS Server is UP @ localhost with ssl: " + PORT)
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
