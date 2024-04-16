package client

import (
	"bytes"
	"fmt"
	"io"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/pterm/pterm"

	"log"
	"net"
	"os"
)

func getFile(conn net.Conn, fname string, myFPass string) {
	file := filepath.Base(fname)
	usersDir := ChooseDir()

	conn.Write([]byte(fmt.Sprintf("download %s\n", file)))

	buffer := make([]byte, 1024)
	n, _ := conn.Read(buffer)
	comStr := strings.Trim(string(buffer[:n]), "\n")
	commandArr := strings.Fields(comStr)

	fileSize, err := strconv.ParseInt(commandArr[2], 10, 64)
	if err != nil || fileSize == -1 {
		log.Println("file size error", err)
		conn.Write([]byte("file size error"))
		return
	}

	conn.Write([]byte("200 Start download!"))

	buf := new(bytes.Buffer)
	io.Copy(buf, io.LimitReader(conn, fileSize))

	arrDec, err := CBCDecrypter(myFPass, buf.Bytes())
	if err != nil {
		log.Println("error with crypt", err)

		return
	}

	outputFile, err := os.Create(usersDir +"/" + file)
	if err != nil {
		log.Println("error create dir", err)
	}
	io.Copy(outputFile, bytes.NewReader(arrDec))
	defer outputFile.Close()
	p, _ := pterm.DefaultProgressbar.WithTotal(5).WithTitle("...Скачивание файла...").Start()

	for i := 0; i < p.Total; i++ {
		p.UpdateTitle("Выгрузка из облака") 
		p.Increment()
		time.Sleep(time.Millisecond * 350)
	}
	checkFileMD5Hash(usersDir+ "/"+ fname)
}
