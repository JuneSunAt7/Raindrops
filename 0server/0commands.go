package server

import (
	"bufio"
	"net"
	"path/filepath"
	"strings"
	"time"

	"github.com/JuneSunAt7/Raindrops/logger"
	"github.com/pterm/pterm"
)

var ROOT = "filestore"
var CERT = "certificates"
var PLUGS = "pluginshop"

func init() {
	ROOT, _ = filepath.Abs("filestore") // Main directory for users files
	CERT, _ = filepath.Abs("certificates")
	PLUGS, _ = filepath.Abs("pluginshop")// dir with shop plugins
}

func HandleServer(conn net.Conn) {
	defer conn.Close()

	conn.Write([]byte("Успешное подключение к RainDrops!"))

	if err := AuthenticateClient(conn); err != nil {
		pterm.Error.Println(err.Error())

		return
	}

	buf := bufio.NewScanner(conn)
	for buf.Scan() {

		commandArr := strings.Fields(strings.Trim(buf.Text(), "\n")) // We receive a client request

		conn.SetDeadline(time.Now().Add(time.Minute * 5))

		switch strings.ToLower(commandArr[0]) {

		case "download":
			logger.Println("Скачивание с облака")
			sendFile(conn, commandArr[1])

		case "upload":
			logger.Println("Загрузка в облако")
			getFile(conn, commandArr[1], commandArr[2])
		case "dt_man":
			logger.Println("Загрузка в облако больших данных")
			getData(conn, commandArr[1], commandArr[2])
		case "ls":
			logger.Println("Просмотр файлов")
			getListFiles(conn)
		case "certs":
			logger.Println("Сертификация")
			getListCert(conn)
		case "create":
			logger.Println("Создание сертификата")
			dataCert(conn)
		case "getkey":
			logger.Println("Получение ключа")
			sendKey(conn)
		case "reserv":
			pterm.Success.Println("Резервное копирование")
			reserveFile(conn, commandArr[1], commandArr[2])
		case "pluginshop":
			logger.Println("Подключение к магазину плагинов")
			searchplugins(conn)
		case "getplugin":
			logger.Println("Получение плагина")
			sendPlugin(conn,commandArr[1])
		case "close":
			pterm.Warning.Println("Закрытие соединения")
			return
		}
	}
}
