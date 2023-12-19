package server

import (
	"bufio"
	"net"
	"path/filepath"
	"strings"
	"time"

	"github.com/JuneSunAt7/netMg/logger"
)

var ROOT = "filestore"
var CERT = "certificates"

func init() {
	ROOT, _ = filepath.Abs("filestore") // Main directory for users files
	CERT, _ = filepath.Abs("certificates")
}

func HandleServer(conn net.Conn) {
	defer conn.Close()

	conn.Write([]byte("Server:Connection Established"))

	if err := AuthenticateClient(conn); err != nil {
		logger.Println(err.Error())
		return
	}

	buf := bufio.NewScanner(conn)
	for buf.Scan() {

		commandArr := strings.Fields(strings.Trim(buf.Text(), "\n")) // We receive a client request

		conn.SetDeadline(time.Now().Add(time.Minute * 5))

		switch strings.ToLower(commandArr[0]) {

		case "download":
			logger.Println("Download Request")
			sendFile(conn, commandArr[1])

		case "upload":
			logger.Println("Upload Request")
			getFile(conn, commandArr[1], commandArr[2])
		case "ls":
			logger.Println("ls")
			getListFiles(conn)
		case "certs":
			logger.Println("certs")
			getListCert(conn)
		case "create":
			logger.Println("create cert")
			dataCert(conn)
		case "getkey":
			logger.Println("Upload file with cert")
			sendKey(conn)
		case "close":
			logger.Println("closed")
			return
		}
	}
}
