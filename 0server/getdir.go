package server

import (
	"fmt"
	"io/ioutil"
	"net"

	"github.com/JuneSunAt7/netMg/logger"
	"github.com/pterm/pterm"
)

func getListFiles(conn net.Conn) {

	files, err := ioutil.ReadDir(ROOT + "/" + Uname) // Read all filenames from filestore.
	if err != nil {
		pterm.Error.Println("Директория не была создана")
		conn.Write([]byte("Директория не была создана"))
		logger.Println(err.Error())
	}

	fileINFO := ""
	for _, file := range files {
		if !file.IsDir() {
			fileINFO += fmt.Sprintf("%-40s%-25s%-10d\n",
				file.Name(),
				file.ModTime().Format("2006-01-02 15:04:05"),
				file.Size())
		}

	}
	conn.Write([]byte(fileINFO))

}
