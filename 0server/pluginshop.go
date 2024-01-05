package server

import (
	"net"
	"os"
	"path/filepath"
	"strings"

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