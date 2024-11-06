package cloudcompute
import (
	"net"
	"strings"
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"
	"crypto/md5"
	"encoding/hex"
	"strconv"



	"log"
	"github.com/harry1453/go-common-file-dialog/cfd"
	crypto "github.com/JuneSunAt7/Raindrops/1client"
	"github.com/pterm/pterm"
)

func UploadBigData(conn net.Conn) {
	fname := ChooseFileData()
	fname = strings.Replace(fname, "\\", "/", -1)
	sendData(conn, fname)
}
func ChooseFileData() string {
	openDialog, err := cfd.NewOpenFileDialog(cfd.DialogConfig{
		Title: "Open A File",
		Role:  "OpenFile",
		FileFilters: []cfd.FileFilter{
			{DisplayName: "CSV files", Pattern: "*.csv"},
			{DisplayName: "DB files", Pattern: "*.db"},
			{DisplayName: "XLSX files", Pattern: "*.xlsx"},
		},

	})
	if err != nil {
		log.Fatal(err)
	}

	if err := openDialog.Show(); err != nil {
		log.Fatal(err)
	}
	result, err := openDialog.GetResult()
	if err == cfd.ErrorCancelled {
		pterm.Error.Println("Вы закрыли окно выбора файла")
	} else if err != nil {
		log.Fatal(err)
	}
	return result
}

func sendData(conn net.Conn, fname string) {
	// That function use module crypto aka AES & MD5 hasing.
	//The server must make sure that the file is encrypted without errors.
	file := filepath.Base(filepath.Clean(fname))
	normalFilename := strings.Replace(file, " ", "_", -1)

	content, err := os.ReadFile(fname)
	if err != nil {
		pterm.Error.Println("Ошибка при загрузке файла")
		return
	}
	hash := md5.Sum([]byte(crypto.PASSWD))
	strPasswd := hex.EncodeToString(hash[:])

	arrEnc, err := crypto.CBCEncrypter(strPasswd, content)
	if err != nil {

		pterm.Error.Println("Ошибка при загрузке файла")
		return
	}
	conn.Write([]byte(fmt.Sprintf("dt_man %s %d\n", normalFilename, len(arrEnc))))

	buf := make([]byte, 1024)
	n, err := conn.Read(buf)
	if err != nil {
		pterm.Error.Println("Ошибка при загрузке файла")
		return
	}

	str := strings.Trim(string(buf[:n]), "\n")
	commandArr := strings.Fields(str)
	if commandArr[0] != "200" {
		pterm.Error.Println("Ошибка сетевого взаимодействия файла")
		return
	}

	io.Copy(conn, bytes.NewReader(arrEnc))

	n, err = conn.Read(buf)
	if err != nil {
		pterm.Error.Println("Ошибка при загрузке файла")
		return
	}
	p, _ := pterm.DefaultProgressbar.WithTotal(10).WithTitle("...Загрузка...").Start()

	for i := 0; i < p.Total; i++ {
		p.UpdateTitle("Загрузка в облако")
		p.Increment()
		time.Sleep(time.Millisecond * 50)
	}
	pterm.Success.Println(strings.Trim(string(buf[:n]), "\n"))

	crypto.CheckFileMD5Hash(fname)

	answerBuf := make([]byte, 1024)
	n, err = conn.Read(answerBuf)
	if err != nil {
		pterm.Error.Println("Ошибка сети")
		return
	}
	pterm.Success.Println(strings.Trim(string(answerBuf[:n]), "\n"))
	time.Sleep(time.Millisecond * 100)
	
	readyBuf := make([]byte, 1024)
	n, err = conn.Read(readyBuf)
	if err != nil {
		pterm.Error.Println("Ошибка сети")
		return
	}
	pterm.Success.Println(strings.Trim(string(readyBuf[:n]), "\n"))
}

func ListStatistics(conn net.Conn) {
	conn.Write([]byte("list_stat\n"))
	
	buf := make([]byte, 1024)
	n, err := conn.Read(buf)
	if err != nil {
		pterm.Error.Println("Ошибка сети")
		return
	}
	pterm.FgBlue.Println(strings.Trim(string(buf[:n]), "\n"))
	fname, _ := pterm.DefaultInteractiveTextInput.Show("Имя файла")
	getFile(conn, fname)

}
func getFile(conn net.Conn, fname string) {
    file := filepath.Base(fname)
    usersDir := crypto.ChooseDir()

    conn.Write([]byte(fmt.Sprintf("down_stat %s\n", file)))

    buffer := make([]byte, 1024)
    n, _ := conn.Read(buffer)
    comStr := strings.Trim(string(buffer[:n]), "\n")
    commandArr := strings.Fields(comStr)

    fileSize, err := strconv.ParseInt(commandArr[2], 10, 64)
    if err != nil || fileSize == -1 {
        log.Println("file size error", err)
        conn.Write([]byte("file size error"))
        return
    }

    conn.Write([]byte("200 Start download!"))

    buf := new(bytes.Buffer)
    io.Copy(buf, io.LimitReader(conn, fileSize))

    outputFile, err := os.Create(usersDir + "/" + file)
    if err != nil {
        log.Println("error create dir", err)
    }
    io.Copy(outputFile, buf)
    defer outputFile.Close()

    p, _ := pterm.DefaultProgressbar.WithTotal(5).WithTitle("...Скачивание файла...").Start()

    for i := 0; i < p.Total; i++ {
        p.UpdateTitle("Выгрузка из облака") 
        p.Increment()
        time.Sleep(time.Millisecond * 350)
    }
    
}