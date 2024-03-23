package server

import (
	"fmt"
	"io"
	"net"
	"os"
	"strconv"
	"net/http"

	"github.com/pterm/pterm"
)
const PORT = ":8080"


func getDownloadLink(name string) string {
    return "http://127.0.0.1:8080/" +Uname +"/"+ name
}

func downloadHandler(w http.ResponseWriter, r *http.Request) {
	fileName := r.URL.Path[len("/download/"):]
	file, err := os.Open(ROOT + "/" + Uname + "/" + fileName)
	if err != nil {
		http.Error(w, "File not found", http.StatusNotFound)
		return
	}
	defer file.Close()

	// Set the content-type header to application/octet-stream
	w.Header().Set("Content-Type", "application/octet-stream")
	http.ServeFile(w, r, ROOT + "/" + Uname + "/" + fileName)
}

func getFile(conn net.Conn, name1 string, fs string) {

	fileSize, err := strconv.ParseInt(fs, 10, 64)
	if err != nil || fileSize <= 0 { // The size must not be less than or equal to zero!
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

	// Use buff size 32 bytes
	io.Copy(outputFile, io.LimitReader(conn, fileSize))
	pterm.Success.Println("Файл  " + name + " загружен в облако")
	fmt.Fprint(conn, "Файл  "+name+" загружен в облако успешно\nСсылка на файл: "+getDownloadLink(name))

	http.HandleFunc("/download/", downloadHandler)
    fmt.Println("File server is running on port" + PORT)
    if err := http.ListenAndServe(PORT, nil); err != nil {
        fmt.Println("Failed to start file server:", err)
    }
}