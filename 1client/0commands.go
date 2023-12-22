package client

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"strings"
	"time"

	"atomicgo.dev/keyboard/keys"
	"github.com/pterm/pterm"
)

var ROOT = "filestore/storeclient"

// local settings
func init() {
	ROOT, _ = filepath.Abs("filestore/storeclient")
}

func Upload(conn net.Conn) {
	fname := ChooseFile()
	fname = strings.Replace(fname, "\\", "/", -1)
	sendFile(conn, fname)
}

func Download(conn net.Conn) {
	fname, _ := pterm.DefaultInteractiveTextInput.Show("Имя файла")
	passwd := PASSWD
	getFile(conn, fname, passwd+"\n")
}

func ListFiles(conn net.Conn) {
	conn.Write([]byte("ls\n"))
	buffer := make([]byte, 4096)
	n, _ := conn.Read(buffer)
	pterm.FgGreen.Println(string(buffer[:n]))

}

func Exit(conn net.Conn) {
	conn.Write([]byte("close\n"))
	pterm.FgGreen.Println("Выход из облака")
}

func CertPasswd(conn net.Conn) {
	var certoptions []string

	certoptions = append(certoptions, fmt.Sprintf("Доступные сертификаты"))
	certoptions = append(certoptions, fmt.Sprintf("Создать сертификат"))
	certoptions = append(certoptions, fmt.Sprintf("Назад"))

	printer := pterm.DefaultInteractiveMultiselect.WithOptions(certoptions)
	printer.Filter = false
	printer.KeyConfirm = keys.Enter
	for {
		selectedOptions, _ := pterm.DefaultInteractiveSelect.WithOptions(certoptions).Show()
		switch selectedOptions {
		case "Доступные сертификаты":
			AllCertificates(conn)
		case "Создать сертификат":
			CreateCert(conn)
		case "Назад":
			return
		}
	}
}

func Autoreserved() {
	var options []string

	options = append(options, fmt.Sprintf("Календарь авторезервирования"))
	options = append(options, fmt.Sprintf("Файлы для авторезервирования"))
	options = append(options, fmt.Sprintf("Контейнеры"))
	options = append(options, fmt.Sprintf("Настройки"))
	options = append(options, fmt.Sprintf("Назад"))

	printer := pterm.DefaultInteractiveMultiselect.WithOptions(options)
	printer.Filter = false
	printer.KeyConfirm = keys.Enter
	for {
		selectedOptions, _ := pterm.DefaultInteractiveSelect.WithOptions(options).Show()
		switch selectedOptions {
		case "Календарь авторезервирования":
			Calendar()
		case "Файлы для авторезервирования":
			Userfiles()
		case "Контейнеры":
			Containers()
		case "Настройки":
			Setting()
		case "Назад":
			return
		}
	}
}

func AutoSendFiles(conn net.Conn) {

	file, err := os.Open(ROOT + "/" + "localSettings" + "/" + "path.ini")
	if err != nil {
		pterm.FgRed.Println("Файлы для резервирования не найдены!")
		return
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	for j := 0; j < len(lines); j++ {

		fname := strings.Replace(lines[j], "\\", "/", -1)
		lastChar := fname[len(fname)-1]
		if string(lastChar) == "/"{// wow if it isnt directory??
			dirname := dirs(fname)
			reserveSend(conn, dirname)
		}else{
		reserveSend(conn, fname)
		}
		if j > 10 {
			time.Sleep(time.Millisecond * 30)
		}
	}
}
