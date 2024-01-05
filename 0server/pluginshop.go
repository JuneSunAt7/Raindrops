package server

import (
	"net"
	"os"
	"path/filepath"
	"strings"
	"fmt"
	"io"
	"github.com/pterm/pterm"

	"github.com/JuneSunAt7/Raindrops/logger"
)

func searchplugins(conn net.Conn){
var files []string

	err := filepath.Walk(PLUGS, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			conn.Write([]byte("Ошибка подключения к магазину"))
			return err
		}

		if !info.IsDir() && filepath.Ext(path) == ".so" {
			
			files = append(files, info.Name())
		}
		return nil
	})
	if err != nil {
		conn.Write([]byte("Ошибка подключения к магазину"))
		
	}
	logger.Println(files)
	if len(files) == 0{
		conn.Write([]byte("Пока что в магазине нет плагинов.\nНо скоро они появятся!"))
	}else{
		logger.Println("Поиск плагинов")
		str := strings.Join(files, "\n")
		conn.Write([]byte(str))
	}

}
func sendPlugin(conn net.Conn, name string){
	inputFile, err := os.Open(PLUGS +  "/" + name)
	if err != nil {
		pterm.Error.Println(err.Error())
		conn.Write([]byte(err.Error()))
		return
	}
	defer inputFile.Close()

	stats, _ := inputFile.Stat()

	conn.Write([]byte(fmt.Sprintf("plugin %s %d\n", name, stats.Size())))

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
	pterm.Success.Println("Плагин ", name, " отправлен")
}