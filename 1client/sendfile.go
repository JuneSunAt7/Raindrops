package client

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"

	"net"
	"strings"

	"os"
	"path/filepath"
	"time"

	"github.com/pterm/pterm"
)

func sendFile(conn net.Conn, fname string) {
	// That function use module crypto aka AES & MD5 hasing.
	//The server must make sure that the file is encrypted without errors.
	file := filepath.Base(filepath.Clean(fname))
	normalFilename := strings.Replace(file, " ", "_", -1)

	content, err := os.ReadFile(fname)
	if err != nil {
		pterm.Error.Println("Ошибка при загрузке файла")
		return
	}
	hash := md5.Sum([]byte(PASSWD))
	strPasswd := hex.EncodeToString(hash[:])

	arrEnc, err := CBCEncrypter(strPasswd, content)
	if err != nil {

		pterm.Error.Println("Ошибка при загрузке файла")
		return
	}
	conn.Write([]byte(fmt.Sprintf("upload %s %d\n", normalFilename, len(arrEnc))))

	buf := make([]byte, 1024)
	n, err := conn.Read(buf)
	if err != nil {
		pterm.Error.Println("Ошибка при загрузке файла")
		return
	}

	str := strings.Trim(string(buf[:n]), "\n")
	commandArr := strings.Fields(str)
	if commandArr[0] != "200" {
		pterm.Error.Println("Ошибка сетевого взаимодействия файла")
		return
	}

	io.Copy(conn, bytes.NewReader(arrEnc))

	n, err = conn.Read(buf)
	if err != nil {
		pterm.Error.Println("Ошибка при загрузке файла")
		return
	}
	p, _ := pterm.DefaultProgressbar.WithTotal(10).WithTitle("...Загрузка...").Start()

	for i := 0; i < p.Total; i++ {
		p.UpdateTitle("Загрузка в облако")
		p.Increment()
		time.Sleep(time.Millisecond * 50)
	}
	pterm.Success.Println(strings.Trim(string(buf[:n]), "\n"))

	CheckFileMD5Hash(fname)
}

func reserveSend(conn net.Conn, fname string) {

	file := filepath.Base(filepath.Clean(fname))

	dir, err := os.Open(fname)
	if err != nil {

		return
	}
	defer dir.Close()

	files, err := dir.Readdir(-1)
	if err != nil {

		return
	}
	for _, file := range files {
		if file.IsDir() {
			dirs(file.Name())
			pterm.Info.Println("Контенеризация папки ", file.Name())
		}
	}
	content, err := os.ReadFile(fname)
	if err != nil {
		pterm.Error.Println("Ошибка при загрузке файла " + file)

		return
	}

	arrEnc, err := CBCEncrypter(PASSWD, content)
	if err != nil {
		pterm.Error.Println("Ошибка при загрузке файла " + file)
		return
	}
	normalFilename := strings.Replace(file, " ", "", -1)

	conn.Write([]byte(fmt.Sprintf("reserv %s %d\n", normalFilename, len(arrEnc))))

	buf := make([]byte, 1024)
	n, err := conn.Read(buf)
	if err != nil {
		pterm.Error.Println("Ошибка при загрузке файла " + file)
		return
	}

	str := strings.Trim(string(buf[:n]), "\n")
	commandArr := strings.Fields(str)
	if commandArr[0] != "200" {
		pterm.Error.Println("Ошибка сетевого взаимодействия файла")
		return
	}

	io.Copy(conn, bytes.NewReader(arrEnc))

	n, err = conn.Read(buf)
	if err != nil {
		pterm.Error.Println("Ошибка при загрузке файла " + file)
		return
	}
	p, _ := pterm.DefaultProgressbar.WithTotal(10).WithTitle("...Загрузка...").Start()

	for i := 0; i < p.Total; i++ {
		p.UpdateTitle("Загрузка в облако")
		p.Increment()
		time.Sleep(time.Millisecond * 50)
	}
	pterm.Success.Println(strings.Trim(string(buf[:n]), "\n"))
}
