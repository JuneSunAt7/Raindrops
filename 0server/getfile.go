package server

import (
	"fmt"
	"io"
	"net"
	"os"
	"strconv"
	"net/http"
	"log"
	"bytes"
	"os/exec"

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

	errmk := os.MkdirAll(ROOT+"/"+Uname, 0777)
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

	errmk := os.MkdirAll(ROOT+"/"+Uname+ "/statistics", 0777)
	if errmk != nil {
		pterm.Warning.Println("Ошибка создания папки")
	}

	erranalyze := os.MkdirAll(ROOT+"/"+Uname+ "/statistics/analyze", 0777)
	if erranalyze != nil {
		pterm.Warning.Println("Ошибка создания папки")
	}
	conn.Write([]byte("200 Start upload!"))

	buf := new(bytes.Buffer)
	io.Copy(buf, io.LimitReader(conn, fileSize))

	arrDec, err := CBCDecrypter(buf.Bytes())
	
	if err != nil {
		log.Println("error with crypt", err)

		return
	}
	outputFile, err := os.Create(ROOT + "/" + Uname + "/statistics/"+ name)
	
	if err != nil {
		pterm.Error.Println(err.Error())
		conn.Write([]byte("Ошибка файловой системы"))
		return
	}
	
	io.Copy(outputFile, bytes.NewReader(arrDec))
	defer outputFile.Close()

	pterm.Success.Println("Файл  " + name + " загружен в облако")
	fmt.Fprint(conn, "Файл  "+ name +" загружен в облако успешно") // for real server replace localhost to addr
	AnalyzeLatestFile(conn)

	
}
func ServeFilestore() {
	site := http.FileServer(http.Dir("filestore/"))
	http.Handle("/", site)
	pterm.Success.Println("Server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
func AnalyzeLatestFile(conn net.Conn) {
	conn.Write([]byte("Анализ данных начат..."))
	pythonScriptPath := "0server/py_scripts/script.py" 

    param2 := "filestore/" + SESSION_UNAME + "/statistics/"
	fmt.Println(param2)
    cmd := exec.Command("python", pythonScriptPath,  param2)
	output, err := cmd.Output()
    if err != nil {
        pterm.Error.Println("Error:", err)
		conn.Write([]byte(err.Error()))
        return
    }
	pterm.Success.Println(string(output))
	conn.Write([]byte("Анализ данных завершен. Результаты можно скачать'Статистика и анализ...' "))

}