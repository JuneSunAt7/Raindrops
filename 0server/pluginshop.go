package server

import(

	"path/filepath"
	"net"
	"os"
	"strings"
)

func searchplugins(conn net.Conn)([]string, error){
var files []string

	err := filepath.Walk(PLUGS, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			conn.Write([]byte("Ошибка подключения к магазину"))
			return err
		}

		if !info.IsDir() && filepath.Ext(path) == ".so" {
			files = append(files, path)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	if len(files) == 0{
		conn.Write([]byte("Пока что в магазине нет плагинов.\nНо скоро они появятся!"))
	}else{
		conn.Write([]byte("Доступные плагины"))
		conn.Write([]byte(strings.Join(files, "\n")))
	}
	return files,nil
}