package server

import (
	"fmt"
	"io"
	"net"
	"os"
	"strconv"
	"net/http"
	"log"

	"github.com/pterm/pterm"
)

func getFile(conn net.Conn, name1 string, fs string) {

	fileSize, err := strconv.ParseInt(fs, 10, 64)
	if err != nil || fileSize == -1 {
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

	io.Copy(outputFile, io.LimitReader(conn, fileSize))
	pterm.Success.Println("Файл  " + name + " загружен в облако")
	fmt.Fprint(conn, "Файл  "+ name +" загружен в облако успешно\nСсылка на облако: http://localhost:8080/"+Uname) // for real server replace localhost to addr

	go ServeFilestore()
}
func getData(conn net.Conn, name1 string, fs string) {
	fileSize, err := strconv.ParseInt(fs, 10, 64)
	if err != nil || fileSize == -1 {
		pterm.Error.Println(err.Error())
		conn.Write([]byte("file size error"))
		return
	}

	name := name1

	errmk := os.Mkdir(ROOT+"/"+Uname+ "/statistics", 0777)
	if errmk != nil {
		pterm.Warning.Println("Ошибка создания папки")
	}

	outputFile, err := os.Create(ROOT + "/" + Uname + "/statistics/"+ name)
	
	if err != nil {
		pterm.Error.Println(err.Error())
		conn.Write([]byte("Ошибка файловой системы"))
		return
	}
	defer outputFile.Close()

	conn.Write([]byte("200 Start upload!"))

	io.Copy(outputFile, io.LimitReader(conn, fileSize))
	pterm.Success.Println("Файл  " + name + " загружен в облако")
	fmt.Fprint(conn, "Файл  "+ name +" загружен в облако успешно") // for real server replace localhost to addr

}
func ServeFilestore() {
	site := http.FileServer(http.Dir("filestore/"))
	http.Handle("/", site)
	pterm.Success.Println("Server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
