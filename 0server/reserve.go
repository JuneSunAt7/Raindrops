package server

import (
	"fmt"
	"io"
	"net"
	"os"
	"strconv"

	"github.com/pterm/pterm"
)

func reserveFile(conn net.Conn, name1 string, fs string) {

	fileSize, err := strconv.ParseInt(fs, 10, 64)
	if err != nil || fileSize == -1 { // The size must not be less than zero!
		pterm.Error.Println(err.Error())
		conn.Write([]byte("file size error"))
		return
	}

	name := name1

	errmk := os.Mkdir(ROOT+"/"+Uname, 0777)
	if errmk != nil {
		pterm.Warning.Println("Ошибка создания папки")
	}

	outputFile, err := os.Create(ROOT + "/" + Uname + "/" + name)
	if err != nil {
		pterm.Error.Println(err.Error())
		conn.Write([]byte("Ошибка файловой системы"))
		return
	}
	defer outputFile.Close()

	conn.Write([]byte("200 Start upload!"))

	//Use buff size 32 bytes
	io.Copy(outputFile, io.LimitReader(conn, fileSize))
	pterm.Success.Println("Файл  " + name + " загружен в облако")
	fmt.Fprint(conn, "Файл  "+name+" загружен в облако успешно\n")

}
