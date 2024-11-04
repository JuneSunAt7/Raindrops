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
	fmt.Println(crypto.PASSWD)
	arrEnc, err := crypto.CBCEncrypter(crypto.PASSWD, content)
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
}