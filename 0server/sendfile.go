package server

import (
	"fmt"
	"io"
	"net"
	"os"
	"strings"

	"github.com/pterm/pterm"
)

func sendFile(conn net.Conn, name string) {

	inputFile, err := os.Open(ROOT + "/" + Uname + "/" + name)
	if err != nil {
		pterm.Error.Println(err.Error())
		conn.Write([]byte(err.Error()))
		return
	}
	defer inputFile.Close()

	stats, _ := inputFile.Stat()

	conn.Write([]byte(fmt.Sprintf("download %s %d\n", name, stats.Size())))

	buf := make([]byte, 1024)
	n, err := conn.Read(buf)
	if err != nil {
		pterm.Warning.Println(err.Error())
		return
	}

	str := strings.Trim(string(buf[:n]), "\n")
	commandArr := strings.Fields(str)
	if commandArr[0] != "200" { // Ansver-code 200(succesfully)
		pterm.Warning.Println(str)
		return
	}

	io.Copy(conn, inputFile)
	pterm.Success.Println("Файл ", name, " отправлен")
}
func sendKey(conn net.Conn) {
	key, err := os.Open(CERT + "/" + Uname + "/" + Uname + ".crt")
	if err != nil {
		pterm.Error.Println(err.Error())
		conn.Write([]byte("Нет ключа шифрования. Создайте его в Сертификаты и пароли"))
		return
	}
	defer key.Close()

	stats, _ := key.Stat()

	conn.Write([]byte(fmt.Sprintf("download  %d\n", stats.Size())))

	buf := make([]byte, 1024)
	n, err := conn.Read(buf)
	if err != nil {
		pterm.Error.Println(err.Error())
		return
	}

	str := strings.Trim(string(buf[:n]), "\n")
	commandArr := strings.Fields(str)
	if commandArr[0] != "200" { // Ansver-code 200(succesfully)
		pterm.Warning.Println(str)
		return
	}

	io.Copy(conn, key)
	pterm.Success.Println("Ключ ", key, " отправлен")
}
